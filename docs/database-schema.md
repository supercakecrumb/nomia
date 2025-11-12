
# Database Schema Design

## Overview

This document defines the PostgreSQL database schema for the baby name statistics platform. The schema is designed for:
- Data integrity through foreign keys and constraints
- Query performance through strategic indexing
- Audit trail with soft deletes and timestamps
- Extensibility for future features

---

## Entity Relationship Diagram

```
┌─────────────────┐
│    countries    │
│─────────────────│
│ id (PK)         │
│ code (UNIQUE)   │
│ name            │
│ source_url      │
│ attribution     │
│ created_at      │
│ updated_at      │
└────────┬────────┘
         │
         │ 1:N
         │
┌────────▼────────┐
│    datasets     │
│─────────────────│
│ id (PK)         │
│ country_id (FK) │
│ filename        │
│ file_path       │
│ file_size       │
│ status          │
│ row_count       │
│ error_message   │
│ uploaded_by     │
│ uploaded_at     │
│ processed_at    │
│ deleted_at      │
└────────┬────────┘
         │
         │ 1:N
         │
┌────────▼────────┐
│     names       │
│─────────────────│
│ id (PK)         │
│ dataset_id (FK) │
│ country_id (FK) │
│ year            │
│ name            │
│ gender          │
│ count           │
│ created_at      │
│ deleted_at      │
│ UNIQUE(dataset, │
│   year, name,   │
│   gender)       │
└─────────────────┘

┌─────────────────┐
│      jobs       │
│─────────────────│
│ id (PK)         │
│ dataset_id (FK) │
│ type            │
│ status          │
│ payload         │
│ attempts        │
│ max_attempts    │
│ last_error      │
│ next_retry_at   │
│ locked_at       │
│ locked_by       │
│ created_at      │
│ started_at      │
│ completed_at    │
└─────────────────┘

┌─────────────────┐
│    api_keys     │
│─────────────────│
│ id (PK)         │
│ key_hash        │
│ name            │
│ role            │
│ expires_at      │
│ last_used_at    │
│ created_at      │
│ revoked_at      │
└─────────────────┘
```

---

## Table Definitions

### 1. countries

Stores metadata about countries whose data is being tracked.

```sql
CREATE TABLE countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(2) NOT NULL UNIQUE,  -- ISO 3166-1 alpha-2
    name VARCHAR(100) NOT NULL,
    source_url TEXT,                   -- URL to government data source
    attribution TEXT,                  -- Required attribution text
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_countries_code ON countries(code);

-- Comments
COMMENT ON TABLE countries IS 'Countries with baby name statistics';
COMMENT ON COLUMN countries.code IS 'ISO 3166-1 alpha-2 country code';
COMMENT ON COLUMN countries.source_url IS 'URL to official government data source';
COMMENT ON COLUMN countries.attribution IS 'Required attribution text for data usage';
```

**Sample Data:**
```sql
INSERT INTO countries (code, name, source_url, attribution) VALUES
('US', 'United States', 'https://www.ssa.gov/oact/babynames/', 'Social Security Administration'),
('GB', 'United Kingdom', 'https://www.ons.gov.uk/peoplepopulationandcommunity/birthsdeathsandmarriages/livebirths', 'Office for National Statistics'),
('CA', 'Canada', 'https://www.statcan.gc.ca/', 'Statistics Canada'),
('AU', 'Australia', 'https://www.abs.gov.au/', 'Australian Bureau of Statistics');
```

### 2. datasets

Tracks uploaded files and their processing status.

