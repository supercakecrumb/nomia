# Quick Start Guide

Get up and running with the Name Affirmation API in 5 minutes.

## Prerequisites

- Docker and Docker Compose installed
- `curl` command-line tool
- The real US SSA data files in `sample-data/real-us-data/`

## Step 1: Start the System

Start all services (API, PostgreSQL, Redis, pgAdmin):

```bash
docker-compose up -d
```

Wait about 10-15 seconds for all services to initialize.

## Step 2: Check Health

Verify the API is running:

```bash
curl http://localhost:8080/health
```

**Expected output:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Step 3: Create US Country

Before uploading data, create the US country entry:

```bash
curl -X POST http://localhost:8080/v1/countries \
  -H "Content-Type: application/json" \
  -d '{
    "code": "US",
    "name": "United States"
  }'
```

**Expected output:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "code": "US",
  "name": "United States",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Save the `id` value** - you'll need it for uploading data!

## Step 4: Upload Data Files

Use the bulk upload script to load all 141 years of US SSA data (1880-2020):

```bash
# Make the script executable (first time only)
chmod +x scripts/bulk-upload.sh

# Upload all files (replace with your actual country_id from Step 3)
./scripts/bulk-upload.sh 550e8400-e29b-41d4-a716-446655440000
```

The script will:
- Upload all 141 files from `sample-data/real-us-data/`
- Show progress for each file
- Wait 2 seconds between uploads to avoid overwhelming the system
- Print a summary when complete

**This will take about 5-10 minutes** depending on your system.

### Manual Upload (Alternative)

If you prefer to upload files manually:

```bash
curl -X POST http://localhost:8080/v1/datasets/upload \
  -F "file=@sample-data/real-us-data/yob2020.txt" \
  -F "country_id=550e8400-e29b-41d4-a716-446655440000"
```

## Step 5: Query the Data

### Check Upload Status

Monitor job progress:

```bash
curl http://localhost:8080/v1/jobs
```

### Search for Names

Find names starting with "Emma":

```bash
curl "http://localhost:8080/v1/names/search?query=Emma&limit=10"
```

**Example output:**
```json
{
  "names": [
    {
      "id": "...",
      "name": "Emma",
      "gender": "F",
      "year": 2020,
      "count": 15581,
      "country_code": "US"
    },
    {
      "id": "...",
      "name": "Emmanuel",
      "gender": "M",
      "year": 2020,
      "count": 2847,
      "country_code": "US"
    }
  ],
  "total": 2,
  "page": 1,
  "page_size": 10
}
```

### Get Top Names by Year

Most popular names in 2020:

```bash
curl "http://localhost:8080/v1/names/top?year=2020&limit=10"
```

### Get Name Statistics

Statistics for a specific name:

```bash
curl "http://localhost:8080/v1/names/stats?name=Emma&country_code=US"
```

### Filter by Gender and Year Range

Female names from 2010-2020:

```bash
curl "http://localhost:8080/v1/names/search?gender=F&year_from=2010&year_to=2020&limit=20"
```

## Step 6: View in pgAdmin (Optional)

Access the database GUI at http://localhost:5050

**Login credentials:**
- Email: `admin@affirm-name.local`
- Password: `admin`

**Connect to database:**
1. Right-click "Servers" → "Register" → "Server"
2. General tab: Name = "Affirm Name DB"
3. Connection tab:
   - Host: `postgres`
   - Port: `5432`
   - Database: `affirm_name`
   - Username: `affirm_user`
   - Password: `affirm_pass`

## Important Notes

### This is a REST API Backend

**There is no web UI yet.** This system provides a REST API for:
- Uploading name datasets
- Querying name statistics
- Searching names by various criteria

You interact with it using:
- `curl` commands (as shown above)
- API clients like Postman or Insomnia
- Your own frontend application

### API Documentation

Full API documentation is available in [`docs/api-specification.md`](docs/api-specification.md)

## Troubleshooting

### Services won't start

```bash
# Check logs
docker-compose logs api

# Restart services
docker-compose down
docker-compose up -d
```

### Upload fails with "country not found"

Make sure you created the country in Step 3 and are using the correct `country_id`.

### Jobs stuck in "pending" status

Check worker logs:

```bash
docker-compose logs api | grep worker
```

The worker pool processes jobs asynchronously. Large files may take a few minutes.

### Database connection errors

Ensure PostgreSQL is running:

```bash
docker-compose ps postgres
```

### Port conflicts

If ports 8080, 5432, 6379, or 5050 are already in use:

1. Edit `docker-compose.yml`
2. Change the port mappings (e.g., `8080:8080` → `8081:8080`)
3. Restart: `docker-compose down && docker-compose up -d`

## Next Steps

- Read the [API Specification](docs/api-specification.md) for all available endpoints
- Check [Architecture](docs/architecture.md) to understand the system design
- Review [Database Schema](docs/database-schema.md) for data structure
- See [TESTING.md](docs/TESTING.md) for testing guidelines

## Stopping the System

```bash
# Stop services (keeps data)
docker-compose down

# Stop and remove all data
docker-compose down -v
```

## Quick Reference

| Service | URL | Purpose |
|---------|-----|---------|
| API | http://localhost:8080 | REST API endpoints |
| Health Check | http://localhost:8080/health | Service status |
| pgAdmin | http://localhost:5050 | Database GUI |
| PostgreSQL | localhost:5432 | Database (internal) |
| Redis | localhost:6379 | Cache (internal) |

## Example Workflow

```bash
# 1. Start system
docker-compose up -d

# 2. Create country and save the ID
COUNTRY_ID=$(curl -X POST http://localhost:8080/v1/countries \
  -H "Content-Type: application/json" \
  -d '{"code":"US","name":"United States"}' \
  | jq -r '.id')

# 3. Upload all data
./scripts/bulk-upload.sh $COUNTRY_ID

# 4. Query names
curl "http://localhost:8080/v1/names/search?query=Emma&limit=5"

# 5. Check statistics
curl "http://localhost:8080/v1/names/stats?name=Emma&country_code=US"
```

Happy querying! 🎉