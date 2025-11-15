-- Affirm Name - Initial Schema Migration
-- Version: 001
-- Description: Create core tables with indexes for name data storage and querying
-- Author: Architecture Team
-- Date: 2025-11-15

-- Enable required PostgreSQL extensions
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- ============================================================================
-- Table: countries
-- Purpose: Store country metadata and data source information
-- ============================================================================

CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    data_source_name VARCHAR(255) NOT NULL,
    data_source_url TEXT NOT NULL,
    data_source_description TEXT,
    data_source_requires_manual_download BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- Index for fast lookups by country code
CREATE UNIQUE INDEX idx_countries_code ON countries(code);

COMMENT ON TABLE countries IS 'Stores metadata about countries and their data sources';
COMMENT ON COLUMN countries.code IS 'ISO-like country code (e.g., US, UK, SE)';
COMMENT ON COLUMN countries.data_source_name IS 'Name of statistical agency or source';
COMMENT ON COLUMN countries.data_source_url IS 'Canonical URL for the data source';
COMMENT ON COLUMN countries.data_source_requires_manual_download IS 'Whether datasets must be manually downloaded';

-- ============================================================================
-- Table: name_datasets
-- Purpose: Track uploaded dataset files and their parsing status
-- ============================================================================

CREATE TABLE name_datasets (
    id SERIAL PRIMARY KEY,
    country_id INTEGER NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    source_file_name VARCHAR(255) NOT NULL,
    source_url TEXT,
    year_from INTEGER,
    year_to INTEGER,
    file_type VARCHAR(50) NOT NULL,
    storage_path TEXT NOT NULL,
    checksum VARCHAR(71), -- 'sha256:' + 64 hex chars
    parser_version VARCHAR(50),
    parse_status VARCHAR(50) NOT NULL DEFAULT 'uploaded' CHECK (parse_status IN ('uploaded', 'parsing', 'parsed', 'failed')),
    uploaded_at TIMESTAMP DEFAULT NOW() NOT NULL,
    uploaded_by VARCHAR(255), -- User ID from JWT
    parsed_at TIMESTAMP,
    error_message TEXT,
    CONSTRAINT chk_year_range CHECK (year_from IS NULL OR year_to IS NULL OR year_from <= year_to)
);

-- Index for querying datasets by country
CREATE INDEX idx_datasets_country ON name_datasets(country_id);

-- Index for filtering by parse status (e.g., finding failed uploads)
CREATE INDEX idx_datasets_status ON name_datasets(parse_status);

-- Index for duplicate detection (checksum-based)
CREATE INDEX idx_datasets_checksum ON name_datasets(checksum) WHERE checksum IS NOT NULL;

COMMENT ON TABLE name_datasets IS 'Represents uploaded dataset files and their processing status';
COMMENT ON COLUMN name_datasets.checksum IS 'SHA-256 hash for deduplication and integrity checking';
COMMENT ON COLUMN name_datasets.parse_status IS 'Current status: uploaded, parsing, parsed, or failed';
COMMENT ON COLUMN name_datasets.uploaded_by IS 'User ID who uploaded the dataset (from JWT)';

-- ============================================================================
-- Table: names
-- Purpose: Core fact table storing atomic name records by year, country, gender
-- ============================================================================

CREATE TABLE names (
    id BIGSERIAL PRIMARY KEY,
    country_id INTEGER NOT NULL REFERENCES countries(id) ON DELETE RESTRICT,
    dataset_id INTEGER NOT NULL REFERENCES name_datasets(id) ON DELETE RESTRICT,
    year INTEGER NOT NULL CHECK (year >= 1800 AND year <= 2100),
    name VARCHAR(255) NOT NULL,
    gender CHAR(1) NOT NULL CHECK (gender IN ('M', 'F', 'U')),
    count INTEGER NOT NULL CHECK (count > 0)
);

-- Critical composite index for filter pipeline (Stage 1: year + country + name filtering)
-- Column order optimized for query pattern: filter by country_id, year, then name
CREATE INDEX idx_names_filter ON names(country_id, year, name, gender);

-- GIN trigram index for glob pattern matching (Stage 1: name_glob filter)
-- Enables efficient ILIKE queries for patterns like 'alex*' or '*сан*'
CREATE INDEX idx_names_name_trgm ON names USING GIN (name gin_trgm_ops);