```sql
CREATE TABLE datasets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    filename VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,           -- Path in storage (local or S3)
    file_size BIGINT NOT NULL,         -- Size in bytes
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    row_count INTEGER,                 -- Number of names parsed
    error_message TEXT,                -- Error details if failed
    uploaded_by VARCHAR(100),          -- User/API key that uploaded
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,          -- When processing completed
    deleted_at TIMESTAMPTZ,            -- Soft delete timestamp
    
    CONSTRAINT chk_status CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'reprocessing')),
    CONSTRAINT chk_file_size CHECK (file_size > 0),
    CONSTRAINT chk_row_count CHECK (row_count IS NULL OR row_count >= 0)
);

-- Indexes
CREATE INDEX idx_datasets_country_id ON datasets(country_id);
CREATE INDEX idx_datasets_status ON datasets(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_datasets_uploaded_at ON datasets(uploaded_at DESC);
CREATE INDEX idx_datasets_deleted_at ON datasets(deleted_at) WHERE deleted_at IS NOT NULL;

-- Comments
COMMENT ON TABLE datasets IS 'Uploaded dataset files and their processing status';
COMMENT ON COLUMN datasets.status IS 'Processing status: pending, processing, completed, failed, reprocessing';
COMMENT ON COLUMN datasets.file_path IS 'Storage path (e.g., uploads/uuid/original.csv or s3://bucket/key)';
COMMENT ON COLUMN datasets.deleted_at IS 'Soft delete timestamp for versioning';
```

**Status Transitions:**
```
pending → processing → completed
pending → processing → failed
completed → reprocessing → completed
completed → reprocessing → failed
```

### 3. names

Core table storing normalized name statistics.

```sql
CREATE TABLE names (
    id BIGSERIAL PRIMARY KEY,
    dataset_id UUID NOT NULL REFERENCES datasets(id) ON DELETE CASCADE,
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    year INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    gender CHAR(1) NOT NULL,
    count INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,            -- Soft delete for reprocessing
    
    CONSTRAINT chk_year CHECK (year BETWEEN 1800 AND 2100),
    CONSTRAINT chk_gender CHECK (gender IN ('M', 'F')),
    CONSTRAINT chk_count CHECK (count > 0),
    CONSTRAINT chk_name_length CHECK (LENGTH(name) > 0 AND LENGTH(name) <= 100),
    
    -- Ensure uniqueness within a dataset
    CONSTRAINT uq_names_dataset_year_name_gender 
        UNIQUE (dataset_id, year, name, gender)
);

-- Indexes for common query patterns
CREATE INDEX idx_names_country_year_gender ON names(country_id, year, gender) 
    WHERE deleted_at IS NULL;

CREATE INDEX idx_names_name_country ON names(name, country_id) 
    WHERE deleted_at IS NULL;

CREATE INDEX idx_names_dataset_id ON names(dataset_id);

CREATE INDEX idx_names_deleted_at ON names(deleted_at) 
    WHERE deleted_at IS NOT NULL;

-- Index for sorting by count (most popular names)
CREATE INDEX idx_names_country_year_count ON names(country_id, year, count DESC) 
    WHERE deleted_at IS NULL;

-- Composite index for trend queries
CREATE INDEX idx_names_name_year_gender ON names(name, year, gender, country_id) 
    WHERE deleted_at IS NULL;

-- Comments
COMMENT ON TABLE names IS 'Normalized baby name statistics';
COMMENT ON COLUMN names.gender IS 'M for Male, F for Female';
COMMENT ON COLUMN names.count IS 'Number of babies with this name in the given year';
COMMENT ON COLUMN names.deleted_at IS 'Soft delete timestamp for reprocessing';
```

**Partitioning Strategy (Future):**

When the table grows beyond 50M rows, consider partitioning:

```sql
-- Partition by country_id (if queries are mostly country-specific)
CREATE TABLE names (
    -- same columns
) PARTITION BY LIST (country_id);

CREATE TABLE names_us PARTITION OF names FOR VALUES IN ('uuid-for-us');
CREATE TABLE names_uk PARTITION OF names FOR VALUES IN ('uuid-for-uk');
-- etc.

-- OR partition by year range (if queries are mostly time-based)
CREATE TABLE names (
    -- same columns
) PARTITION BY RANGE (year);

CREATE TABLE names_1970_1989 PARTITION OF names FOR VALUES FROM (1970) TO (1990);
CREATE TABLE names_1990_2009 PARTITION OF names FOR VALUES FROM (1990) TO (2010);
CREATE TABLE names_2010_2029 PARTITION OF names FOR VALUES FROM (2010) TO (2030);
```

