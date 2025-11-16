#!/bin/bash
set -e

echo "🚀 Importing US name data..."

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

# Ensure database is running
echo "📊 Checking database status..."
if ! docker-compose ps | grep -q "affirm-name-db.*Up"; then
    echo "🔄 Starting database..."
    docker-compose up -d
    echo "⏳ Waiting for database to be ready..."
    sleep 5
fi

# Verify database is accessible
echo "🔍 Verifying database connection..."
if ! docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo "❌ Database is not ready. Please check docker-compose logs"
    exit 1
fi

echo "✅ Database is ready"

# Set DATABASE_URL environment variable
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/affirm_name?sslmode=disable"

# Run the import tool
echo ""
echo "📥 Starting import process..."
cd backend
go run cmd/import/main.go

echo ""
echo "✅ Import complete!"
echo ""
echo "📈 Verifying imported data..."
docker-compose exec -T postgres psql -U postgres -d affirm_name -c "
SELECT 
    year, 
    COUNT(*) as total_names,
    COUNT(DISTINCT name) as unique_names,
    SUM(count) as total_occurrences
FROM names 
GROUP BY year 
ORDER BY year;
"

echo ""
echo "🎉 Done! US name data has been successfully imported."