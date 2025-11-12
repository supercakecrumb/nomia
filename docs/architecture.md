
# Baby Name Statistics Platform - Architecture Document

## Executive Summary

This document outlines the architecture for a web application that aggregates and visualizes baby name statistics from multiple countries. The system ingests government-published datasets (CSV format), normalizes them into a unified schema, and provides APIs for querying and analyzing name trends.

**Key Design Principles:**
- **Extensibility**: Easy addition of new country parsers
- **Reliability**: Transactional data integrity with versioning
- **Performance**: Efficient handling of millions of records
- **Transparency**: Full audit trail of uploads and processing
- **Maintainability**: Clean separation of concerns

---

## 1. System Overview

### 1.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend Layer                          │
│                  (React + TypeScript - Future)                  │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTPS/JSON
┌────────────────────────────▼────────────────────────────────────┐
│                         API Gateway                             │
│                    (Go HTTP Server + Router)                    │
│  ┌──────────────┬──────────────┬──────────────┬──────────────┐ │
│  │   Upload     │    Query     │    Trend     │   Admin      │ │
│  │   Handler    │   Handler    │   Handler    │   Handler    │ │
│  └──────┬───────┴──────┬───────┴──────┬───────┴──────┬───────┘ │
└─────────┼──────────────┼──────────────┼──────────────┼─────────┘
          │              │              │              │
          │              │              │              │
┌─────────▼──────────────▼──────────────▼──────────────▼─────────┐
│                      Service Layer                              │
│  ┌──────────────┬──────────────┬──────────────┬──────────────┐ │
│  │   Upload     │    Name      │    Trend     │   Country    │ │
│  │   Service    │   Service    │   Service    │   Service    │ │
│  └──────┬───────┴──────┬───────┴──────┬───────┴──────┬───────┘ │
└─────────┼──────────────┼──────────────┼──────────────┼─────────┘
          │              │              │              │
          ▼              │              │              │
┌─────────────────┐      │              │              │
│   Job Queue     │      │              │              │
│   (PostgreSQL   │      │              │              │
│    based)       │      │              │              │
└────────┬────────┘      │              │              │
         │               │              │              │
         ▼               ▼              ▼              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Data Access Layer                          │
│  ┌──────────────┬──────────────┬──────────────┬──────────────┐ │
│  │   Dataset    │    Name      │    Job       │   Country    │ │
│  │   Repository │  Repository  │  Repository  │  Repository  │ │
│  └──────┬───────┴──────┬───────┴──────┬───────┴──────┬───────┘ │
└─────────┼──────────────┼──────────────┼──────────────┼─────────┘
          │              │              │              │
          └──────────────┴──────────────┴──────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                      PostgreSQL Database                        │
│  ┌──────────────┬──────────────┬──────────────┬──────────────┐ │
│  │  countries   │   datasets   │    names     │     jobs     │ │
│  └──────────────┴──────────────┴──────────────┴──────────────┘ │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                      File Storage Layer                         │
│              (Local FS or S3-compatible storage)                │
│                    /uploads/{dataset_id}/                       │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                      Background Workers                         │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Parser Worker Pool (configurable concurrency)           │  │
│  │  - Polls job queue                                        │  │
│  │  - Executes country-specific parsers                      │  │
│  │  - Updates job status and dataset metadata                │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Component Responsibilities

**API Gateway**
- HTTP request routing and validation
- Authentication/authorization middleware
- Request/response serialization
- Rate limiting and CORS handling

**Service Layer**
- Business logic orchestration
- Transaction management
- Cross-cutting concerns (logging, metrics)
- Service-to-service communication

**Data Access Layer**
- Database query abstraction
- Connection pooling
- Query optimization
- Transaction handling

**Background Workers**
- Asynchronous job processing
- File parsing and validation
- Data normalization and insertion
- Error handling and retry logic

**File Storage**
- Original file persistence for audit
- Temporary upload staging
- Configurable backend (local/S3)

---

## 2. Data Flow

### 2.1 Upload and Ingestion Flow