### 4. jobs

Job queue for asynchronous processing.

```sql
CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dataset_id UUID REFERENCES datasets(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'queued',
    payload JSONB,                     -- Job-specific data
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 3,
    last_error TEXT,
    next_retry_at TIMESTAMPTZ,
    locked_at TIMESTAMPTZ,             -- When job was locked by worker
    locked_by VARCHAR(100),            -- Worker ID that locked the job
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    
    CONSTRAINT chk_type CHECK (type IN ('parse_dataset', 'reprocess_dataset')),
    CONSTRAINT chk_status CHECK (status IN ('queued', 'running', 'completed', 'failed')),
    CONSTRAINT chk_attempts CHECK (attempts >= 0 AND attempts <= max_attempts)
);

-- Indexes for job queue operations
CREATE INDEX idx_jobs_status_next_retry ON jobs(status, next_retry_at) 
    WHERE status = 'queued' AND (next_retry_at IS NULL OR next_retry_at <= NOW());

CREATE INDEX idx_jobs_dataset_id ON jobs(dataset_id);

CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC);

CREATE INDEX idx_jobs_locked_at ON jobs(locked_at) 
    WHERE locked_at IS NOT NULL AND status = 'running';

-- Comments
COMMENT ON TABLE jobs IS 'Asynchronous job queue for background processing';
COMMENT ON COLUMN jobs.type IS 'Job type: parse_dataset, reprocess_dataset';
COMMENT ON COLUMN jobs.status IS 'Job status: queued, running, completed, failed';
COMMENT ON COLUMN jobs.locked_at IS 'Timestamp when job was locked by a worker';
COMMENT ON COLUMN jobs.locked_by IS 'Worker ID (hostname or container ID)';
```

**Job Payload Examples:**

```json
// parse_dataset job
{
  "dataset_id": "uuid",
  "country_code": "US",
  "parser_type": "us_ssa"
}

// reprocess_dataset job
{
  "dataset_id": "uuid",
  "country_code": "US",
  "parser_type": "us_ssa",
  "reason": "parser_bug_fix"
}
```

### 5. api_keys

API key management for authentication.

```sql
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key_hash VARCHAR(255) NOT NULL UNIQUE,  -- bcrypt hash of the key
    name VARCHAR(100) NOT NULL,              -- Descriptive name
    role VARCHAR(20) NOT NULL DEFAULT 'viewer',
    expires_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,
    
    CONSTRAINT chk_role CHECK (role IN ('admin', 'viewer'))
);

-- Indexes
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash) 
    WHERE revoked_at IS NULL;

CREATE INDEX idx_api_keys_expires_at ON api_keys(expires_at) 
    WHERE expires_at IS NOT NULL AND revoked_at IS NULL;

-- Comments
COMMENT ON TABLE api_keys IS 'API keys for authentication';
COMMENT ON COLUMN api_keys.key_hash IS 'bcrypt hash of the API key';
COMMENT ON COLUMN api_keys.role IS 'User role: admin (full access), viewer (read-only)';
COMMENT ON COLUMN api_keys.revoked_at IS 'Timestamp when key was revoked';
```

---

## Materialized Views (Optional)

For performance optimization of common aggregations:

### name_statistics

Pre-computed statistics for popular queries.

```sql
CREATE MATERIALIZED VIEW name_statistics AS
SELECT 
    country_id,
    year,
    gender,
    name,
    SUM(count) as total_count,
    COUNT(DISTINCT dataset_id) as dataset_count
FROM names
WHERE deleted_at IS NULL
GROUP BY country_id, year, gender, name;

-- Indexes on materialized view
CREATE INDEX idx_name_stats_country_year ON name_statistics(country_id, year);
CREATE INDEX idx_name_stats_name ON name_statistics(name);

-- Refresh strategy: after each dataset processing
REFRESH MATERIALIZED VIEW CONCURRENTLY name_statistics;
```

