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
    WHERE status = 'queued';
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

-- Create update timestamp trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for countries table
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
-- INSERT INTO countries (code, name, source_url, attribution) VALUES
-- ('US', 'United States', 'https://www.ssa.gov/oact/babynames/', 'Social Security Administration'),
-- ('GB', 'United Kingdom', 'https://www.ons.gov.uk/peoplepopulationandcommunity/birthsdeathsandmarriages/livebirths', 'Office for National Statistics'),
-- ('CA', 'Canada', 'https://www.statcan.gc.ca/', 'Statistics Canada'),
-- ('AU', 'Australia', 'https://www.abs.gov.au/', 'Australian Bureau of Statistics');