```
┌─────────┐
│  Admin  │
│  User   │
└────┬────┘
     │ 1. POST /api/v1/datasets/upload
     │    (multipart/form-data)
     │    - file: CSV file
     │    - country_id: UUID
     │    - metadata: JSON
     ▼
┌─────────────────┐
│  Upload Handler │
│                 │
│ 2. Validate:    │
│    - File size  │
│    - MIME type  │
│    - Country    │
│      exists     │
└────┬────────────┘
     │ 3. Create dataset record
     │    status = 'pending'
     ▼
┌─────────────────┐
│  File Storage   │
│                 │
│ 4. Save file to │
│    /uploads/    │
│    {dataset_id}/│
│    original.csv │
└────┬────────────┘
     │ 5. Create job record
     │    type = 'parse_dataset'
     ▼
┌─────────────────┐
│   Job Queue     │
│  (jobs table)   │
│                 │
│ 6. Job inserted │
│    status =     │
│    'queued'     │
└────┬────────────┘
     │
     │ 7. Return 202 Accepted
     │    { dataset_id, job_id }
     │
     ▼
┌─────────────────┐
│  HTTP Response  │
└─────────────────┘

     ┌──────────────────────────────────┐
     │  Background Worker (async)       │
     │                                  │
     │ 8. Poll for queued jobs          │
     │                                  │
     │ 9. Lock job (status = 'running')│
     │                                  │
     │ 10. Load file from storage       │
     │                                  │
     │ 11. Detect country parser        │
     │     (based on country_id)        │
     │                                  │
     │ 12. Parse CSV:                   │
     │     - Stream rows                │
     │     - Normalize fields           │
     │     - Validate data              │
     │     - Batch insert (1000 rows)   │
     │                                  │
     │ 13. Update dataset:              │
     │     - row_count                  │
     │     - status = 'completed'       │
     │     - processed_at               │
     │                                  │
     │ 14. Update job:                  │
     │     - status = 'completed'       │
     │     - completed_at               │
     │                                  │
     │ 15. On error:                    │
     │     - Rollback transaction       │
     │     - status = 'failed'          │
     │     - error_message              │
     │     - Retry logic (max 3)        │
     └──────────────────────────────────┘
```

### 2.2 Query Flow

```
┌─────────┐
│  Client │
└────┬────┘
     │ GET /api/v1/names?country=US&year=2020&gender=F&limit=100
     ▼
┌─────────────────┐
│  Query Handler  │
│                 │
│ 1. Parse params │
│ 2. Validate     │
│ 3. Set defaults │
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│  Name Service   │
│                 │
│ 4. Build query  │
│    with filters │
│ 5. Apply        │
│    pagination   │
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│ Name Repository │
│                 │
│ 6. Execute SQL  │
│    with indexes │
│ 7. Fetch rows   │
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│   PostgreSQL    │
│                 │
│ 8. Query with:  │
│    - WHERE      │
│    - ORDER BY   │
│    - LIMIT      │
│    - OFFSET     │
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│  HTTP Response  │
│                 │
│ 9. JSON:        │
│    {            │
│      data: [],  │
│      meta: {    │
│        total,   │
│        page,    │
│        limit    │
│      }          │
│    }            │
└─────────────────┘
```

### 2.3 Trend Analysis Flow

```
┌─────────┐
│  Client │
└────┬────┘
     │ GET /api/v1/trends/Emma?country=US&start_year=1970&end_year=2020
     ▼
┌─────────────────┐
│  Trend Handler  │
│                 │
│ 1. Parse params │
│ 2. Validate     │
│    name format  │
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│  Trend Service  │
│                 │
│ 3. Query names  │
│    grouped by   │
│    year, gender │
│ 4. Calculate    │
│    percentages  │
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│ Name Repository │
│                 │
│ 5. Aggregate    │
│    query with   │
│    GROUP BY     │
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│   PostgreSQL    │
│                 │
│ 6. Execute:     │
│    SELECT year, │
│      gender,    │
│      SUM(count) │
│    GROUP BY     │
│      year,      │
│      gender     │
└────┬────────────┘
     │
     ▼
┌─────────────────┐
│  HTTP Response  │
│                 │
│ 7. JSON:        │
│    {            │
│      name,      │
│      trends: [  │
│        {year,   │
│         gender, │
│         count,  │
│         rank}   │
│      ]          │
│    }            │
└─────────────────┘
```

---

## 3. Scalability Considerations

### 3.1 Data Scalability

**Current Scale Estimates:**
- 50 countries × 50 years × 10,000 names/year = ~25M records
- Average row size: ~100 bytes
- Total data size: ~2.5 GB (uncompressed)
- With indexes: ~5-7 GB

**Growth Projections:**
- New countries: +500k records/country
- Annual updates: +500k records/year
- 5-year projection: ~50M records, ~10 GB

**Scaling Strategies:**

1. **Partitioning** (when >50M records)
   - Partition `names` table by country_id or year
   - Use PostgreSQL declarative partitioning
   - Benefits: Faster queries, easier archival

2. **Indexing Strategy**
   - Composite indexes for common query patterns
   - Partial indexes for active datasets
   - Regular VACUUM and ANALYZE

3. **Caching Layer** (future)
   - Redis for frequently accessed aggregations
   - Cache popular name trends
   - TTL: 24 hours for static historical data

4. **Read Replicas** (when read load >1000 QPS)
   - Primary for writes
   - Replicas for read queries
   - Connection pooling with pgBouncer

### 3.2 Code Scalability

**Modular Architecture:**