### gender_probabilities

Pre-computed gender probabilities for names across all years.

```sql
CREATE MATERIALIZED VIEW gender_probabilities AS
SELECT 
    name,
    country_id,
    SUM(CASE WHEN gender = 'M' THEN count ELSE 0 END) as male_count,
    SUM(CASE WHEN gender = 'F' THEN count ELSE 0 END) as female_count,
    SUM(count) as total_count,
    ROUND(
        SUM(CASE WHEN gender = 'M' THEN count ELSE 0 END)::NUMERIC / 
        NULLIF(SUM(count), 0) * 100, 
        2
    ) as male_probability,
    ROUND(
        SUM(CASE WHEN gender = 'F' THEN count ELSE 0 END)::NUMERIC / 
        NULLIF(SUM(count), 0) * 100, 
        2
    ) as female_probability
FROM names
WHERE deleted_at IS NULL
GROUP BY name, country_id
HAVING SUM(count) >= 100;  -- Minimum threshold for statistical significance

-- Indexes
CREATE INDEX idx_gender_prob_name ON gender_probabilities(name);
CREATE INDEX idx_gender_prob_country ON gender_probabilities(country_id);

-- Refresh after dataset processing
REFRESH MATERIALIZED VIEW CONCURRENTLY gender_probabilities;
```

---

## Functions and Triggers

### 1. Update Timestamp Trigger

Automatically update `updated_at` column.

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_countries_updated_at
    BEFORE UPDATE ON countries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

### 2. Job Lock Function

Atomically lock a job for processing.

```sql
CREATE OR REPLACE FUNCTION lock_next_job(worker_id VARCHAR(100))
RETURNS TABLE (
    job_id UUID,
    job_type VARCHAR(50),
    job_payload JSONB
) AS $$
BEGIN
    RETURN QUERY
    UPDATE jobs
    SET 
        status = 'running',
        locked_at = NOW(),
        locked_by = worker_id,
        started_at = COALESCE(started_at, NOW()),
        attempts = attempts + 1
    WHERE id = (
        SELECT id
        FROM jobs
        WHERE status = 'queued'
          AND (next_retry_at IS NULL OR next_retry_at <= NOW())
        ORDER BY created_at ASC
        LIMIT 1
        FOR UPDATE SKIP LOCKED
    )
    RETURNING id, type, payload;
END;
$$ LANGUAGE plpgsql;
```

### 3. Cleanup Old Jobs Function

Remove completed jobs older than 30 days.

```sql
CREATE OR REPLACE FUNCTION cleanup_old_jobs()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM jobs
    WHERE status IN ('completed', 'failed')
      AND completed_at < NOW() - INTERVAL '30 days';
    
    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Schedule with pg_cron (if available) or run manually
-- SELECT cron.schedule('cleanup-jobs', '0 2 * * *', 'SELECT cleanup_old_jobs()');
```

---

## Migration Files

### Migration: 001_initial_schema.up.sql

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create countries table
CREATE TABLE countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(2) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    source_url TEXT,
    attribution TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_countries_code ON countries(code);

-- Create datasets table
CREATE TABLE datasets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    filename VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    row_count INTEGER,
    error_message TEXT,
    uploaded_by VARCHAR(100),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT chk_status CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'reprocessing')),
    CONSTRAINT chk_file_size CHECK (file_size > 0),
    CONSTRAINT chk_row_count CHECK (row_count IS NULL OR row_count >= 0)
);

