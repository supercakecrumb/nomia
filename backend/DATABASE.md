# Database Setup Guide

## Prerequisites
- Docker and Docker Compose installed

## Quick Start

### 1. Start PostgreSQL
```bash
docker-compose up -d
```

### 2. Verify Database is Running
```bash
docker-compose ps
docker-compose logs postgres
```

### 3. Access Database
```bash
# Using psql from host (if installed)
psql -h localhost -U postgres -d affirm_name

# Using psql from container
docker-compose exec postgres psql -U postgres -d affirm_name
```

### 4. Stop Database
```bash
docker-compose down
```

### 5. Stop and Remove Data (CAUTION)
```bash
docker-compose down -v
```

## Database Connection
- **Host**: localhost
- **Port**: 5432
- **Database**: affirm_name
- **User**: postgres
- **Password**: postgres
- **Connection String**: `postgresql://postgres:postgres@localhost:5432/affirm_name?sslmode=disable`

## Data Persistence
- Data is stored in Docker volume `postgres_data`
- Volume persists even when container is removed
- To backup data, use `docker-compose exec postgres pg_dump`

## Migrations
- Migrations in `migrations/` folder are automatically run on first startup
- Files are mounted read-only to `/docker-entrypoint-initdb.d`
- Only runs on fresh database (empty data directory)

## Troubleshooting

### Port Already in Use
```bash
# Find process using port 5432
lsof -i :5432
# Stop existing PostgreSQL or change port in docker-compose.yml
```

### Reset Database
```bash
docker-compose down -v  # Removes volume
docker-compose up -d     # Creates fresh database
```

### View Logs
```bash
docker-compose logs -f postgres