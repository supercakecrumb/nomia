package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/supercakecrumb/affirm-name/internal/config"
	"github.com/supercakecrumb/affirm-name/internal/database"
	"github.com/supercakecrumb/affirm-name/internal/logging"
	"github.com/supercakecrumb/affirm-name/internal/repository"
	"github.com/supercakecrumb/affirm-name/internal/service"
	"github.com/supercakecrumb/affirm-name/internal/storage"
	"github.com/supercakecrumb/affirm-name/internal/worker"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logging.InitLogger(cfg.Logging.Level, cfg.Logging.Format)
	logger.Info("Starting background worker",
		slog.String("version", "1.0.0"),
		slog.Int("concurrency", cfg.Worker.Concurrency),
		slog.Duration("poll_interval", cfg.Worker.PollInterval),
		slog.Int("max_retries", cfg.Worker.MaxRetries),
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
	datasetRepo := repository.NewDatasetRepository(db.Sqlx)
	jobRepo := repository.NewJobRepository(db.Sqlx)
	nameRepo := repository.NewNameRepository(db.Sqlx)

	// Initialize services
	parserService := service.NewParserService(datasetRepo, nameRepo, jobRepo, stor, logger)

	// Initialize processor
	processor := worker.NewProcessor(
		jobRepo,
		datasetRepo,
		parserService,
		stor,
		cfg.Worker.MaxRetries,
		logger,
	)

	// Initialize worker pool
	pool := worker.NewPool(
		processor,
		cfg.Worker.Concurrency,
		cfg.Worker.PollInterval,
		logger,
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker pool
	if err := pool.Start(ctx); err != nil {
		logger.Error("Failed to start worker pool", slog.Any("error", err))
		os.Exit(1)
	}

	// Wait for interrupt signal to gracefully shutdown the worker
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down worker...")

	// Cancel context to stop workers
	cancel()

	// Stop worker pool (waits for all workers to finish)
	pool.Stop()

	logger.Info("Worker exited")
}