CREATE INDEX idx_datasets_country_id ON datasets(country_id);
CREATE INDEX idx_datasets_status ON datasets(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_datasets_uploaded_at ON datasets(uploaded_at DESC);
CREATE INDEX idx_datasets_deleted_at ON datasets(deleted_at) WHERE deleted_at IS NOT NULL;

-- Create names table
CREATE TABLE names (
    id BIGSERIAL PRIMARY KEY,
    dataset_id UUID NOT NULL REFERENCES datasets(id) ON DELETE CASCADE,
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    year INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    gender CHAR(1) NOT NULL,
    count INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT chk_year CHECK (year BETWEEN 1800 AND 2100),
    CONSTRAINT chk_gender CHECK (gender IN ('M', 'F')),
    CONSTRAINT chk_count CHECK (count > 0),
    CONSTRAINT chk_name_length CHECK (LENGTH(name) > 0 AND LENGTH(name) <= 100),
    CONSTRAINT uq_names_dataset_year_name_gender UNIQUE (dataset_id, year, name, gender)
);

CREATE INDEX idx_names_country_year_gender ON names(country_id, year, gender) WHERE deleted_at IS NULL;
CREATE INDEX idx_names_name_country ON names(name, country_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_names_dataset_id ON names(dataset_id);
CREATE INDEX idx_names_deleted_at ON names(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX idx_names_country_year_count ON names(country_id, year, count DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_names_name_year_gender ON names(name, year, gender, country_id) WHERE deleted_at IS NULL;

-- Create jobs table
CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dataset_id UUID REFERENCES datasets(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'queued',
    payload JSONB,
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 3,
    last_error TEXT,
    next_retry_at TIMESTAMPTZ,
    locked_at TIMESTAMPTZ,
    locked_by VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    
    CONSTRAINT chk_type CHECK (type IN ('parse_dataset', 'reprocess_dataset')),
    CONSTRAINT chk_status CHECK (status IN ('queued', 'running', 'completed', 'failed')),
    CONSTRAINT chk_attempts CHECK (attempts >= 0 AND attempts <= max_attempts)
);

CREATE INDEX idx_jobs_status_next_retry ON jobs(status, next_retry_at) 
    WHERE status = 'queued' AND (next_retry_at IS NULL OR next_retry_at <= NOW());
CREATE INDEX idx_jobs_dataset_id ON jobs(dataset_id);
CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC);
CREATE INDEX idx_jobs_locked_at ON jobs(locked_at) WHERE locked_at IS NOT NULL AND status = 'running';

-- Create api_keys table
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key_hash VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'viewer',
    expires_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ,
    
    CONSTRAINT chk_role CHECK (role IN ('admin', 'viewer'))
);

CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash) WHERE revoked_at IS NULL;
CREATE INDEX idx_api_keys_expires_at ON api_keys(expires_at) 
    WHERE expires_at IS NOT NULL AND revoked_at IS NULL;

-- Create update timestamp trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_countries_updated_at
    BEFORE UPDATE ON countries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Create job lock function
CREATE OR REPLACE FUNCTION lock_next_job(worker_id VARCHAR(100))
RETURNS TABLE (
    job_id UUID,
    job_type VARCHAR(50),
    job_payload JSONB
) AS $$
BEGIN
    RETURN QUERY
    UPDATE jobs
    SET 
        status = 'running',
        locked_at = NOW(),
        locked_by = worker_id,
        started_at = COALESCE(started_at, NOW()),
        attempts = attempts + 1
    WHERE id = (
        SELECT id
        FROM jobs
        WHERE status = 'queued'
          AND (next_retry_at IS NULL OR next_retry_at <= NOW())
        ORDER BY created_at ASC
        LIMIT 1
        FOR UPDATE SKIP LOCKED
    )
    RETURNING id, type, payload;
END;
$$ LANGUAGE plpgsql;

-- Insert sample countries
INSERT INTO countries (code, name, source_url, attribution) VALUES
('US', 'United States', 'https://www.ssa.gov/oact/babynames/', 'Social Security Administration'),
('GB', 'United Kingdom', 'https://www.ons.gov.uk/peoplepopulationandcommunity/birthsdeathsandmarriages/livebirths', 'Office for National Statistics'),
('CA', 'Canada', 'https://www.statcan.gc.ca/', 'Statistics Canada'),
('AU', 'Australia', 'https://www.abs.gov.au/', 'Australian Bureau of Statistics');
```

### Migration: 001_initial_schema.down.sql

```sql
-- Drop tables in reverse order
DROP TABLE IF EXISTS api_keys CASCADE;
DROP TABLE IF EXISTS jobs CASCADE;
DROP TABLE IF EXISTS names CASCADE;
DROP TABLE IF EXISTS datasets CASCADE;
DROP TABLE IF EXISTS countries CASCADE;

-- Drop functions
DROP FUNCTION IF EXISTS lock_next_job(VARCHAR);
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop extension
DROP EXTENSION IF EXISTS "pgcrypto";
```

---

## Index Strategy

### Query Pattern Analysis

**Most Common Queries:**

1. **List names by country/year/gender**
   ```sql
   SELECT * FROM names 
   WHERE country_id = ? AND year = ? AND gender = ?
   ORDER BY count DESC
   LIMIT 100;
   ```
   **Index:** `idx_names_country_year_gender`, `idx_names_country_year_count`

2. **Search names by prefix**
   ```sql
   SELECT * FROM names 
   WHERE country_id = ? AND name ILIKE 'Em%'
   ORDER BY count DESC;
   ```
   **Index:** Consider `pg_trgm` extension for fuzzy search

3. **Trend analysis for a name**
   ```sql
   SELECT year, gender, SUM(count) 
   FROM names 
   WHERE name = ? AND country_id = ?
   GROUP BY year, gender
   ORDER BY year;
   ```
   **Index:** `idx_names_name_year_gender`

4. **Gender probability**
   ```sql
   SELECT gender, SUM(count) 
   FROM names 
   WHERE name = ?
   GROUP BY gender;
   ```
   **Index:** `idx_names_name_country` or materialized view

### Index Maintenance

```sql
-- Analyze tables after bulk inserts
ANALYZE names;
ANALYZE datasets;

-- Reindex if fragmented
REINDEX TABLE names;

-- Monitor index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan ASC;

-- Remove unused indexes
-- (Indexes with idx_scan = 0 after significant usage period)
```

---

## Performance Considerations

### 1. Bulk Insert Optimization

```sql
-- Disable indexes during bulk load
DROP INDEX idx_names_country_year_gender;
DROP INDEX idx_names_name_country;
-- ... drop other indexes

-- Use COPY for bulk insert (10x faster than INSERT)
COPY names (dataset_id, country_id, year, name, gender, count)
FROM '/path/to/data.csv'
WITH (FORMAT csv, HEADER true);

-- Recreate indexes
CREATE INDEX idx_names_country_year_gender ON names(country_id, year, gender) 
    WHERE deleted_at IS NULL;
-- ... recreate other indexes

-- Analyze table
ANALYZE names;
```

### 2. Query Optimization

```sql
-- Use EXPLAIN ANALYZE to identify slow queries
EXPLAIN (ANALYZE, BUFFERS) 
SELECT * FROM names 
WHERE country_id = 'uuid' AND year = 2020
ORDER BY count DESC
LIMIT 100;

-- Monitor slow queries
ALTER DATABASE affirm_name SET log_min_duration_statement = 1000;  -- Log queries >1s
```

### 3. Connection Pooling

```go
// Configure connection pool
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)
db.SetConnMaxIdleTime(time.Minute * 10)
```

### 4. Vacuum and Maintenance

```sql
-- Regular vacuum to reclaim space
VACUUM ANALYZE names;

