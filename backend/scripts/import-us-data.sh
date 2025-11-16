#!/bin/bash
set -e

echo "========================================="
echo "Affirm Name - US Data Import"
echo "========================================="
echo ""

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

# Ensure database is running
echo "1. Checking database status..."
if ! docker-compose ps | grep -q "affirm-name-db.*Up"; then
    echo "   Starting database..."
    docker-compose up -d
    echo "   Waiting for database to be ready..."
    sleep 5
fi

# Check database connection
echo "2. Verifying database connection..."
if ! docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo "   ❌ ERROR: Database not ready"
    exit 1
fi
echo "   ✅ Database is ready"

# Set DATABASE_URL environment variable
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/affirm_name?sslmode=disable"

# Run the import tool
echo ""
echo "3. Importing US name data..."
cd backend

# Determine import parameters based on arguments
IMPORT_ARGS="-country=US -dir=../names-example"

if [ "$1" = "all" ]; then
    echo "   Mode: Importing ALL available years"
    # No year filtering - import all files found
elif [ -n "$1" ] && [ -n "$2" ]; then
    echo "   Mode: Importing years $1 to $2"
    IMPORT_ARGS="$IMPORT_ARGS -year-from=$1 -year-to=$2"
elif [ -n "$1" ]; then
    echo "   Mode: Importing single year $1"
    IMPORT_ARGS="$IMPORT_ARGS -year-from=$1 -year-to=$1"
else
    echo "   Mode: Importing available years in names-example/"
fi

# Check if verbose flag is passed
if [ "$3" = "-v" ] || [ "$3" = "--verbose" ]; then
    IMPORT_ARGS="$IMPORT_ARGS -verbose"
fi

# Run import
echo ""
go run cmd/import/main.go $IMPORT_ARGS

# Show statistics
echo ""
echo "4. Import complete! Database statistics:"
echo ""
docker-compose exec -T postgres psql -U postgres -d affirm_name -c "
SELECT 
    c.name as country,
    COUNT(DISTINCT n.year) as years,
    MIN(n.year) as first_year,
    MAX(n.year) as last_year,
    COUNT(DISTINCT n.name) as unique_names,
    COUNT(*) as total_records,
    SUM(n.count) as total_occurrences
FROM names n
JOIN countries c ON n.country_id = c.id
GROUP BY c.name
ORDER BY c.name;
"

echo ""
echo "Year-by-year breakdown:"
docker-compose exec -T postgres psql -U postgres -d affirm_name -c "
SELECT 
    n.year,
    COUNT(*) as total_names,
    COUNT(DISTINCT n.name) as unique_names,
    SUM(n.count) as total_occurrences
FROM names n
JOIN countries c ON n.country_id = c.id
WHERE c.code = 'US'
GROUP BY n.year 
ORDER BY n.year DESC
LIMIT 10;
"

echo ""
echo "========================================="
echo "✅ Import complete!"
echo "========================================="
echo ""
echo "Usage examples:"
echo "  Import all years:    bash backend/scripts/import-us-data.sh all"
echo "  Import range:        bash backend/scripts/import-us-data.sh 2020 2024"
echo "  Import single year:  bash backend/scripts/import-us-data.sh 2023"
echo "  Verbose output:      bash backend/scripts/import-us-data.sh all -v"
echo ""