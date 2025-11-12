package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Register pgx driver for database/sql
	"github.com/jmoiron/sqlx"
	"github.com/supercakecrumb/affirm-name/internal/config"
)

// DB wraps both pgxpool and sqlx connections
type DB struct {
	Pool *pgxpool.Pool
	Sqlx *sqlx.DB
}

// New creates a new database connection pool
func New(cfg config.DatabaseConfig) (*DB, error) {
	// Parse connection string and create pool config
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool settings
	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MinConns = int32(cfg.MaxIdle)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime
	poolConfig.MaxConnIdleTime = cfg.ConnMaxIdleTime

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Also create sqlx connection for repositories that need it
	sqlxDB, err := sqlx.Connect("pgx", cfg.URL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to create sqlx connection: %w", err)
	}

	// Configure sqlx connection pool
	sqlxDB.SetMaxOpenConns(cfg.MaxConnections)
	sqlxDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlxDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlxDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	slog.Info("database connection established",
		"max_connections", cfg.MaxConnections,
		"max_idle", cfg.MaxIdle,
	)

	return &DB{
		Pool: pool,
		Sqlx: sqlxDB,
	}, nil
}

// Close closes the database connection pools
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
	if db.Sqlx != nil {
		db.Sqlx.Close()
	}
	slog.Info("database connection closed")
}

// HealthCheck performs a health check on the database
func (db *DB) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.Pool.Ping(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Stats returns connection pool statistics
func (db *DB) Stats() *pgxpool.Stat {
	return db.Pool.Stat()
}
