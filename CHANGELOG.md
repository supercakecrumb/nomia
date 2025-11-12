# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 0.1.0 - 2025-11-12

### Added

- Initial project setup with Go modules and dependency management
- PostgreSQL database integration with connection pooling
- Database migration system using golang-migrate
- RESTful API server with Chi router
- Structured logging with zerolog
- Configuration management with environment variables
- Docker and Docker Compose setup for development and production
- API endpoints for name data management:
  - GET /api/v1/names - List names with pagination and filtering
  - GET /api/v1/names/:id - Get name details
  - GET /api/v1/names/search - Search names by pattern
  - GET /api/v1/names/popular - Get popular names by year and country
- Country management endpoints:
  - GET /api/v1/countries - List all countries
  - GET /api/v1/countries/:code - Get country details
- Dataset management endpoints:
  - GET /api/v1/datasets - List datasets
  - GET /api/v1/datasets/:id - Get dataset details
- File upload and processing system:
  - POST /api/v1/upload - Upload dataset files
  - Support for local and S3 storage backends
- Background job processing system:
  - Worker pool for concurrent job processing
  - Job status tracking and monitoring
  - GET /api/v1/jobs - List jobs
  - GET /api/v1/jobs/:id - Get job status
- CSV parser framework with pluggable parsers:
  - US SSA baby names parser
  - Parser registry for extensibility
  - Name normalization and validation
- Database schema with optimized indexes:
  - Countries table with ISO codes
  - Datasets table for tracking data sources
  - Names table with composite indexes for performance
  - Jobs table for async processing
- Comprehensive test suite:
  - Unit tests for parsers and normalizers
  - Integration tests for API endpoints
  - Worker integration tests
- API middleware:
  - CORS support
  - Request logging
  - Panic recovery
- Documentation:
  - Architecture documentation
  - Database schema documentation
  - API specification
  - Implementation plan
  - Testing guide
  - Worker system documentation
- Development tools:
  - Docker initialization scripts
  - Environment configuration examples
  - Sample test fixtures

### Security

- Input validation for all API endpoints
- SQL injection prevention with parameterized queries
- File upload size limits and type validation
- Secure configuration management with environment variables