-- Auto-vacuum configuration
ALTER TABLE names SET (
    autovacuum_vacuum_scale_factor = 0.1,
    autovacuum_analyze_scale_factor = 0.05
);
```

---

## Data Integrity

### Constraints Summary

| Table | Constraint | Purpose |
|-------|-----------|---------|
| countries | UNIQUE(code) | Prevent duplicate country codes |
| datasets | FK(country_id) | Ensure country exists |
| datasets | CHECK(status) | Valid status values only |
| names | FK(dataset_id) | Ensure dataset exists |
| names | FK(country_id) | Ensure country exists |
| names | UNIQUE(dataset, year, name, gender) | No duplicates within dataset |
| names | CHECK(gender) | Only M or F allowed |
| names | CHECK(year) | Reasonable year range |
| jobs | FK(dataset_id) | Ensure dataset exists |
| jobs | CHECK(status) | Valid status values only |

### Referential Integrity

- **ON DELETE RESTRICT**: Prevents deletion of countries with datasets
- **ON DELETE CASCADE**: Automatically deletes names when dataset is deleted
- **Soft Deletes**: Use `deleted_at` for versioning, not physical deletion

---

## Backup Strategy

### Daily Backups

```bash
# Full database backup
pg_dump -Fc affirm_name > backup_$(date +%Y%m%d).dump

