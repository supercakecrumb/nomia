# Testing Guide for Phase 2 - Country Management API

## Prerequisites

1. PostgreSQL database running
2. Environment variables configured (copy `.env.example` to `.env` and update)
3. Go 1.21+ installed

## Setup

### 1. Start PostgreSQL

```bash
# Using Docker
docker run --name affirm-postgres \
  -e POSTGRES_USER=affirm \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=affirm_name \
  -p 5432:5432 \
  -d postgres:15-alpine
```

### 2. Configure Environment

```bash
cp .env.example .env
# Edit .env with your database credentials
```

Example `.env`:
```
DATABASE_URL=postgres://affirm:secret@localhost:5432/affirm_name?sslmode=disable
SERVER_PORT=8080
LOG_LEVEL=debug
LOG_FORMAT=json
```

### 3. Run Migrations

```bash
go run cmd/migrate/main.go up
```

Expected output:
```
Migrations applied successfully
```

### 4. Start API Server

```bash
go run cmd/api/main.go
```

Expected output:
```json
{"time":"2024-01-15T10:00:00Z","level":"INFO","msg":"Starting API server","version":"1.0.0","port":"8080"}
{"time":"2024-01-15T10:00:00Z","level":"INFO","msg":"Connected to database"}
{"time":"2024-01-15T10:00:00Z","level":"INFO","msg":"Starting HTTP server","port":"8080"}
```

## API Testing

### Health Checks

#### Health Check
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "checks": {
    "database": "connected"
  }
}
```

#### Readiness Check
```bash
curl http://localhost:8080/ready
```

Expected response:
```json
{
  "ready": true
}
```

### Country CRUD Operations

#### 1. Create Country

```bash
curl -X POST http://localhost:8080/v1/countries \
  -H "Content-Type: application/json" \
  -d '{
    "code": "US",
    "name": "United States",
    "source_url": "https://www.ssa.gov/oact/babynames/",
    "attribution": "Social Security Administration"
  }'
```

Expected response (201 Created):
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "US",
    "name": "United States",
    "source_url": "https://www.ssa.gov/oact/babynames/",
    "attribution": "Social Security Administration",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

#### 2. List Countries

```bash
curl http://localhost:8080/v1/countries
```

Expected response (200 OK):
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "code": "US",
      "name": "United States",
      "source_url": "https://www.ssa.gov/oact/babynames/",
      "attribution": "Social Security Administration",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "meta": {
    "total": 4,
    "limit": 100,
    "offset": 0,
    "has_more": false
  }
}
```

#### 3. Get Country by ID

```bash
curl http://localhost:8080/v1/countries/550e8400-e29b-41d4-a716-446655440000
```

Expected response (200 OK):
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "US",
    "name": "United States",
    "source_url": "https://www.ssa.gov/oact/babynames/",
    "attribution": "Social Security Administration",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z",
    "stats": {
      "dataset_count": 0,
      "total_names": 0,
      "year_range": {
        "min": 0,
        "max": 0
      }
    }
  }
}
```

#### 4. Get Country by Code

```bash
curl http://localhost:8080/v1/countries/US
```

Same response as above.

#### 5. Update Country

```bash
curl -X PATCH http://localhost:8080/v1/countries/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "United States of America",
    "attribution": "U.S. Social Security Administration"
  }'
```

Expected response (200 OK):
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "US",
    "name": "United States of America",
    "source_url": "https://www.ssa.gov/oact/babynames/",
    "attribution": "U.S. Social Security Administration",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:05:00Z"
  }
}
```

#### 6. Delete Country

```bash
curl -X DELETE http://localhost:8080/v1/countries/550e8400-e29b-41d4-a716-446655440000
```

Expected response (204 No Content): Empty body

### Error Cases

#### 1. Create Duplicate Country Code

```bash
curl -X POST http://localhost:8080/v1/countries \
  -H "Content-Type: application/json" \
  -d '{
    "code": "US",
    "name": "United States"
  }'
```

Expected response (409 Conflict):
```json
{
  "error": {
    "code": "conflict",
    "message": "Country code already exists"
  }
}
```

#### 2. Invalid Country Code

```bash
curl -X POST http://localhost:8080/v1/countries \
  -H "Content-Type: application/json" \
  -d '{
    "code": "USA",
    "name": "United States"
  }'
```

Expected response (422 Unprocessable Entity):
```json
{
  "error": {
    "code": "validation_error",
    "message": "invalid input: ..."
  }
}
```

#### 3. Country Not Found

```bash
curl http://localhost:8080/v1/countries/00000000-0000-0000-0000-000000000000
```

Expected response (404 Not Found):
```json
{
  "error": {
    "code": "not_found",
    "message": "Country not found"
  }
}
```

#### 4. Delete Country with Datasets

```bash
# First create a country and dataset, then try to delete
curl -X DELETE http://localhost:8080/v1/countries/550e8400-e29b-41d4-a716-446655440000
```

Expected response (409 Conflict):
```json
{
  "error": {
    "code": "conflict",
    "message": "Cannot delete country with associated datasets"
  }
}
```

## Pagination Testing

```bash
# Get first page
curl "http://localhost:8080/v1/countries?limit=2&offset=0"

# Get second page
curl "http://localhost:8080/v1/countries?limit=2&offset=2"
```