```
affirm-name/
├── cmd/
│   ├── api/           # HTTP API server
│   ├── worker/        # Background job processor
│   └── migrate/       # Database migration tool
├── internal/
│   ├── api/
│   │   ├── handlers/  # HTTP handlers
│   │   ├── middleware/# Auth, logging, etc.
│   │   └── router/    # Route definitions
│   ├── service/
│   │   ├── upload/    # Upload orchestration
│   │   ├── name/      # Name queries
│   │   ├── trend/     # Trend analysis
│   │   └── country/   # Country management
│   ├── repository/
│   │   ├── dataset/   # Dataset CRUD
│   │   ├── name/      # Name CRUD
│   │   ├── job/       # Job queue
│   │   └── country/   # Country CRUD
│   ├── parser/
│   │   ├── registry.go    # Parser registration
│   │   ├── interface.go   # Parser interface
│   │   ├── normalizer.go  # Common normalization
│   │   └── parsers/
│   │       ├── us_ssa.go  # US SSA format
│   │       ├── uk_ons.go  # UK ONS format
│   │       └── ...        # Other countries
│   ├── storage/
│   │   ├── interface.go   # Storage abstraction
│   │   ├── local.go       # Local filesystem
│   │   └── s3.go          # S3-compatible
│   ├── worker/
│   │   ├── pool.go        # Worker pool
│   │   └── processor.go   # Job processing
│   ├── model/         # Domain models
│   └── config/        # Configuration
├── migrations/        # SQL migrations
├── docs/             # Documentation
└── tests/
    ├── integration/  # Integration tests
    └── fixtures/     # Test data
```

**Parser Extensibility:**

```go
// Parser interface for country-specific implementations
type Parser interface {
    // Parse reads CSV and yields normalized records
    Parse(ctx context.Context, reader io.Reader) (<-chan Record, <-chan error)
    
    // Validate checks if file matches expected format
    Validate(reader io.Reader) error
    
    // Metadata returns parser information
    Metadata() ParserMetadata
}

// Registry for parser discovery
type Registry struct {
    parsers map[string]Parser
}

func (r *Registry) Register(countryCode string, parser Parser) {
    r.parsers[countryCode] = parser
}

func (r *Registry) Get(countryCode string) (Parser, error) {
    parser, ok := r.parsers[countryCode]
    if !ok {
        return nil, ErrParserNotFound
    }
    return parser, nil
}
```

**Adding a New Country Parser:**

1. Implement `Parser` interface in `internal/parser/parsers/new_country.go`
2. Register in `internal/parser/registry.go` init function
3. Add country record to database
4. No changes to core ingestion logic required

### 3.3 Performance Optimization

**Database Optimizations:**

1. **Batch Inserts**
   - Insert 1000 rows per transaction
   - Use `COPY` for bulk loading (10x faster than INSERT)
   - Disable indexes during bulk load, rebuild after

2. **Query Optimization**
   - Use EXPLAIN ANALYZE for slow queries
   - Composite indexes for filter combinations
   - Materialized views for complex aggregations

3. **Connection Pooling**
   - Max connections: 100
   - Idle connections: 10
   - Connection lifetime: 1 hour

**API Optimizations:**

1. **Pagination**
   - Cursor-based for stable sorting
   - Default limit: 100, max: 1000
   - Include total count in metadata

2. **Response Compression**
   - gzip compression for JSON responses
   - Reduces bandwidth by ~70%

3. **Request Validation**
   - Early validation before DB queries
   - Schema validation with JSON Schema

---

## 4. Extendability Plan

### 4.1 New Country Support

**Steps to Add a Country:**

1. **Create Parser Implementation**
   ```go
   // internal/parser/parsers/france_insee.go
   type FranceINSEEParser struct {
       normalizer *Normalizer
   }
   
   func (p *FranceINSEEParser) Parse(ctx context.Context, r io.Reader) (<-chan Record, <-chan error) {
       // Country-specific parsing logic
   }
   ```

2. **Register Parser**
   ```go
   // internal/parser/registry.go
   func init() {
       registry.Register("FR", &parsers.FranceINSEEParser{})
   }
   ```

3. **Add Country Metadata**
   ```sql
   INSERT INTO countries (code, name, source_url, attribution)
   VALUES ('FR', 'France', 'https://insee.fr', 'INSEE');
   ```

4. **Test with Sample Data**
   - Add test fixtures
   - Run integration tests
   - Verify normalization

**Parser Variations Handled:**

- Different column names (name, prenom, vorname)
- Different gender encodings (M/F, Male/Female, 1/2)
- Different separators (comma, semicolon, tab)
- Different encodings (UTF-8, Latin-1, Windows-1252)
- Header vs. no header
- Multiple files per year vs. single file

### 4.2 New Data Formats

**Current Support:** CSV

**Future Extensions:**

1. **Excel (XLS/XLSX)**
   - Use `github.com/xuri/excelize` library
   - Convert to CSV internally
   - Same normalization pipeline

2. **JSON/XML**
   - Parse to intermediate format
   - Feed to normalizer
   - Same storage model

3. **API Integration**
   - Fetch data from government APIs
   - Schedule periodic updates
   - Automated ingestion

### 4.3 Reprocessing Strategy

**Scenarios Requiring Reprocessing:**

1. **Parser Bug Fixes**
   - Incorrect gender mapping
   - Name encoding issues
   - Data validation errors

2. **Schema Changes**
   - New normalization rules
   - Additional metadata fields

3. **Data Corrections**
   - Government dataset updates
   - Error corrections

**Reprocessing Workflow:**

```
1. Mark dataset as 'reprocessing'
2. Soft delete existing names (deleted_at = NOW())
3. Create new job with type = 'reprocess'
4. Parse file again with updated parser
5. Insert new records
6. Update dataset status to 'completed'
7. Hard delete old records after verification period (30 days)
```

