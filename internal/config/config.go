package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Storage  StorageConfig
	Server   ServerConfig
	Worker   WorkerConfig
	Logging  LoggingConfig
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	URL             string
	MaxConnections  int
	MaxIdle         int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// StorageConfig holds file storage settings
type StorageConfig struct {
	Type string // "local" or "s3"
	Path string // for local storage

	// S3 settings
	S3Bucket   string
	S3Region   string
	S3Endpoint string
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// WorkerConfig holds background worker settings
type WorkerConfig struct {
	Concurrency  int
	PollInterval time.Duration
	MaxRetries   int
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Level  string // debug, info, warn, error
	Format string // json or text
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", ""),
			MaxConnections:  getEnvInt("DATABASE_MAX_CONNECTIONS", 100),
			MaxIdle:         getEnvInt("DATABASE_MAX_IDLE", 10),
			ConnMaxLifetime: getEnvDuration("DATABASE_CONN_MAX_LIFETIME", time.Hour),
			ConnMaxIdleTime: getEnvDuration("DATABASE_CONN_MAX_IDLE_TIME", 10*time.Minute),
		},
		Storage: StorageConfig{
			Type:       getEnv("STORAGE_TYPE", "local"),
			Path:       getEnv("STORAGE_PATH", "./uploads"),
			S3Bucket:   getEnv("S3_BUCKET", ""),
			S3Region:   getEnv("S3_REGION", "us-east-1"),
			S3Endpoint: getEnv("S3_ENDPOINT", ""),
		},
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Worker: WorkerConfig{
			Concurrency:  getEnvInt("WORKER_CONCURRENCY", 4),
			PollInterval: getEnvDuration("WORKER_POLL_INTERVAL", 5*time.Second),
			MaxRetries:   getEnvInt("WORKER_MAX_RETRIES", 3),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that required configuration is present
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.Storage.Type != "local" && c.Storage.Type != "s3" {
		return fmt.Errorf("STORAGE_TYPE must be 'local' or 's3', got: %s", c.Storage.Type)
	}

	if c.Storage.Type == "local" && c.Storage.Path == "" {
		return fmt.Errorf("STORAGE_PATH is required when STORAGE_TYPE is 'local'")
	}

	if c.Storage.Type == "s3" && c.Storage.S3Bucket == "" {
		return fmt.Errorf("S3_BUCKET is required when STORAGE_TYPE is 's3'")
	}

	if c.Logging.Level != "debug" && c.Logging.Level != "info" &&
		c.Logging.Level != "warn" && c.Logging.Level != "error" {
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}

	if c.Logging.Format != "json" && c.Logging.Format != "text" {
		return fmt.Errorf("LOG_FORMAT must be 'json' or 'text'")
	}

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets an integer environment variable or returns a default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDuration gets a duration environment variable or returns a default value
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
