#!/bin/sh
set -e

# =============================================================================
# Docker Initialization Script
# =============================================================================
# This script is used to initialize the application in Docker containers.
# It waits for PostgreSQL to be ready and runs migrations.
# =============================================================================

echo "==> Docker initialization script started"

# Configuration
MAX_RETRIES=30
RETRY_INTERVAL=2
DB_HOST="${DB_HOST:-postgres}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-affirm_user}"
DB_NAME="${DB_NAME:-affirm_name}"

# Function to check if PostgreSQL is ready
wait_for_postgres() {
    echo "==> Waiting for PostgreSQL to be ready..."
    
    retries=0
    until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" > /dev/null 2>&1; do
        retries=$((retries + 1))
        
        if [ $retries -ge $MAX_RETRIES ]; then
            echo "ERROR: PostgreSQL did not become ready in time"
            exit 1
        fi
        
        echo "    PostgreSQL is unavailable - sleeping (attempt $retries/$MAX_RETRIES)"
        sleep $RETRY_INTERVAL
    done
    
    echo "==> PostgreSQL is ready!"
}

# Function to check if database exists and is accessible
check_database() {
    echo "==> Checking database connectivity..."
    
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" > /dev/null 2>&1; then
        echo "==> Database is accessible!"
        return 0
    else
        echo "ERROR: Cannot connect to database"
        return 1
    fi
}

# Function to run migrations
run_migrations() {
    echo "==> Running database migrations..."
    
    if [ -f "/app/migrate" ]; then
        /app/migrate up
        echo "==> Migrations completed successfully!"
    else
        echo "WARNING: Migration binary not found, skipping migrations"
    fi
}

# Function to seed initial data (optional)
seed_data() {
    if [ "$SEED_DATA" = "true" ]; then
        echo "==> Seeding initial data..."
        # Add seed data commands here if needed
        echo "==> Data seeding completed!"
    fi
}

# Main execution
main() {
    echo "==> Starting initialization process..."
    echo "    Database Host: $DB_HOST"
    echo "    Database Port: $DB_PORT"
    echo "    Database User: $DB_USER"
    echo "    Database Name: $DB_NAME"
    echo ""
    
    # Wait for PostgreSQL
    wait_for_postgres
    
    # Check database connectivity
    if ! check_database; then
        exit 1
    fi
    
    # Run migrations
    run_migrations
    
    # Seed data if requested
    seed_data
    
    echo ""
    echo "==> Initialization completed successfully!"
    echo "==> Application is ready to start"
}

# Run main function
main