**Implementation:**

```sql
-- Soft delete existing records
UPDATE names 
SET deleted_at = NOW() 
WHERE dataset_id = $1 AND deleted_at IS NULL;

-- Reprocess creates new records
-- Old records remain for rollback

-- After verification, hard delete
DELETE FROM names 
WHERE dataset_id = $1 
  AND deleted_at < NOW() - INTERVAL '30 days';
```

---

## 5. Failure and Error Handling

### 5.1 Upload Failures

**Failure Scenarios:**

1. **Invalid File Format**
   - Response: 400 Bad Request
   - Error: "Invalid CSV format: missing required columns"
   - Action: User must fix file and re-upload

2. **File Too Large**
   - Response: 413 Payload Too Large
   - Error: "File exceeds maximum size of 100MB"
   - Action: User must split file or compress

3. **Country Not Found**
   - Response: 404 Not Found
   - Error: "Country with ID {id} does not exist"
   - Action: User must create country first

4. **Storage Failure**
   - Response: 500 Internal Server Error
   - Error: "Failed to save file to storage"
   - Action: Retry upload, check storage health

### 5.2 Parsing Failures

**Failure Scenarios:**

1. **Malformed CSV**
   - Job status: 'failed'
   - Error: "Row 1234: invalid column count"
   - Action: Log error, mark dataset as failed
   - Recovery: User fixes file and re-uploads

2. **Encoding Issues**
   - Job status: 'failed'
   - Error: "Invalid UTF-8 sequence at byte 5678"
   - Action: Attempt auto-detection, fallback to Latin-1
   - Recovery: Manual encoding specification

3. **Data Validation Errors**
   - Job status: 'partial'
   - Error: "Skipped 10 rows with invalid data"
   - Action: Log skipped rows, continue processing
   - Recovery: Review error log, fix source data

4. **Database Constraint Violations**
   - Job status: 'failed'
   - Error: "Duplicate key violation: (country, year, name, gender)"
   - Action: Rollback transaction
   - Recovery: Check for duplicate dataset upload

### 5.3 Partial Import Handling

**Strategy: All-or-Nothing with Staging**

```sql
BEGIN;

-- Create staging table
CREATE TEMP TABLE names_staging (LIKE names INCLUDING ALL);

-- Insert parsed data to staging
INSERT INTO names_staging (...) VALUES (...);

-- Validate staging data
SELECT COUNT(*) FROM names_staging WHERE name IS NULL;
-- If validation fails, ROLLBACK

-- Move from staging to production
INSERT INTO names SELECT * FROM names_staging;

COMMIT;
```

**Benefits:**
- No partial data in production
- Fast rollback on error
- Validation before commit
- Atomic dataset insertion

### 5.4 Retry Logic

**Job Retry Strategy:**

```go
type Job struct {
    ID          uuid.UUID
    Status      string  // queued, running, failed, completed
    Attempts    int     // Current attempt number
    MaxAttempts int     // Maximum retry attempts (default: 3)
    LastError   string  // Error message from last attempt
    NextRetryAt time.Time // When to retry next
}

// Exponential backoff: 1min, 5min, 15min
func calculateBackoff(attempt int) time.Duration {
    return time.Duration(math.Pow(5, float64(attempt))) * time.Minute
}
```

**Retry Conditions:**
- Transient database errors (connection lost)
- Temporary storage unavailability
- Rate limiting from external services

**No Retry Conditions:**
- Invalid file format (permanent error)
- Data validation failures (requires user action)
- Authorization failures

### 5.5 Error Logging and Monitoring

**Structured Logging:**

```go
log.WithFields(log.Fields{
    "job_id":     jobID,
    "dataset_id": datasetID,
    "country":    countryCode,
    "row":        rowNumber,
    "error":      err.Error(),
}).Error("Failed to parse row")
```

**Metrics to Track:**
- Upload success/failure rate
- Average parsing time per dataset
- Job queue depth
- Error rate by error type
- Storage usage

**Alerting Thresholds:**
- Job failure rate >10%
- Job queue depth >100
- Average parsing time >5 minutes
- Storage usage >80%

---

## 6. Security and Access Control

### 6.1 Authentication Strategy

**Phase 1: API Key (MVP)**

```
Authorization: Bearer <api_key>
```

- Admin users receive API keys
- Keys stored hashed in database
- Keys have expiration dates
- Keys can be revoked

**Phase 2: JWT (Future)**

```
Authorization: Bearer <jwt_token>
```

- OAuth2/OIDC integration
- Short-lived access tokens (15 min)
- Refresh tokens (30 days)
- Role-based claims in JWT

### 6.2 Authorization Model

**Roles:**

1. **Admin**
   - Upload datasets
   - Manage countries
   - View all data
   - Reprocess datasets

2. **Viewer** (future)
   - Query names
   - View trends
   - Read-only access

**Endpoint Protection:**

