package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/supercakecrumb/affirm-name/internal/api/handlers"
	"github.com/supercakecrumb/affirm-name/internal/api/middleware"
	"github.com/supercakecrumb/affirm-name/internal/config"
	"github.com/supercakecrumb/affirm-name/internal/database"
	"github.com/supercakecrumb/affirm-name/internal/logging"
	"github.com/supercakecrumb/affirm-name/internal/repository"
	"github.com/supercakecrumb/affirm-name/internal/service"
	"github.com/supercakecrumb/affirm-name/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logging.InitLogger(cfg.Logging.Level, cfg.Logging.Format)
	logger.Info("Starting API server",
		slog.String("version", "1.0.0"),
		slog.String("port", cfg.Server.Port),
	)

	// Connect to database
	db, err := database.New(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Connected to database")

	// Initialize storage
	var stor storage.Storage
	if cfg.Storage.Type == "s3" {
		s3Storage, err := storage.NewS3Storage(context.Background(), storage.S3Config{
			Bucket:   cfg.Storage.S3Bucket,
			Region:   cfg.Storage.S3Region,
			Endpoint: cfg.Storage.S3Endpoint,
		})
		if err != nil {
			logger.Error("Failed to initialize S3 storage", slog.Any("error", err))
			os.Exit(1)
		}
		stor = s3Storage
		logger.Info("Initialized S3 storage", slog.String("bucket", cfg.Storage.S3Bucket))
	} else {
		localStorage, err := storage.NewLocalStorage(cfg.Storage.Path)
		if err != nil {
			logger.Error("Failed to initialize local storage", slog.Any("error", err))
			os.Exit(1)
		}
		stor = localStorage
		logger.Info("Initialized local storage", slog.String("path", cfg.Storage.Path))
	}

	// Initialize repositories
	countryRepo := repository.NewCountryRepository(db.Pool)
	datasetRepo := repository.NewDatasetRepository(db.Sqlx)
	jobRepo := repository.NewJobRepository(db.Sqlx)
	nameRepo := repository.NewNameRepository(db.Sqlx)

	// Initialize services
	countryService := service.NewCountryService(countryRepo)
	uploadService := service.NewUploadService(stor, datasetRepo, jobRepo, countryRepo)
	nameService := service.NewNameService(nameRepo, countryRepo)

	// Initialize handlers
	countryHandler := handlers.NewCountryHandler(countryService)
	uploadHandler := handlers.NewUploadHandler(uploadService)
	datasetHandler := handlers.NewDatasetHandler(datasetRepo)
	jobHandler := handlers.NewJobHandler(jobRepo)
	nameHandler := handlers.NewNameHandler(nameService)

	// Set Gin mode based on log level
	if cfg.Logging.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Apply middleware
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())

	// Health check endpoints
	router.GET("/health", healthCheck(db))
	router.GET("/ready", readinessCheck(db))

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Country routes
		countries := v1.Group("/countries")
		{
			countries.POST("", countryHandler.Create)
			countries.GET("", countryHandler.List)
			countries.GET("/:id", countryHandler.Get)
			countries.PATCH("/:id", countryHandler.Update)
			countries.DELETE("/:id", countryHandler.Delete)
		}

		// Dataset routes
		datasets := v1.Group("/datasets")
		{
			datasets.POST("/upload", uploadHandler.Upload)
			datasets.GET("/upload/limits", uploadHandler.GetUploadLimits)
			datasets.GET("", datasetHandler.List)
			datasets.GET("/:id", datasetHandler.Get)
			datasets.DELETE("/:id", datasetHandler.Delete)
		}

		// Job routes
		jobs := v1.Group("/jobs")
		{
			jobs.GET("", jobHandler.List)
			jobs.GET("/:id", jobHandler.Get)
		}

		// Name routes
		names := v1.Group("/names")
		{
			names.GET("", nameHandler.List)
			names.GET("/search", nameHandler.Search)
			names.GET("/:name", nameHandler.GetByName)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting HTTP server",
			slog.String("port", cfg.Server.Port),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("Server exited")
}

// healthCheck returns a handler for the health check endpoint
func healthCheck(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database connection
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		dbStatus := "connected"
		if err := db.Pool.Ping(ctx); err != nil {
			dbStatus = "disconnected"
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"checks": gin.H{
					"database": dbStatus,
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"checks": gin.H{
				"database": dbStatus,
			},
		})
	}
}

// readinessCheck returns a handler for the readiness check endpoint
func readinessCheck(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database connection
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := db.Pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"ready": false,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"ready": true,
		})
	}
}
