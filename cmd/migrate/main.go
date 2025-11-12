package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/supercakecrumb/affirm-name/internal/config"
	"github.com/supercakecrumb/affirm-name/internal/logging"
)

const (
	migrationsTable = "schema_migrations"
)

func main() {
	// Initialize logger
	logging.InitLogger("info", "text")

	// Check command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <up|down|status>")
		os.Exit(1)
	}

	command := os.Args[1]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Connect to database
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.Database.URL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(ctx, conn); err != nil {
		slog.Error("failed to create migrations table", "error", err)
		os.Exit(1)
	}

	// Execute command
	switch command {
	case "up":
		if err := migrateUp(ctx, conn); err != nil {
			slog.Error("migration up failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migrations applied successfully")
	case "down":
		if err := migrateDown(ctx, conn); err != nil {
			slog.Error("migration down failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migration rolled back successfully")
	case "status":
		if err := showStatus(ctx, conn); err != nil {
			slog.Error("failed to show status", "error", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage: migrate <up|down|status>")
		os.Exit(1)
	}
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist
func createMigrationsTable(ctx context.Context, conn *pgx.Conn) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`
	_, err := conn.Exec(ctx, query)
	return err
}

// getAppliedMigrations returns a list of applied migration versions
func getAppliedMigrations(ctx context.Context, conn *pgx.Conn) (map[string]bool, error) {
	rows, err := conn.Query(ctx, "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// getMigrationFiles returns a sorted list of migration files
func getMigrationFiles(direction string) ([]string, error) {
	migrationsDir := "migrations"
	pattern := fmt.Sprintf("*.%s.sql", direction)

	files, err := filepath.Glob(filepath.Join(migrationsDir, pattern))
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

// extractVersion extracts the version from a migration filename
func extractVersion(filename string) string {
	base := filepath.Base(filename)
	// Remove .up.sql or .down.sql suffix
	parts := strings.Split(base, ".")
	if len(parts) >= 3 {
		return parts[0]
	}
	return base
}

// migrateUp applies all pending migrations
func migrateUp(ctx context.Context, conn *pgx.Conn) error {
	// Get applied migrations
	applied, err := getAppliedMigrations(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get migration files
	files, err := getMigrationFiles("up")
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	if len(files) == 0 {
		slog.Info("no migration files found")
		return nil
	}

	// Apply pending migrations
	for _, file := range files {
		version := extractVersion(file)

		if applied[version] {
			slog.Debug("migration already applied", "version", version)
			continue
		}

		slog.Info("applying migration", "version", version, "file", file)

		// Read migration file
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Execute migration in a transaction
		tx, err := conn.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Execute migration SQL
		if _, err := tx.Exec(ctx, string(content)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to execute migration %s: %w", version, err)
		}

		// Record migration
		if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", version, err)
		}

		slog.Info("migration applied", "version", version)
	}

	return nil
}

// migrateDown rolls back the last applied migration
func migrateDown(ctx context.Context, conn *pgx.Conn) error {
	// Get applied migrations
	applied, err := getAppliedMigrations(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) == 0 {
		slog.Info("no migrations to roll back")
		return nil
	}

	// Get migration files
	upFiles, err := getMigrationFiles("up")
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Find the last applied migration
	var lastVersion string
	for i := len(upFiles) - 1; i >= 0; i-- {
		version := extractVersion(upFiles[i])
		if applied[version] {
			lastVersion = version
			break
		}
	}

	if lastVersion == "" {
		slog.Info("no migrations to roll back")
		return nil
	}

	// Find corresponding down file
	downFile := filepath.Join("migrations", lastVersion+".down.sql")
	if _, err := os.Stat(downFile); os.IsNotExist(err) {
		return fmt.Errorf("down migration file not found: %s", downFile)
	}

	slog.Info("rolling back migration", "version", lastVersion, "file", downFile)

	// Read migration file
	content, err := os.ReadFile(downFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %w", downFile, err)
	}

	// Execute migration in a transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute migration SQL
	if _, err := tx.Exec(ctx, string(content)); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to execute migration %s: %w", lastVersion, err)
	}

	// Remove migration record
	if _, err := tx.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", lastVersion); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to remove migration record %s: %w", lastVersion, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit rollback %s: %w", lastVersion, err)
	}

	slog.Info("migration rolled back", "version", lastVersion)
	return nil
}

// showStatus shows the current migration status
func showStatus(ctx context.Context, conn *pgx.Conn) error {
	// Get applied migrations
	applied, err := getAppliedMigrations(ctx, conn)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get migration files
	files, err := getMigrationFiles("up")
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	fmt.Println("\nMigration Status:")
	fmt.Println("=================")

	if len(files) == 0 {
		fmt.Println("No migration files found")
		return nil
	}

	for _, file := range files {
		version := extractVersion(file)
		status := "pending"
		if applied[version] {
			status = "applied"
		}
		fmt.Printf("%-50s %s\n", version, status)
	}

	fmt.Printf("\nTotal: %d migrations, %d applied, %d pending\n",
		len(files), len(applied), len(files)-len(applied))

	return nil
}