```
POST   /api/v1/datasets/upload     -> Admin only
DELETE /api/v1/datasets/:id        -> Admin only
POST   /api/v1/countries            -> Admin only
GET    /api/v1/names                -> Public (rate limited)
GET    /api/v1/trends/:name         -> Public (rate limited)
```

### 6.3 Input Validation

**File Upload Validation:**

```go
type UploadValidator struct {
    MaxFileSize   int64  // 100 MB
    AllowedTypes  []string // text/csv, application/csv
    MaxFilenameLen int   // 255 characters
}

func (v *UploadValidator) Validate(file multipart.File, header *multipart.FileHeader) error {
    // Check file size
    if header.Size > v.MaxFileSize {
        return ErrFileTooLarge
    }
    
    // Check MIME type
    buffer := make([]byte, 512)
    file.Read(buffer)
    mimeType := http.DetectContentType(buffer)
    if !contains(v.AllowedTypes, mimeType) {
        return ErrInvalidFileType
    }
    
    // Check filename
    if len(header.Filename) > v.MaxFilenameLen {
        return ErrFilenameTooLong
    }
    
    // Reset file pointer
    file.Seek(0, 0)
    
    return nil
}
```

**Query Parameter Validation:**

```go
type NameQueryParams struct {
    Country  string `validate:"required,len=2"`
    Year     int    `validate:"required,min=1970,max=2030"`
    Gender   string `validate:"omitempty,oneof=M F"`
    Limit    int    `validate:"min=1,max=1000"`
    Offset   int    `validate:"min=0"`
}
```

### 6.4 Rate Limiting

**Strategy: Token Bucket per API Key**

```
Rate Limits:
- Admin endpoints: 100 req/min
- Public endpoints: 1000 req/min per IP
- Burst: 2x rate limit
```

**Implementation:**

```go
// Middleware using golang.org/x/time/rate
func RateLimitMiddleware(limiter *rate.Limiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": limiter.Reserve().Delay().Seconds(),
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 6.5 Data Security

**At Rest:**
- Database encryption (PostgreSQL TDE)
- File storage encryption (AES-256)
- Encrypted backups

**In Transit:**
- TLS 1.3 for all HTTP traffic
- Certificate pinning for internal services
- Secure WebSocket connections (future)

**Secrets Management:**
- Environment variables for configuration
- Vault/AWS Secrets Manager for production
- No secrets in code or logs

---

## 7. Operational Requirements

### 7.1 Deployment Architecture

**Containerized Deployment (Docker)**

```dockerfile
# Dockerfile for API server
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o api ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]
```

**Docker Compose (Development)**

```yaml
version: '3.8'
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: affirm_name
      POSTGRES_USER: affirm
      POSTGRES_PASSWORD: secret
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
  
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://affirm:secret@postgres:5432/affirm_name
      STORAGE_TYPE: local
      STORAGE_PATH: /data/uploads
    volumes:
      - upload_data:/data/uploads
    depends_on:
      - postgres
  
  worker:
    build: .
    command: ./worker
    environment:
      DATABASE_URL: postgres://affirm:secret@postgres:5432/affirm_name
      STORAGE_TYPE: local
      STORAGE_PATH: /data/uploads
      WORKER_CONCURRENCY: 4
    volumes:
      - upload_data:/data/uploads
    depends_on:
      - postgres

volumes:
  postgres_data:
  upload_data:
```

**Production Deployment Options:**

1. **Kubernetes**
   - API: Deployment with HPA (2-10 replicas)
   - Worker: Deployment with fixed replicas (4)
   - PostgreSQL: StatefulSet or managed service (RDS)
   - Storage: PersistentVolume or S3

2. **VM-based**
   - API: Load balanced across 2+ instances
   - Worker: Dedicated instance(s)
   - PostgreSQL: Primary + replica
   - Storage: NFS or object storage

### 7.2 Configuration Management

**Environment Variables:**

```bash
# Database
DATABASE_URL=postgres://user:pass@host:5432/dbname
DATABASE_MAX_CONNECTIONS=100
DATABASE_MAX_IDLE=10

# Storage
STORAGE_TYPE=s3  # or 'local'
STORAGE_PATH=/data/uploads  # for local
S3_BUCKET=affirm-name-uploads
S3_REGION=us-east-1
S3_ENDPOINT=https://s3.amazonaws.com

# Server
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

# Worker
WORKER_CONCURRENCY=4
WORKER_POLL_INTERVAL=5s
WORKER_MAX_RETRIES=3

# Auth
API_KEY_HASH_COST=12
JWT_SECRET=<secret>
JWT_EXPIRY=15m

# Logging
LOG_LEVEL=info  # debug, info, warn, error
LOG_FORMAT=json # json or text

# Monitoring
METRICS_ENABLED=true
METRICS_PORT=9090
```

**Configuration Struct:**

```go
type Config struct {
    Database DatabaseConfig
    Storage  StorageConfig
    Server   ServerConfig
    Worker   WorkerConfig
    Auth     AuthConfig
    Logging  LoggingConfig
    Metrics  MetricsConfig
}

