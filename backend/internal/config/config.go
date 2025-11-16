package config

import (
	"github.com/spf13/viper"
	"github.com/supercakecrumb/affirm-name/internal/db"
)

// Config holds the application configuration
type Config struct {
	FixtureMode bool
	Port        string
	DatabaseURL string
	FrontendURL string
	DB          *db.DB // Database connection pool
}

// Load reads configuration from .env file and environment variables
func Load() (*Config, error) {
	// Set default values
	viper.SetDefault("FIXTURE_MODE", false)
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/affirm_name?sslmode=disable")
	viper.SetDefault("FRONTEND_URL", "http://localhost:5173")

	// Read from .env file
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Attempt to read the config file
	// It's okay if the .env file doesn't exist, we'll use defaults and env vars
	_ = viper.ReadInConfig()

	// Environment variables take precedence
	viper.AutomaticEnv()

	// Build config struct
	cfg := &Config{
		FixtureMode: viper.GetBool("FIXTURE_MODE"),
		Port:        viper.GetString("PORT"),
		DatabaseURL: viper.GetString("DATABASE_URL"),
		FrontendURL: viper.GetString("FRONTEND_URL"),
	}

	return cfg, nil
}