# Restore
pg_restore -d affirm_name backup_20240101.dump
```

### Continuous Archiving (WAL)

```sql
-- Enable WAL archiving
ALTER SYSTEM SET wal_level = replica;
ALTER SYSTEM SET archive_mode = on;
ALTER SYSTEM SET archive_command = 'cp %p /archive/%f';

-- Point-in-time recovery
pg_basebackup -D /backup/base -Ft -z -P
```

### Backup Retention

- Daily backups: 30 days
- Weekly backups: 12 weeks
- Monthly backups: 12 months

---

## Monitoring Queries

### Database Size

```sql
SELECT 
    pg_size_pretty(pg_database_size('affirm_name')) as database_size;
```

### Table Sizes

```sql
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

### Index Usage

```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan ASC;
```

### Active Connections

```sql
SELECT 
    datname,
    count(*) as connections,
    max(state) as state
FROM pg_stat_activity
WHERE datname = 'affirm_name'
GROUP BY datname;
```

### Long Running Queries

```sql
SELECT 
    pid,
    now() - query_start as duration,
    state,
    query
FROM pg_stat_activity
WHERE state != 'idle'
  AND now() - query_start > interval '1 minute'
ORDER BY query_start;
```

---

## Security Considerations

### 1. Row-Level Security (Future)

For multi-tenancy support:

```sql
-- Enable RLS on names table
ALTER TABLE names ENABLE ROW LEVEL SECURITY;

-- Policy: Users can only see their organization's data
CREATE POLICY names_org_isolation ON names
    FOR SELECT
    USING (country_id IN (
        SELECT country_id FROM user_organizations 
        WHERE user_id = current_user_id()
    ));
```

### 2. Encryption

- **At Rest**: Enable PostgreSQL TDE (Transparent Data Encryption)
- **In Transit**: Require SSL connections
- **Sensitive Data**: Hash API keys with bcrypt

```sql
-- Require SSL connections
ALTER SYSTEM SET ssl = on;
ALTER SYSTEM SET ssl_cert_file = '/path/to/server.crt';
ALTER SYSTEM SET ssl_key_file = '/path/to/server.key';
```

### 3. Access Control

```sql
-- Create read-only role for API
CREATE ROLE api_readonly;
GRANT CONNECT ON DATABASE affirm_name TO api_readonly;
GRANT USAGE ON SCHEMA public TO api_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO api_readonly;

-- Create read-write role for workers
CREATE ROLE api_worker;
GRANT CONNECT ON DATABASE affirm_name TO api_worker;
GRANT USAGE ON SCHEMA public TO api_worker;
GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO api_worker;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO api_worker;
```

---

## Appendix: Sample Queries

### Most Popular Names by Year