func LoadConfig() (*Config, error) {
    // Load from environment variables
    // Validate required fields
    // Set defaults for optional fields
}
```

### 7.3 Logging and Observability

**Structured Logging with Levels:**

```go
// Use zerolog or logrus
log.Info().
    Str("dataset_id", datasetID).
    Str("country", countryCode).
    Int("row_count", rowCount).
    Msg("Dataset processing completed")

log.Error().
    Err(err).
    Str("job_id", jobID).
    Str("file_path", filePath).
    Msg("Failed to parse CSV file")
```

**Log Aggregation:**
- Stdout/stderr for containerized environments
- Collected by Fluentd/Fluent Bit
- Stored in Elasticsearch or Loki
- Visualized in Kibana or Grafana

**Metrics (Prometheus format):**

```go
// Counter: Total uploads
uploads_total{status="success|failure",country="US"}

// Histogram: Upload processing time
upload_processing_duration_seconds{country="US"}

// Gauge: Current job queue depth
job_queue_depth{status="queued|running"}

// Gauge: Active worker count
worker_active_count

// Counter: API requests
api_requests_total{method="GET",endpoint="/names",status="200"}

// Histogram: API response time
api_response_duration_seconds{endpoint="/names"}
```

**Health Checks:**

```
GET /health
Response: 200 OK
{
  "status": "healthy",
  "database": "connected",
  "storage": "accessible",
  "version": "1.0.0",
  "uptime": "24h30m"
}