## Verification Checklist

- [ ] Server starts successfully
- [ ] Health check returns healthy status
- [ ] Readiness check returns ready status
- [ ] Can create a country
- [ ] Can list countries with pagination
- [ ] Can get country by ID
- [ ] Can get country by code
- [ ] Can update country
- [ ] Can delete country
- [ ] Duplicate country code returns 409
- [ ] Invalid input returns 422
- [ ] Not found returns 404
- [ ] Cannot delete country with datasets returns 409
- [ ] Request logging works
- [ ] Error responses follow API spec
- [ ] CORS headers are present

## Database Verification

```sql
-- Connect to database
psql -U affirm -d affirm_name

-- Check countries
SELECT * FROM countries;

-- Check migrations
SELECT * FROM schema_migrations;
```

## Cleanup

```bash
# Stop server: Ctrl+C

# Drop database (if needed)
go run cmd/migrate/main.go down

# Stop PostgreSQL container
docker stop affirm-postgres
docker rm affirm-postgres
```

## Next Steps

After Phase 2 is complete and tested:
1. Phase 3: File upload and storage
2. Phase 4: Parser framework
3. Phase 5: Background workers
4. Phase 6: Name query API
5. Phase 7: Trend analysis

## Troubleshooting

### Database Connection Failed
- Check PostgreSQL is running: `docker ps`
- Verify DATABASE_URL in .env
- Check firewall/network settings

### Port Already in Use
- Change SERVER_PORT in .env
- Kill process using port: `lsof -ti:8080 | xargs kill`

### Migration Errors
- Check database exists
- Verify migration files are present
- Run `go run cmd/migrate/main.go down` and retry
## Name Query API Testing (Phase 6)

### List Names

Query names for a specific country and year:

```bash
# Basic query
curl "http://localhost:8080/v1/names?country=US&year=2020"

# With gender filter
curl "http://localhost:8080/v1/names?country=US&year=2020&gender=F"

# With name prefix search
curl "http://localhost:8080/v1/names?country=US&year=2020&name=Em"

# With minimum count filter
curl "http://localhost:8080/v1/names?country=US&year=2020&min_count=1000"

# With sorting by name
curl "http://localhost:8080/v1/names?country=US&year=2020&sort=name:asc"

# With pagination
curl "http://localhost:8080/v1/names?country=US&year=2020&limit=10&offset=0"

# Combined filters
curl "http://localhost:8080/v1/names?country=US&year=2020&gender=F&name=Em&sort=count:desc&limit=10"
```

Expected response:
```json
{
  "data": [
    {
      "name": "Emma",
      "gender": "F",
      "count": 15581,
      "year": 2020,
      "country_code": "US",
      "rank": 1
    }
  ],
  "meta": {
    "total": 32033,
    "limit": 100,
    "offset": 0,
    "has_more": true,
    "filters": {
      "country": "US",
      "year": 2020,
      "gender": "F"
    }
  }
}
```

### Search Names

Search for names across all years:

```bash
# Basic search
curl "http://localhost:8080/v1/names/search?q=Em"

# With country filter
curl "http://localhost:8080/v1/names/search?q=Em&country=US"

# With gender filter
curl "http://localhost:8080/v1/names/search?q=Em&gender=F"
```

Expected response:
```json
{
  "data": [
    {
      "name": "Emma",
      "total_count": 1500000,
      "countries": ["US", "GB", "CA"],
      "year_range": {
        "min": 1970,
        "max": 2020
      },
      "gender_distribution": {
        "M": 0.5,
        "F": 99.5
      }
    }
  ],
  "meta": {
    "total": 2,
    "limit": 100,
    "offset": 0,
    "has_more": false
  }
}
```

### Get Name Details

Get detailed information about a specific name:

```bash
# Get name details
curl "http://localhost:8080/v1/names/Emma"

# With country filter
curl "http://localhost:8080/v1/names/Emma?country=US"
```

Expected response:
```json
{
  "data": {
    "name": "Emma",
    "total_count": 1500000,
    "countries": [
      {
        "code": "US",
        "name": "United States",
        "count": 1200000,
        "year_range": {
          "min": 1970,
          "max": 2020
        }
      }
    ],
    "gender_distribution": {
      "M": {
        "count": 7500,
        "percentage": 0.5
      },
      "F": {
        "count": 1492500,
        "percentage": 99.5
      }
    },
    "popularity_trend": "increasing",
    "peak_year": 2020,
    "peak_count": 15581
  }
}
```

### Name Query Error Cases

#### Missing Required Parameters
```bash
curl "http://localhost:8080/v1/names?year=2020"
```

Response (400 Bad Request):
```json
{
  "error": {
    "code": "validation_error",
    "message": "country parameter is required"
  }
}
```

#### Invalid Year
```bash
curl "http://localhost:8080/v1/names?country=US&year=1800"
```

Response (422 Unprocessable Entity):
```json
{
  "error": {
    "code": "validation_error",
    "message": "year must be between 1970 and 2030"
  }
}
```

#### Name Not Found
```bash
curl "http://localhost:8080/v1/names/NonExistentName"
```

Response (404 Not Found):
```json
{
  "error": {
    "code": "not_found",
    "message": "Name not found"
  }
}
```

---