-- Index for dataset-level queries (e.g., auditing, re-ingestion)
CREATE INDEX idx_names_dataset ON names(dataset_id);

-- Optional: Add if needed for performance after benchmarking
-- Covering index to avoid table lookups during aggregation
-- CREATE INDEX idx_names_aggregation ON names(country_id, year, name, gender) INCLUDE (count);

COMMENT ON TABLE names IS 'Core fact table storing name occurrence counts by year, country, and gender';
COMMENT ON COLUMN names.gender IS 'Gender marker: M (male), F (female), U (unknown/nonbinary)';
COMMENT ON COLUMN names.count IS 'Number of occurrences for this name/year/gender/country combination';
COMMENT ON INDEX idx_names_filter IS 'Composite index for efficient filtering by country, year, and name';
COMMENT ON INDEX idx_names_name_trgm IS 'Trigram index for glob pattern matching (ILIKE optimization)';

-- ============================================================================
-- Optional: Audit Log Table
-- Purpose: Track all administrative actions for compliance and debugging
-- ============================================================================

CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    event VARCHAR(100) NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW() NOT NULL,
    user_id VARCHAR(255),
    user_email VARCHAR(255),
    ip_address INET,
    user_agent TEXT,
    resource_type VARCHAR(50),
    resource_id INTEGER,
    details JSONB,
    result VARCHAR(50) NOT NULL CHECK (result IN ('success', 'failure', 'error'))
);

-- Index for querying audit logs by event type
CREATE INDEX idx_audit_event ON audit_log(event, timestamp DESC);

-- Index for querying audit logs by user
CREATE INDEX idx_audit_user ON audit_log(user_id, timestamp DESC);

-- Index for querying audit logs by resource
CREATE INDEX idx_audit_resource ON audit_log(resource_type, resource_id, timestamp DESC);

COMMENT ON TABLE audit_log IS 'Audit trail for all administrative actions';
COMMENT ON COLUMN audit_log.details IS 'JSON object with additional context (file size, checksum, etc.)';

-- ============================================================================
-- Sample Data (Optional - for development/testing)
-- ============================================================================

-- Insert sample countries
INSERT INTO countries (code, name, data_source_name, data_source_url, data_source_description, data_source_requires_manual_download)
VALUES
    ('US', 'United States', 'Social Security Administration', 'https://www.ssa.gov/oact/babynames/', 'Annual baby name data from SSA records since 1880. Includes names with at least 5 occurrences per year.', true),
    ('UK', 'United Kingdom', 'Office for National Statistics', 'https://www.ons.gov.uk/peoplepopulationandcommunity/birthsdeathsandmarriages/livebirths', 'Baby names data for England and Wales.', true),
    ('SE', 'Sweden', 'Statistics Sweden', 'https://www.scb.se/', NULL, true);

-- ============================================================================
-- Database Configuration
-- ============================================================================

-- Set UTF-8 encoding (should be set at database creation, but verify)
-- CREATE DATABASE affirm_name WITH ENCODING 'UTF8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8';

-- ============================================================================
-- Performance Tuning Notes
-- ============================================================================

-- For production, consider these settings (adjust based on workload):
-- 
-- shared_buffers = 256MB (or 25% of RAM)
-- effective_cache_size = 1GB (or 50% of RAM)
-- maintenance_work_mem = 128MB
-- work_mem = 16MB
-- 
-- For trigram searches:
-- pg_trgm.similarity_threshold = 0.3 (default)
--
-- Monitor query performance with:
-- EXPLAIN ANALYZE <query>
--
-- Check index usage with:
-- SELECT schemaname, tablename, indexname, idx_scan, idx_tup_read, idx_tup_fetch
-- FROM pg_stat_user_indexes
-- WHERE schemaname = 'public';

-- ============================================================================
-- Migration Complete
-- ============================================================================

-- Verify extensions
SELECT extname, extversion FROM pg_extension WHERE extname = 'pg_trgm';

-- Verify tables
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
ORDER BY table_name;

-- Verify indexes
SELECT tablename, indexname, indexdef 
FROM pg_indexes 
WHERE schemaname = 'public'
ORDER BY tablename, indexname;