GET /ready
Response: 200 OK (ready to serve traffic)
Response: 503 Service Unavailable (not ready)
```

**Tracing (OpenTelemetry):**
- Distributed tracing for request flows
- Trace upload → job creation → parsing → storage
- Identify bottlenecks and slow queries

### 7.4 Backup and Disaster Recovery

**Database Backups:**
- Daily full backups
- Continuous WAL archiving
- Point-in-time recovery capability
- Retention: 30 days
- Test restore monthly

**File Storage Backups:**
- Replicate to secondary storage
- Versioning enabled (S3 versioning)
- Retention: 90 days
- Lifecycle policies for archival

**Recovery Procedures:**

1. **Database Failure**
   - Promote replica to primary
   - Restore from backup if needed
   - RTO: 15 minutes, RPO: 5 minutes

2. **Storage Failure**
   - Switch to backup storage
   - Restore files from backup
   - RTO: 30 minutes, RPO: 24 hours

3. **Complete System Failure**
   - Deploy to new infrastructure
   - Restore database from backup
   - Restore files from backup
   - RTO: 2 hours, RPO: 24 hours

---

## 8. Architectural Decisions and Trade-offs

### 8.1 Key Decisions

**Decision 1: PostgreSQL-based Job Queue**

*Rationale:*
- Simplifies infrastructure (no separate queue service)
- ACID guarantees for job state
- Easy to query and monitor
- Sufficient for expected load (<1000 jobs/day)

*Trade-offs:*
- Not as scalable as dedicated queue (RabbitMQ, SQS)
- Polling overhead on database
- May need migration if load increases significantly

*When to reconsider:* >10,000 jobs/day or <1s latency required

**Decision 2: Soft Delete with Versioning**

*Rationale:*
- Enables safe reprocessing
- Audit trail for data changes
- Easy rollback on errors
- Compliance with data retention policies

*Trade-offs:*
- Increased storage requirements
- Queries must filter deleted records
- Periodic cleanup required

*Mitigation:* Partial indexes on non-deleted records

**Decision 3: Synchronous File Upload, Async Processing**

*Rationale:*
- Fast response to user (202 Accepted)
- Handles large files without timeout
- Better resource utilization
- Enables retry logic

*Trade-offs:*
- More complex than synchronous
- Requires job status polling
- Need background worker infrastructure

*Alternative considered:* Fully synchronous (rejected due to timeout risk)

**Decision 4: Country-Specific Parsers**

*Rationale:*
- Handles format variations cleanly
- Easy to add new countries
- Testable in isolation
- Clear ownership of parsing logic

*Trade-offs:*
- More code to maintain
- Duplication of common logic (mitigated by normalizer)

*Alternative considered:* Generic configurable parser (rejected due to complexity)

**Decision 5: Local/S3 Storage Abstraction**

*Rationale:*
- Flexibility for different environments
- Easy testing with local storage
- Production-ready with S3
- No vendor lock-in

*Trade-offs:*
- Abstraction layer adds complexity
- Must maintain two implementations

*Mitigation:* Shared interface, comprehensive tests

### 8.2 Assumptions

1. **Data Volume**: Datasets are <100MB, <1M rows per file
2. **Upload Frequency**: <100 uploads per day
3. **Query Load**: <1000 QPS for read endpoints
4. **Data Freshness**: Historical data, no real-time requirements
5. **User Base**: <1000 concurrent users initially
6. **Geographic Distribution**: Single region deployment sufficient
7. **Data Retention**: Indefinite retention of all datasets
8. **Compliance**: No GDPR/CCPA concerns (public government data)

### 8.3 Future Considerations

**When to Add:**

1. **Caching Layer (Redis)**
   - Trigger: Query response time >500ms
   - Benefit: 10x faster for popular queries
   - Cost: Additional infrastructure, cache invalidation complexity

2. **Read Replicas**
   - Trigger: Database CPU >70% sustained
   - Benefit: Horizontal read scaling
   - Cost: Replication lag, increased infrastructure

3. **CDN for API**
   - Trigger: Geographic latency >200ms
   - Benefit: Faster global access
   - Cost: CDN costs, cache invalidation

4. **Message Queue (RabbitMQ/SQS)**
   - Trigger: >10,000 jobs/day
   - Benefit: Better scalability, decoupling
   - Cost: Additional service, operational complexity

5. **Multi-tenancy**
   - Trigger: Need for organization-specific datasets
   - Benefit: SaaS business model
   - Cost: Schema changes, access control complexity

---

## 9. Implementation Roadmap

### Phase 1: Foundation (Weeks 1-2)

**Goals:**
- Database schema and migrations
- Basic project structure
- Configuration management

**Deliverables:**
1. PostgreSQL schema with all tables
2. Migration tool (golang-migrate)
3. Project scaffolding (cmd/, internal/)
4. Configuration loading
5. Database connection pooling
6. Basic logging setup

**Success Criteria:**
- Migrations run successfully
- Can connect to database
- Configuration loads from env vars

### Phase 2: Core API (Weeks 3-4)

**Goals:**
- REST API framework
- Country CRUD operations
- Basic authentication

**Deliverables:**
1. HTTP server with routing (Gin/Echo)
2. Country endpoints (CRUD)
3. API key authentication middleware
4. Request validation
5. Error handling
6. Health check endpoints

**Success Criteria:**
- Can create/read/update/delete countries
- Authentication works
- API returns proper error codes

### Phase 3: File Upload & Storage (Week 5)

**Goals:**
- File upload handling
- Storage abstraction
- Dataset management

**Deliverables:**
1. Upload endpoint (multipart/form-data)
2. File validation
3. Local storage implementation
4. S3 storage implementation
5. Dataset CRUD operations
6. Job queue table and basic operations

**Success Criteria:**
- Can upload files
- Files stored correctly
- Dataset records created
- Jobs queued

### Phase 4: Parser Framework (Week 6)

**Goals:**
- Parser interface and registry
- Normalizer implementation
- First country parser (US SSA)

**Deliverables:**
1. Parser interface definition
2. Parser registry
3. Normalizer with gender/name mapping
4. US SSA parser implementation
5. CSV streaming reader
6. Batch insertion logic

**Success Criteria:**
- US SSA files parse correctly
- Data normalized properly
- Can add new parsers easily

### Phase 5: Background Worker (Week 7)

**Goals:**
- Job processing system
- Worker pool
- Error handling and retry

**Deliverables:**
1. Worker pool implementation
2. Job polling and locking
3. Parser execution
4. Transaction management
5. Retry logic with backoff
6. Job status updates

**Success Criteria:**
- Jobs process asynchronously
- Errors handled gracefully
- Retries work correctly
- No data corruption

### Phase 6: Query API (Week 8)

**Goals:**
- Name listing endpoint
- Filtering and pagination
- Performance optimization

**Deliverables:**
1. Name query endpoint
2. Filter implementation (country, year, gender)
3. Pagination (offset-based)
4. Sorting (by count, name)
5. Database indexes
6. Query optimization

**Success Criteria:**
- Can query names with filters
- Pagination works
- Response time <100ms for typical queries

### Phase 7: Trend Analysis (Week 9)

**Goals:**
- Trend endpoint
- Aggregation queries
- Gender probability calculation

**Deliverables:**
1. Trend endpoint implementation
2. Aggregation by year/gender
3. Rank calculation
4. Gender probability endpoint
5. Caching for popular names (optional)

**Success Criteria:**
- Trend data accurate
- Performance acceptable (<500ms)
- Gender probabilities calculated correctly

### Phase 8: Additional Parsers (Week 10)

**Goals:**
- Add 2-3 more country parsers
- Validate extensibility

**Deliverables:**
1. UK ONS parser
2. Canada parser
3. Australia parser
4. Parser tests
5. Sample datasets

**Success Criteria:**
- All parsers work correctly
- Easy to add new parsers
- Data quality validated

### Phase 9: Testing & Documentation (Week 11)

**Goals:**
- Comprehensive test coverage
- API documentation
- Deployment guides

**Deliverables:**
1. Unit tests (>80% coverage)
2. Integration tests
3. API documentation (OpenAPI spec)
4. Deployment guide
5. Operations runbook
6. Sample data and scripts

**Success Criteria:**
- All tests pass
- Documentation complete
- Can deploy to production

### Phase 10: Production Readiness (Week 12)

**Goals:**
- Monitoring and alerting
- Performance tuning
- Security hardening

**Deliverables:**
1. Prometheus metrics
2. Grafana dashboards
3. Alert rules
4. Rate limiting
5. Security audit
6. Load testing results
7. Backup procedures

**Success Criteria:**
- Monitoring in place
- Performance targets met
- Security validated
- Ready for production traffic

---

## 10. Testing Strategy

### 10.1 Unit Tests

**Coverage Targets:**
- Parsers: 90%
- Services: 85%
- Repositories: 80%
- Handlers: 75%

**Key Test Areas:**

1. **Parser Tests**
   ```go
   func TestUSSSAParser_Parse(t *testing.T) {
       tests := []struct {
           name     string
           input    string
           expected []Record
           wantErr  bool
       }{
           {
               name: "valid CSV",
               input: "Emma,F,20000\nLiam,M,19000\n",
               expected: []Record{
                   {Name: "Emma", Gender: "F", Count: 20000},
                   {Name: "Liam", Gender: "M", Count: 19000},
               },
               wantErr: false,
           },
           {
               name: "invalid gender",
               input: "Emma,X,20000\n",
               wantErr: true,
           },
       }
       // Test implementation
   }
   ```

2. **Normalizer Tests**
   - Gender mapping (Male→M, Female→F, etc.)
   - Name trimming and case handling
   - Invalid data rejection

3. **Service Tests**
   - Business logic validation
   - Error handling
   - Transaction management

4. **Repository Tests**
   - CRUD operations
   - Query building
   - Constraint handling

### 10.2 Integration Tests

**Test Database:**
- Use testcontainers for PostgreSQL
- Isolated test database per test suite
- Automatic cleanup

**Test Scenarios:**

1. **End-to-End Upload Flow**
   ```go
   func TestUploadFlow(t *testing.T) {
       // 1. Create country
       // 2. Upload file
       // 3. Wait for job completion
       // 4. Verify data in database
       // 5. Query names
       // 6. Verify results
   }
   ```

2. **Parser Integration**
   - Parse real sample files
   - Verify data normalization
   - Check database state

3. **API Integration**
   - HTTP request/response testing
   - Authentication flow
   - Error responses

4. **Worker Integration**
   - Job processing
   - Retry logic
   - Error handling

### 10.3 Performance Tests

**Load Testing (k6 or Locust):**

```javascript
// k6 script
import http from 'k6/http';
import { check } from 'k6';