```sql
SELECT 
    c.name as country,
    n.year,
    n.gender,
    n.name,
    n.count,
    ROW_NUMBER() OVER (PARTITION BY n.country_id, n.year, n.gender ORDER BY n.count DESC) as rank
FROM names n
JOIN countries c ON c.id = n.country_id
WHERE n.country_id = 'uuid-for-us'
  AND n.year = 2020
  AND n.deleted_at IS NULL
ORDER BY n.gender, n.count DESC
LIMIT 10;
```

### Name Trend Over Time

```sql
SELECT 
    year,
    gender,
    SUM(count) as total_count,
    RANK() OVER (PARTITION BY year ORDER BY SUM(count) DESC) as rank_in_year
FROM names
WHERE name = 'Emma'
  AND country_id = 'uuid-for-us'
  AND deleted_at IS NULL
GROUP BY year, gender
ORDER BY year, gender;
```

### Gender Distribution for a Name

```sql
SELECT 
    name,
    SUM(CASE WHEN gender = 'M' THEN count ELSE 0 END) as male_count,
    SUM(CASE WHEN gender = 'F' THEN count ELSE 0 END) as female_count,
    ROUND(
        100.0 * SUM(CASE WHEN gender = 'M' THEN count ELSE 0 END) / 
        NULLIF(SUM(count), 0),
        2
    ) as male_percentage
FROM names
WHERE name = 'Jordan'
  AND country_id = 'uuid-for-us'
  AND deleted_at IS NULL
GROUP BY name;
```

### Dataset Processing Statistics

```sql
SELECT 
    c.name as country,
    d.status,
    COUNT(*) as dataset_count,
    SUM(d.row_count) as total_rows,
    AVG(EXTRACT(EPOCH FROM (d.processed_at - d.uploaded_at))) as avg_processing_seconds
FROM datasets d
JOIN countries c ON c.id = d.country_id
WHERE d.deleted_at IS NULL
GROUP BY c.name, d.status
ORDER BY c.name, d.status;
```

### Job Queue Health

```sql
SELECT 
    status,
    COUNT(*) as job_count,
    AVG(attempts) as avg_attempts,
    MAX(created_at) as latest_job,
    MIN(CASE WHEN status = 'queued' THEN created_at END) as oldest_queued
FROM jobs
GROUP BY status
ORDER BY 
    CASE status 
        WHEN 'queued' THEN 1
        WHEN 'running' THEN 2
        WHEN 'failed' THEN 3
        WHEN 'completed' THEN 4
    END;
```

---

## Migration Management

### Using golang-migrate

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create new migration
migrate create -ext sql -dir migrations -seq add_user_table

# Run migrations
migrate -database "postgres://user:pass@localhost:5432/affirm_name?sslmode=disable" \
        -path migrations up

# Rollback last migration
migrate -database "postgres://user:pass@localhost:5432/affirm_name?sslmode=disable" \
        -path migrations down 1

# Check migration version
migrate -database "postgres://user:pass@localhost:5432/affirm_name?sslmode=disable" \
        -path migrations version
```

### Migration Best Practices

1. **Always test migrations on a copy of production data**
2. **Make migrations reversible (provide .down.sql)**
3. **Keep migrations small and focused**
4. **Never modify existing migrations after deployment**
5. **Use transactions for DDL operations**
6. **Test rollback procedures**

### Migration Template

```sql
-- migrations/XXX_description.up.sql
BEGIN;

-- Your changes here
ALTER TABLE names ADD COLUMN new_field VARCHAR(100);

COMMIT;
```

```sql
-- migrations/XXX_description.down.sql
BEGIN;

-- Reverse your changes
ALTER TABLE names DROP COLUMN new_field;

COMMIT;
```

---

## Conclusion

This database schema provides:

✅ **Data Integrity**: Foreign keys, constraints, and validation
✅ **Performance**: Strategic indexes for common query patterns
✅ **Scalability**: Partitioning strategy for growth
✅ **Auditability**: Soft deletes and timestamps
✅ **Flexibility**: Extensible for future features
✅ **Reliability**: Transaction support and backup strategy

The schema is production-ready and can handle millions of records efficiently while maintaining data quality and consistency.