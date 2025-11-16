package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/supercakecrumb/affirm-name/internal/config"
	"github.com/supercakecrumb/affirm-name/internal/db"
	"github.com/supercakecrumb/affirm-name/internal/handlers"
	"github.com/supercakecrumb/affirm-name/internal/middleware"
	"go.uber.org/zap"
)

func main() {
	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 3. Initialize database connection
	if !cfg.FixtureMode {
		database, err := db.New(cfg.DatabaseURL)
		if err != nil {
			logger.Fatal("Failed to connect to database", zap.Error(err))
		}
		defer database.Close()
		cfg.DB = database
		logger.Info("Database connected successfully")
	}

	// 4. Create chi router
	r := chi.NewRouter()

	// 5. Add middleware
	r.Use(middleware.Logger(logger))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// 6. Register routes
	r.Get("/api/meta/years", handlers.MetaYears(cfg))
	r.Get("/api/meta/countries", handlers.MetaCountries(cfg))
	r.Get("/api/names", handlers.NamesList(cfg))
	r.Get("/api/names/trend", handlers.NameTrend(cfg))

	// Log startup information
	logger.Info("Server starting",
		zap.String("port", cfg.Port),
		zap.Bool("fixture_mode", cfg.FixtureMode),
		zap.String("frontend_url", cfg.FrontendURL),
	)

	// 7. Start server
	addr := ":" + cfg.Port
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