export let options = {
    stages: [
        { duration: '2m', target: 100 },  // Ramp up
        { duration: '5m', target: 100 },  // Steady state
        { duration: '2m', target: 0 },    // Ramp down
    ],
};

export default function() {
    let res = http.get('http://api/v1/names?country=US&year=2020');
    check(res, {
        'status is 200': (r) => r.status === 200,
        'response time < 200ms': (r) => r.timings.duration < 200,
    });
}
```

**Performance Targets:**
- Name query: p95 < 200ms
- Trend query: p95 < 500ms
- Upload endpoint: p95 < 1s
- Throughput: 1000 req/s

### 10.4 Test Data

**Fixtures:**

```
tests/fixtures/
├── countries/
│   └── sample_countries.sql
├── datasets/
│   ├── us_ssa_2020.csv
│   ├── uk_ons_2020.csv
│   └── invalid_format.csv
└── expected/
    ├── us_ssa_2020_normalized.json
    └── uk_ons_2020_normalized.json
```

**Test Data Generation:**
- Script to generate synthetic datasets
- Various sizes (100 rows, 10k rows, 100k rows)
- Edge cases (special characters, long names, etc.)

### 10.5 Continuous Integration

**CI Pipeline (GitHub Actions):**

```yaml
name: CI
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -html=coverage.out -o coverage.html
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
  
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: golangci/golangci-lint-action@v3
        with:
          version: latest
```

---

## 11. Conclusion

This architecture provides a solid foundation for the baby name statistics platform with:

✅ **Extensibility**: Easy to add new countries and data sources
✅ **Reliability**: Transactional integrity and error handling
✅ **Performance**: Optimized for millions of records
✅ **Maintainability**: Clean separation of concerns
✅ **Scalability**: Clear path to handle growth
✅ **Observability**: Comprehensive logging and monitoring

**Next Steps:**
1. Review and approve this architecture
2. Set up development environment
3. Begin Phase 1 implementation
4. Iterate based on feedback

**Key Success Factors:**
- Follow the phased roadmap
- Maintain test coverage
- Document as you build
- Regular code reviews
- Monitor performance metrics

---

## Appendix A: Glossary

- **Dataset**: A single uploaded file containing name statistics for a country
- **Parser**: Country-specific code that converts CSV to normalized format
- **Normalizer**: Common logic for standardizing data (gender, names, etc.)
- **Job**: Asynchronous task for processing uploaded files
- **Worker**: Background process that executes jobs
- **Soft Delete**: Marking records as deleted without physical removal
- **Staging Table**: Temporary table for validating data before insertion

## Appendix B: References

- PostgreSQL Documentation: https://www.postgresql.org/docs/
- Go Best Practices: https://go.dev/doc/effective_go
- REST API Design: https://restfulapi.net/
- OpenAPI Specification: https://swagger.io/specification/
- Prometheus Metrics: https://prometheus.io/docs/practices/naming/

## Appendix C: Contact and Support

- Architecture Questions: [Team Lead]
- Implementation Support: [Senior Developer]
- Infrastructure: [DevOps Team]
- Security: [Security Team]