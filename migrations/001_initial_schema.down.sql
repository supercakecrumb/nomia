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