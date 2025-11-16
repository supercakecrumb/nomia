-- Affirm Name - US Data Seed Migration
-- Version: 002
-- Description: Insert US country record for SSA name data
-- Author: Data Team
-- Date: 2025-11-16

-- ============================================================================
-- Insert US Country Record
-- ============================================================================

INSERT INTO countries (code, name, data_source_name, data_source_url, data_source_description, data_source_requires_manual_download)
VALUES (
    'US',
    'United States',
    'Social Security Administration',
    'https://www.ssa.gov/oact/babynames/',
    'Annual baby name data from SSA records since 1880. Includes all names with at least 5 occurrences.',
    true
) ON CONFLICT (code) DO NOTHING;

-- Verify the insert
SELECT id, code, name, data_source_name FROM countries WHERE code = 'US';