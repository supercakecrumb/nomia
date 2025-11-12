
# Implementation Plan

## Overview

This document provides a comprehensive implementation plan for the baby name statistics platform, including detailed roadmap, parser design, upload service architecture, operational requirements, security patterns, and testing strategy.

---

## Table of Contents

1. [File Parsing and Ingestion Layer](#file-parsing-and-ingestion-layer)
2. [Upload Service Design](#upload-service-design)
3. [Operational Requirements](#operational-requirements)
4. [Security and Access Control](#security-and-access-control)
5. [Implementation Roadmap](#implementation-roadmap)
6. [Testing Strategy](#testing-strategy)

---

## File Parsing and Ingestion Layer

### Architecture Overview

The parsing layer is designed for extensibility, allowing easy addition of country-specific parsers while maintaining a common normalization pipeline.

```
┌─────────────────────────────────────────────────────────────┐
│                    Parser Architecture                      │
└─────────────────────────────────────────────────────────────┘

┌──────────────┐
│  CSV File    │
└──────┬───────┘
       │
       ▼
┌──────────────────────────────────────────────────────────┐
│              Parser Registry                             │
│  - Selects parser based on country_id                    │
│  - Returns country-specific parser implementation        │
└──────┬───────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────┐
│         Country-Specific Parser                          │
│  - US SSA Parser                                         │
│  - UK ONS Parser                                         │
│  - Canada Parser                                         │
│  - etc.                                                  │
│                                                          │
│  Responsibilities:                                       │
│  - Read CSV with correct format                         │
│  - Extract fields (name, gender, count, year)           │
│  - Handle country-specific quirks                       │
└──────┬───────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────┐
│              Normalizer                                  │
│  - Standardize gender (Male→M, Female→F)                │
│  - Trim and validate names                              │
│  - Validate counts and years                            │
│  - Handle encoding issues                               │
└──────┬───────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────┐
│              Batch Inserter                              │
│  - Accumulate records (batch size: 1000)                │
│  - Use staging table for validation                     │
│  - Atomic commit to production table                    │
│  - Handle constraint violations                         │
└──────┬───────────────────────────────────────────────────┘
       │
       ▼
┌──────────────┐
│  PostgreSQL  │
│  names table │
└──────────────┘
```

### Parser Interface

```go
// internal/parser/interface.go

package parser

import (
    "context"
    "io"
    "time"
)

// Record represents a normalized name record
type Record struct {
    Year    int
    Name    string
    Gender  string // "M" or "F"
    Count   int
}

// ParserMetadata provides information about the parser
type ParserMetadata struct {
    CountryCode string
    Name        string
    Version     string
    Description string
    SampleFormat string // Example CSV format
}

// Parser interface for country-specific implementations
type Parser interface {
    // Parse reads CSV and yields normalized records
    // Returns channels for records and errors
    // Closes channels when done
    Parse(ctx context.Context, reader io.Reader) (<-chan Record, <-chan error)
    
    // Validate checks if file matches expected format
    // Returns error if format is invalid
    Validate(reader io.Reader) error
    
    // Metadata returns parser information
    Metadata() ParserMetadata
}

// ParseResult contains parsing statistics
type ParseResult struct {
    RowsProcessed int
    RowsSkipped   int
    Errors        []ParseError
    Duration      time.Duration
}

// ParseError represents a parsing error
type ParseError struct {
    Row     int
    Column  string
    Message string
    Value   string
}
```

### Parser Registry

```go
// internal/parser/registry.go

package parser

import (
    "fmt"
    "sync"
)

var (
    globalRegistry = NewRegistry()
)

// Registry manages parser implementations
type Registry struct {
    mu      sync.RWMutex
    parsers map[string]Parser
}

// NewRegistry creates a new parser registry
func NewRegistry() *Registry {
    return &Registry{
        parsers: make(map[string]Parser),
    }
}

// Register adds a parser for a country
func (r *Registry) Register(countryCode string, parser Parser) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.parsers[countryCode] = parser
}

// Get retrieves a parser for a country
func (r *Registry) Get(countryCode string) (Parser, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    parser, ok := r.parsers[countryCode]
    if !ok {
        return nil, fmt.Errorf("parser not found for country: %s", countryCode)
    }
    return parser, nil
}

// List returns all registered country codes
func (r *Registry) List() []string {
    r.mu.RLock()
    defer r.mu.RUnlock()
    
    codes := make([]string, 0, len(r.parsers))
    for code := range r.parsers {
        codes = append(codes, code)
    }
    return codes
}

// Global registry functions
func Register(countryCode string, parser Parser) {
    globalRegistry.Register(countryCode, parser)
}

func Get(countryCode string) (Parser, error) {
    return globalRegistry.Get(countryCode)
}
```

### Normalizer

```go
// internal/parser/normalizer.go

package parser

import (
    "fmt"
    "strings"
    "unicode"
)

// Normalizer handles common data normalization
type Normalizer struct {
    genderMap map[string]string
}

// NewNormalizer creates a new normalizer
func NewNormalizer() *Normalizer {
    return &Normalizer{
        genderMap: map[string]string{
            "M":      "M",
            "F":      "F",
            "Male":   "M",
            "Female": "F",
            "MALE":   "M",
            "FEMALE": "F",
            "m":      "M",
            "f":      "F",
            "1":      "M",
            "2":      "F",
            "Boy":    "M",
            "Girl":   "F",
        },
    }
}

// NormalizeGender converts various gender representations to M or F
func (n *Normalizer) NormalizeGender(gender string) (string, error) {
    normalized, ok := n.genderMap[strings.TrimSpace(gender)]
    if !ok {
        return "", fmt.Errorf("invalid gender: %s", gender)
    }
    return normalized, nil
}

// NormalizeName cleans and validates a name
func (n *Normalizer) NormalizeName(name string) (string, error) {
    // Trim whitespace
    name = strings.TrimSpace(name)
    
    // Check length
    if len(name) == 0 {
        return "", fmt.Errorf("name cannot be empty")
    }
    if len(name) > 100 {
        return "", fmt.Errorf("name too long: %d characters", len(name))
    }
    
    // Check for valid characters (letters, spaces, hyphens, apostrophes)
    for _, r := range name {
        if !unicode.IsLetter(r) && r != ' ' && r != '-' && r != '\'' {
            return "", fmt.Errorf("invalid character in name: %c", r)
        }
    }
    
    return name, nil
}

// ValidateYear checks if year is in valid range
func (n *Normalizer) ValidateYear(year int) error {
    if year < 1800 || year > 2100 {
        return fmt.Errorf("year out of range: %d", year)
    }
    return nil
}

// ValidateCount checks if count is positive
func (n *Normalizer) ValidateCount(count int) error {
    if count <= 0 {
        return fmt.Errorf("count must be positive: %d", count)
    }
    return nil
}
```

### Example Parser: US SSA

```go
// internal/parser/parsers/us_ssa.go

package parsers

import (
    "context"
    "encoding/csv"
    "fmt"
    "io"
    "strconv"
    "strings"
    
    "github.com/affirm-name/internal/parser"
)

// USSSAParser parses US Social Security Administration format
// Format: name,gender,count (no header, comma-separated)
// Example: Emma,F,15581
type USSSAParser struct {
    normalizer *parser.Normalizer
}

// NewUSSSAParser creates a new US SSA parser
func NewUSSSAParser() *USSSAParser {
    return &USSSAParser{
        normalizer: parser.NewNormalizer(),
    }
}

// Parse implements Parser interface
func (p *USSSAParser) Parse(ctx context.Context, reader io.Reader) (<-chan parser.Record, <-chan error) {
    records := make(chan parser.Record, 100)
    errors := make(chan error, 10)
    
    go func() {
        defer close(records)
        defer close(errors)
        
        csvReader := csv.NewReader(reader)
        csvReader.FieldsPerRecord = 3
        csvReader.TrimLeadingSpace = true
        
        rowNum := 0
        for {
            select {
            case <-ctx.Done():
                errors <- ctx.Err()
                return
            default:
            }
            
            row, err := csvReader.Read()
            if err == io.EOF {
                return
            }
            if err != nil {
                errors <- fmt.Errorf("row %d: %w", rowNum, err)
                continue
            }
            
            rowNum++
            
            // Parse fields
            name := strings.TrimSpace(row[0])
            genderRaw := strings.TrimSpace(row[1])
            countStr := strings.TrimSpace(row[2])
            
            // Normalize name
            name, err = p.normalizer.NormalizeName(name)
            if err != nil {
                errors <- fmt.Errorf("row %d: %w", rowNum, err)
                continue
            }
            
            // Normalize gender
            gender, err := p.normalizer.NormalizeGender(genderRaw)
            if err != nil {
                errors <- fmt.Errorf("row %d: %w", rowNum, err)
                continue
            }
            
            // Parse count
            count, err := strconv.Atoi(countStr)
            if err != nil {
                errors <- fmt.Errorf("row %d: invalid count: %w", rowNum, err)
                continue
            }
            
            // Validate count
            if err := p.normalizer.ValidateCount(count); err != nil {
                errors <- fmt.Errorf("row %d: %w", rowNum, err)
                continue
            }
            
            // Send record
            records <- parser.Record{
                Name:   name,
                Gender: gender,
                Count:  count,
            }
        }
    }()
    
    return records, errors
}

// Validate implements Parser interface
func (p *USSSAParser) Validate(reader io.Reader) error {
    csvReader := csv.NewReader(reader)
    csvReader.FieldsPerRecord = 3
    
    // Read first row to validate format
    row, err := csvReader.Read()
    if err != nil {
        return fmt.Errorf("invalid CSV format: %w", err)
    }
    
    if len(row) != 3 {
        return fmt.Errorf("expected 3 columns, got %d", len(row))
    }
    
    // Validate first row has expected types
    if _, err := strconv.Atoi(strings.TrimSpace(row[2])); err != nil {
        return fmt.Errorf("third column must be numeric: %w", err)
    }
    
    return nil
}

// Metadata implements Parser interface
func (p *USSSAParser) Metadata() parser.ParserMetadata {
    return parser.ParserMetadata{
        CountryCode:  "US",
        Name:         "US Social Security Administration",
        Version:      "1.0.0",
        Description:  "Parses US SSA baby names format (name,gender,count)",
        SampleFormat: "Emma,F,15581\nLiam,M,19659",
    }
}

// Register parser on init
func init() {
    parser.Register("US", NewUSSSAParser())
}
```

### Batch Inserter

```go
// internal/parser/inserter.go

package parser

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/google/uuid"
)

// BatchInserter handles batch insertion of records
type BatchInserter struct {
    db        *sql.DB
    batchSize int
}

// NewBatchInserter creates a new batch inserter
func NewBatchInserter(db *sql.DB, batchSize int) *BatchInserter {
    return &BatchInserter{
        db:        db,
        batchSize: batchSize,
    }
}

// Insert processes records and inserts them in batches
func (b *BatchInserter) Insert(
    ctx context.Context,
    datasetID uuid.UUID,
    countryID uuid.UUID,
    year int,
    records <-chan Record,
) (int, error) {
    tx, err := b.db.BeginTx(ctx, nil)
    if err != nil {
        return 0, fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback()
    
    // Create staging table
    _, err = tx.ExecContext(ctx, `
        CREATE TEMP TABLE names_staging (
            dataset_id UUID NOT NULL,
            country_id UUID NOT NULL,
            year INTEGER NOT NULL,
            name VARCHAR(100) NOT NULL,
            gender CHAR(1) NOT NULL,
            count INTEGER NOT NULL
        ) ON COMMIT DROP
    `)
    if err != nil {
        return 0, fmt.Errorf("create staging table: %w", err)
    }
    
    // Prepare insert statement
    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO names_staging (dataset_id, country_id, year, name, gender, count)
        VALUES ($1, $2, $3, $4, $5, $6)
    `)
    if err != nil {
        return 0, fmt.Errorf("prepare statement: %w", err)
    }
    defer stmt.Close()
    
    // Insert records in batches
    totalRows := 0
    batch := make([]Record, 0, b.batchSize)
    
    for record := range records {
        batch = append(batch, record)
        
        if len(batch) >= b.batchSize {
            if err := b.insertBatch(ctx, stmt, datasetID, countryID, year, batch); err != nil {
                return totalRows, err
            }
            totalRows += len(batch)
            batch = batch[:0]
        }
    }
    
    // Insert remaining records
    if len(batch) > 0 {
        if err := b.insertBatch(ctx, stmt, datasetID, countryID, year, batch); err != nil {
            return totalRows, err
        }
        totalRows += len(batch)
    }
    
    // Validate staging data
    var invalidCount int
    err = tx.QueryRowContext(ctx, `
        SELECT COUNT(*) FROM names_staging
        WHERE name IS NULL OR gender NOT IN ('M', 'F') OR count <= 0
    `).Scan(&invalidCount)
    if err != nil {
        return totalRows, fmt.Errorf("validate staging data: %w", err)
    }
    if invalidCount > 0 {
        return totalRows, fmt.Errorf("found %d invalid records in staging", invalidCount)
    }
    
    // Move from staging to production
    _, err = tx.ExecContext(ctx, `
        INSERT INTO names (dataset_id, country_id, year, name, gender, count)
        SELECT dataset_id, country_id, year, name, gender, count
        FROM names_staging
    `)
    if err != nil {
        return totalRows, fmt.Errorf("insert to production: %w", err)
    }
    
    // Commit transaction
    if err := tx.Commit(); err != nil {
        return totalRows, fmt.Errorf("commit transaction: %w", err)
    }
    
    return totalRows, nil
}

// insertBatch inserts a batch of records
func (b *BatchInserter) insertBatch(
    ctx context.Context,
    stmt *sql.Stmt,
    datasetID uuid.UUID,
    countryID uuid.UUID,
    year int,
    batch []Record,
) error {
    for _, record := range batch {
        _, err := stmt.ExecContext(ctx,
            datasetID,
            countryID,
            year,
            record.Name,
            record.Gender,
            record.Count,
        )
        if err != nil {
            return fmt.Errorf("insert record: %w", err)
        }
    }
    return nil
}
```

### Streaming vs Batch Trade-offs

**Streaming Approach (Chosen):**
- ✅ Memory efficient (constant memory usage)
- ✅ Can handle files of any size
- ✅ Early error detection
- ✅ Progress tracking possible
- ❌ Slightly more complex code

**Batch Approach (Alternative):**
- ✅ Simpler code
- ✅ Faster for small files
- ❌ High memory usage for large files
- ❌ All-or-nothing processing
- ❌ No progress tracking

**Decision:** Use streaming with batch inserts (1000 rows per batch) for optimal balance.

---

## Upload Service Design

### Upload Flow Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                    Upload Service Flow                       │
└──────────────────────────────────────────────────────────────┘

1. HTTP Request
   │
   ▼
┌──────────────────────────────────────────────────────────────┐
│  Upload Handler                                              │
│  - Validate multipart form                                   │
│  - Check file size (<100MB)                                  │
│  - Verify MIME type (text/csv)                              │
│  - Validate country_id exists                               │
└──────┬───────────────────────────────────────────────────────┘
       │
       ▼
2. Create Dataset Record
   │
   ▼
┌──────────────────────────────────────────────────────────────┐
│  Dataset Repository                                          │
│  - INSERT INTO datasets                                      │
│  - status = 'pending'                                        │
│  - Generate UUID                                             │
│  - Record metadata                                           │
└──────┬───────────────────────────────────────────────────────┘
       │
       ▼
3. Save File to Storage
   │
   ▼
┌──────────────────────────────────────────────────────────────┐
│  Storage Service                                             │
│  - Create directory: uploads/{dataset_id}/                   │
│  - Save file: original.csv                                   │
│  - Calculate checksum                                        │
│  - Update dataset.file_path                                  │
└──────┬───────────────────────────────────────────────────────┘
       │
       ▼
4. Create Job
   │
   ▼
┌──────────────────────────────────────────────────────────────┐
│  Job Repository                                              │
│  - INSERT INTO jobs                                          │
│  - type = 'parse_dataset'                                    │
│  - status = 'queued'                                         │
│  - payload = {dataset_id, country_code}                      │
└──────┬───────────────────────────────────────────────────────┘
       │
       ▼
5. Return Response (202 Accepted)
   │
   ▼
┌──────────────────────────────────────────────────────────────┐
│  HTTP Response                                               │
│  {                                                           │
│    "dataset_id": "uuid",                                     │
│    "job_id": "uuid",                                         │
│    "status": "pending"                                       │
│  }                                                           │
└──────────────────────────────────────────────────────────────┘

   ┌────────────────────────────────────────────────────────┐
   │  Background Worker (Async)                             │
   │                                                        │
   │  6. Poll for queued jobs                               │
   │  7. Lock job (FOR UPDATE SKIP LOCKED)                  │
   │  8. Update dataset status = 'processing'               │
   │  9. Load file from storage                             │
   │  10. Get parser from registry                          │
   │  11. Parse and insert data                             │
   │  12. Update dataset status = 'completed'               │
   │  13. Update job status = 'completed'                   │
   └────────────────────────────────────────────────────────┘
```

### Upload Handler Implementation

```go
// internal/api/handlers/upload.go

package handlers

import (
    "fmt"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

const (
    MaxFileSize = 100 * 1024 * 1024 // 100 MB
)

type UploadHandler struct {
    uploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
    return &UploadHandler{
        uploadService: uploadService,
    }
}

func (h *UploadHandler) Upload(c *gin.Context) {
    // Parse multipart form
    if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "invalid_request",
                "message": "Failed to parse multipart form",
            },
        })
        return
    }
    
    // Get file
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "missing_file",
                "message": "File is required",
            },
        })
        return
    }
    defer file.Close()
    
    // Validate file size
    if header.Size > MaxFileSize {
        c.JSON(http.StatusRequestEntityTooLarge, gin.H{
            "error": gin.H{
                "code":    "file_too_large",
                "message": fmt.Sprintf("File exceeds maximum size of %d MB", MaxFileSize/(1024*1024)),
            },
        })
        return
    }
    
    // Get country_id
    countryIDStr := c.PostForm("country_id")
    if countryIDStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "missing_country_id",
                "message": "country_id is required",
            },
        })
        return
    }
    
    countryID, err := uuid.Parse(countryIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "invalid_country_id",
                "message": "country_id must be a valid UUID",
            },
        })
        return
    }
    
    // Get optional metadata
    metadata := c.PostForm("metadata")
    
    // Get uploader from context (set by auth middleware)
    uploader := c.GetString("user_id")
    
    // Upload file
    result, err := h.uploadService.Upload(c.Request.Context(), service.UploadParams{
        File:      file,
        Filename:  header.Filename,
        FileSize:  header.Size,
        CountryID: countryID,
        Metadata:  metadata,
        UploadedBy: uploader,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "upload_failed",
                "message": err.Error(),
            },
        })
        return
    }
    
    // Return 202 Accepted
    c.JSON(http.StatusAccepted, gin.H{
        "data": gin.H{
            "dataset_id": result.DatasetID,
            "job_id":     result.JobID,
            "status":     "pending",
            "message":    "Dataset uploaded successfully. Processing will begin shortly.",
        },
    })
}
```

### Storage Service

```go
// internal/storage/interface.go

package storage

import (
    "context"
    "io"
    
    "github.com/google/uuid"
)

// Storage interface for file storage abstraction
type Storage interface {
    // Save saves a file and returns its path
    Save(ctx context.Context, datasetID uuid.UUID, filename string, reader io.Reader) (string, error)
    
    // Load loads a file
    Load(ctx context.Context, path string) (io.ReadCloser, error)
    
    // Delete deletes a file
    Delete(ctx context.Context, path string) error
    
    // Exists checks if a file exists
    Exists(ctx context.Context, path string) (bool, error)
}
```

```go
// internal/storage/local.go

package storage

import (
    "context"
    "fmt"
    "io"
    "os"
    "path/filepath"
    
    "github.com/google/uuid"
)

// LocalStorage implements Storage for local filesystem
type LocalStorage struct {
    basePath string
}

// NewLocalStorage creates a new local storage
func NewLocalStorage(basePath string) *LocalStorage {
    return &LocalStorage{
        basePath: basePath,
    }
}

// Save implements Storage interface
func (s *LocalStorage) Save(ctx context.Context, datasetID uuid.UUID, filename string, reader io.Reader) (string, error) {
    // Create directory
    dir := filepath.Join(s.basePath, datasetID.String())
    if err := os.MkdirAll(dir, 0755); err != nil {
        return "", fmt.Errorf("create directory: %w", err)
    }
    
    // Create file
    path := filepath.Join(dir, filename)
    file, err := os.Create(path)
    if err != nil {
        return "", fmt.Errorf("create file: %w", err)
    }
    defer file.Close()
    
    // Copy data
    if _, err := io.Copy(file, reader); err != nil {
        return "", fmt.Errorf("write file: %w", err)
    }
    
    return path, nil
}

// Load implements Storage interface
func (s *LocalStorage) Load(ctx context.Context, path string) (io.ReadCloser, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, fmt.Errorf("open file: %w", err)
    }
    return file, nil
}

// Delete implements Storage interface
func (s *LocalStorage) Delete(ctx context.Context, path string) error {
    if err := os.Remove(path); err != nil {
        return fmt.Errorf("delete file: %w", err)
    }
    return nil
}

// Exists implements Storage interface
func (s *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}
```

### Worker Pool

```go
// internal/worker/pool.go

package worker

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"
)

// Pool manages a pool of workers
type Pool struct {
    workers     int
    pollInterval time.Duration
    processor   *Processor
    wg          sync.WaitGroup
    stopCh      chan struct{}
}

// NewPool creates a new worker pool
func NewPool(workers int, pollInterval time.Duration, processor *Processor) *Pool {
    return &Pool{
        workers:      workers,
        pollInterval: pollInterval,
        processor:    processor,
        stopCh:       make(chan struct{}),
    }
}

// Start starts the worker pool
func (p *Pool) Start(ctx context.Context) {
    log.Printf("Starting worker pool with %d workers", p.workers)
    
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(ctx, i)
    }
}

// Stop stops the worker pool
func (p *Pool) Stop() {
    log.Println("Stopping worker pool")
    close(p.stopCh)
    p.wg.Wait()
    log.Println("Worker pool stopped")
}

// worker is a single worker goroutine
func (p *Pool) worker(ctx context.Context, id int) {
    defer p.wg.Done()
    
    log.Printf("Worker %d started", id)
    ticker := time.NewTicker(p.pollInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            log.Printf("Worker %d stopped: context cancelled", id)
            return
        case <-p.stopCh:
            log.Printf("Worker %d stopped: pool shutdown", id)
            return
        case <-ticker.C:
            if err := p.processNextJob(ctx, id); err != nil {
                log.Printf("Worker %d error: %v", id, err)
            }
        }
    }
}

// processNextJob processes the next available job
func (p *Pool) processNextJob(ctx context.Context, workerID int) error {
    // Lock and get next job
    job, err := p.processor.LockNextJob(ctx, fmt.Sprintf("worker-%d", workerID))
    if err != nil {
        return err
    }
    if job == nil {
        // No jobs available
        return nil
    }
    
    log.Printf("Worker %d processing job %s", workerID, job.ID)
    
    // Process job
    if err := p.processor.ProcessJob(ctx, job); err != nil {
        log.Printf("Worker %d failed to process job %s: %v", workerID, job.ID, err)
        return err
    }
    
    log.Printf("Worker %d completed job %s", workerID, job.ID)
    return nil

}
```

### Job Processor

```go
// internal/worker/processor.go

package worker

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "github.com/affirm-name/internal/model"
    "github.com/affirm-name/internal/parser"
    "github.com/affirm-name/internal/repository"
    "github.com/affirm-name/internal/storage"
)

// Processor handles job processing
type Processor struct {
    db              *sql.DB
    jobRepo         *repository.JobRepository
    datasetRepo     *repository.DatasetRepository
    storage         storage.Storage
    parserRegistry  *parser.Registry
    batchInserter   *parser.BatchInserter
}

// NewProcessor creates a new job processor
func NewProcessor(
    db *sql.DB,
    jobRepo *repository.JobRepository,
    datasetRepo *repository.DatasetRepository,
    storage storage.Storage,
    parserRegistry *parser.Registry,
    batchInserter *parser.BatchInserter,
) *Processor {
    return &Processor{
        db:             db,
        jobRepo:        jobRepo,
        datasetRepo:    datasetRepo,
        storage:        storage,
        parserRegistry: parserRegistry,
        batchInserter:  batchInserter,
    }
}

// LockNextJob locks and returns the next available job
func (p *Processor) LockNextJob(ctx context.Context, workerID string) (*model.Job, error) {
    return p.jobRepo.LockNext(ctx, workerID)
}

// ProcessJob processes a job
func (p *Processor) ProcessJob(ctx context.Context, job *model.Job) error {
    startTime := time.Now()
    
    // Parse job payload
    var payload struct {
        DatasetID   uuid.UUID `json:"dataset_id"`
        CountryCode string    `json:"country_code"`
    }
    if err := json.Unmarshal(job.Payload, &payload); err != nil {
        return p.failJob(ctx, job, fmt.Errorf("invalid payload: %w", err))
    }
    
    // Get dataset
    dataset, err := p.datasetRepo.GetByID(ctx, payload.DatasetID)
    if err != nil {
        return p.failJob(ctx, job, fmt.Errorf("get dataset: %w", err))
    }
    
    // Update dataset status
    if err := p.datasetRepo.UpdateStatus(ctx, dataset.ID, "processing"); err != nil {
        return p.failJob(ctx, job, fmt.Errorf("update dataset status: %w", err))
    }
    
    // Load file from storage
    file, err := p.storage.Load(ctx, dataset.FilePath)
    if err != nil {
        return p.failJob(ctx, job, fmt.Errorf("load file: %w", err))
    }
    defer file.Close()
    
    // Get parser
    parser, err := p.parserRegistry.Get(payload.CountryCode)
    if err != nil {
        return p.failJob(ctx, job, fmt.Errorf("get parser: %w", err))
    }
    
    // Parse file
    records, errors := parser.Parse(ctx, file)
    
    // Collect errors in background
    var parseErrors []string
    go func() {
        for err := range errors {
            parseErrors = append(parseErrors, err.Error())
        }
    }()
    
    // Insert records
    rowCount, err := p.batchInserter.Insert(
        ctx,
        dataset.ID,
        dataset.CountryID,
        2020, // TODO: Extract year from metadata or filename
        records,
    )
    if err != nil {
        return p.failJob(ctx, job, fmt.Errorf("insert records: %w", err))
    }
    
    // Update dataset
    if err := p.datasetRepo.UpdateCompleted(ctx, dataset.ID, rowCount); err != nil {
        return p.failJob(ctx, job, fmt.Errorf("update dataset: %w", err))
    }
    
    // Complete job
    duration := time.Since(startTime)
    result := map[string]interface{}{
        "rows_processed": rowCount,
        "rows_skipped":   len(parseErrors),
        "errors":         parseErrors,
        "duration_seconds": duration.Seconds(),
    }
    
    if err := p.jobRepo.Complete(ctx, job.ID, result); err != nil {
        return fmt.Errorf("complete job: %w", err)
    }
    
    return nil
}

// failJob marks a job as failed
func (p *Processor) failJob(ctx context.Context, job *model.Job, err error) error {
    // Check if should retry
    if job.Attempts < job.MaxAttempts {
        // Calculate backoff
        backoff := time.Duration(1<<uint(job.Attempts)) * time.Minute
        nextRetry := time.Now().Add(backoff)
        
        if err := p.jobRepo.Retry(ctx, job.ID, err.Error(), nextRetry); err != nil {
            return fmt.Errorf("retry job: %w", err)
        }
        return nil
    }
    
    // Max attempts reached, fail permanently
    if err := p.jobRepo.Fail(ctx, job.ID, err.Error()); err != nil {
        return fmt.Errorf("fail job: %w", err)
    }
    
    // Update dataset status
    if err := p.datasetRepo.UpdateStatus(ctx, job.DatasetID, "failed"); err != nil {
        return fmt.Errorf("update dataset status: %w", err)
    }
    
    return err
}
```

---

## Operational Requirements

### Deployment

#### Docker Compose (Development)

```yaml
# docker-compose.yml
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
      - ./migrations:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U affirm"]
      interval: 10s
      timeout: 5s
      retries: 5

  api:
    build:
      context: .
      dockerfile: Dockerfile
      target: api
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://affirm:secret@postgres:5432/affirm_name?sslmode=disable
      STORAGE_TYPE: local
      STORAGE_PATH: /data/uploads
      SERVER_PORT: 8080
      LOG_LEVEL: debug
    volumes:
      - upload_data:/data/uploads
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

  worker:
    build:
      context: .
      dockerfile: Dockerfile
      target: worker
    environment:
      DATABASE_URL: postgres://affirm:secret@postgres:5432/affirm_name?sslmode=disable
      STORAGE_TYPE: local
      STORAGE_PATH: /data/uploads
      WORKER_CONCURRENCY: 4
      WORKER_POLL_INTERVAL: 5s
      LOG_LEVEL: debug
    volumes:
      - upload_data:/data/uploads
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/dashboards:/etc/grafana/provisioning/dashboards
      - ./grafana/datasources:/etc/grafana/provisioning/datasources

volumes:
  postgres_data:
  upload_data:
  prometheus_data:
  grafana_data:
```

#### Dockerfile

```dockerfile
# Multi-stage Dockerfile

# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build API server
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Build worker
RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/worker

# Build migrate tool
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate

# API stage
FROM alpine:latest AS api

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./api"]

# Worker stage
FROM alpine:latest AS worker

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/worker .

CMD ["./worker"]

# Migrate stage
FROM alpine:latest AS migrate

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations

ENTRYPOINT ["./migrate"]
```

### Configuration

#### Environment Variables

```bash
# .env.example

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/affirm_name?sslmode=disable
DATABASE_MAX_CONNECTIONS=100
DATABASE_MAX_IDLE=10
DATABASE_CONN_MAX_LIFETIME=1h
DATABASE_CONN_MAX_IDLE_TIME=10m

# Storage
STORAGE_TYPE=local  # or 's3'
STORAGE_PATH=/data/uploads  # for local
S3_BUCKET=affirm-name-uploads  # for s3
S3_REGION=us-east-1
S3_ENDPOINT=https://s3.amazonaws.com
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret

# Server
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=120s
SERVER_MAX_HEADER_BYTES=1048576

# Worker
WORKER_CONCURRENCY=4
WORKER_POLL_INTERVAL=5s
WORKER_MAX_RETRIES=3

# Auth
API_KEY_HASH_COST=12
JWT_SECRET=your_secret_key_here
JWT_EXPIRY=15m

# Logging
LOG_LEVEL=info  # debug, info, warn, error
LOG_FORMAT=json  # json or text

# Metrics
METRICS_ENABLED=true
METRICS_PORT=9090

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://app.affirm-name.com
CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Authorization,Content-Type
CORS_MAX_AGE=86400

# Rate Limiting
RATE_LIMIT_ADMIN=100  # requests per minute
RATE_LIMIT_PUBLIC=1000  # requests per minute
```

### Monitoring

#### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'affirm-name-api'
    static_configs:
      - targets: ['api:9090']
    metrics_path: '/metrics'

  - job_name: 'affirm-name-worker'
    static_configs:
      - targets: ['worker:9090']
    metrics_path: '/metrics'

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres-exporter:9187']
```

#### Grafana Dashboard

```json
{
  "dashboard": {
    "title": "Affirm Name Platform",
    "panels": [
      {
        "title": "API Request Rate",
        "targets": [
          {
            "expr": "rate(api_requests_total[5m])"
          }
        ]
      },
      {
        "title": "API Response Time (p95)",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(api_response_duration_seconds_bucket[5m]))"
          }
        ]
      },
      {
        "title": "Job Queue Depth",
        "targets": [
          {
            "expr": "job_queue_depth"
          }
        ]
      },
      {
        "title": "Upload Success Rate",
        "targets": [
          {
            "expr": "rate(uploads_total{status=\"success\"}[5m]) / rate(uploads_total[5m])"
          }
        ]
      }
    ]
  }
}
```

### Logging

#### Structured Logging Example

```go
// internal/logging/logger.go

package logging

import (
    "os"
    
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// InitLogger initializes the global logger
func InitLogger(level string, format string) {
    // Set log level
    switch level {
    case "debug":
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    case "info":
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    case "warn":
        zerolog.SetGlobalLevel(zerolog.WarnLevel)
    case "error":
        zerolog.SetGlobalLevel(zerolog.ErrorLevel)
    default:
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    }
    
    // Set format
    if format == "text" {
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
    }
    
    // Add caller information
    log.Logger = log.With().Caller().Logger()
}

// Example usage
func ExampleLogging() {
    log.Info().
        Str("dataset_id", "uuid").
        Str("country", "US").
        Int("row_count", 32033).
        Msg("Dataset processing completed")
    
    log.Error().
        Err(err).
        Str("job_id", "uuid").
        Str("file_path", "/path/to/file").
        Msg("Failed to parse CSV file")
}
```

### Health Checks

```go
// internal/api/handlers/health.go

package handlers

import (
    "context"
    "database/sql"
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
)

type HealthHandler struct {
    db      *sql.DB
    storage storage.Storage
    version string
    startTime time.Time
}

func NewHealthHandler(db *sql.DB, storage storage.Storage, version string) *HealthHandler {
    return &HealthHandler{
        db:        db,
        storage:   storage,
        version:   version,
        startTime: time.Now(),
    }
}

func (h *HealthHandler) Health(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
    defer cancel()
    
    checks := make(map[string]string)
    healthy := true
    
    // Check database
    if err := h.db.PingContext(ctx); err != nil {
        checks["database"] = "disconnected"
        healthy = false
    } else {
        checks["database"] = "connected"
    }
    
    // Check storage
    if _, err := h.storage.Exists(ctx, "."); err != nil {
        checks["storage"] = "inaccessible"
        healthy = false
    } else {
        checks["storage"] = "accessible"
    }
    
    status := "healthy"
    statusCode := http.StatusOK
    if !healthy {
        status = "unhealthy"
        statusCode = http.StatusServiceUnavailable
    }
    
    c.JSON(statusCode, gin.H{
        "status":          status,
        "version":         h.version,
        "uptime_seconds":  time.Since(h.startTime).Seconds(),
        "timestamp":       time.Now().UTC().Format(time.RFC3339),
        "checks":          checks,
    })
}

func (h *HealthHandler) Ready(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()
    
    // Check if database is ready
    if err := h.db.PingContext(ctx); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "ready": false,
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "ready": true,
    })
}
```

---

## Security and Access Control

### API Key Management

```go
// internal/auth/apikey.go

package auth

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "time"
    
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
)

// APIKey represents an API key
type APIKey struct {
    ID         uuid.UUID
    KeyHash    string
    Name       string
    Role       string
    ExpiresAt  *time.Time
    LastUsedAt *time.Time
    CreatedAt  time.Time
    RevokedAt  *time.Time
}

// APIKeyService manages API keys
type APIKeyService struct {
    repo *repository.APIKeyRepository
    cost int
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(repo *repository.APIKeyRepository, cost int) *APIKeyService {
    return &APIKeyService{
        repo: repo,
        cost: cost,
    }
}

// Generate generates a new API key
func (s *APIKeyService) Generate(ctx context.Context, name string, role string, expiresAt *time.Time) (string, *APIKey, error) {
    // Generate random key
    keyBytes := make([]byte, 32)
    if _, err := rand.Read(keyBytes); err != nil {
        return "", nil, fmt.Errorf("generate random key: %w", err)
    }
    key := "ak_" + base64.URLEncoding.EncodeToString(keyBytes)
    
    // Hash key
    hash, err := bcrypt.GenerateFromPassword([]byte(key), s.cost)
    if err != nil {
        return "", nil, fmt.Errorf("hash key: %w", err)
    }
    
    // Create API key record
    apiKey := &APIKey{
        ID:        uuid.New(),
        KeyHash:   string(hash),
        Name:      name,
        Role:      role,
        ExpiresAt: expiresAt,
        CreatedAt: time.Now(),
    }
    
    // Save to database
    if err := s.repo.Create(ctx, apiKey); err != nil {
        return "", nil, fmt.Errorf("save api key: %w", err)
    }
    
    return key, apiKey, nil
}

// Verify verifies an API key
func (s *APIKeyService) Verify(ctx context.Context, key string) (*APIKey, error) {
    // Get all active API keys (in production, use indexed lookup)
    apiKeys, err := s.repo.ListActive(ctx)
    if err != nil {
        return nil, fmt.Errorf("list api keys: %w", err)
    }
    
    // Check each key
    for _, apiKey := range apiKeys {
        if err := bcrypt.CompareHashAndPassword([]byte(apiKey.KeyHash), []byte(key)); err == nil {
            // Check expiration
            if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
                return nil, fmt.Errorf("api key expired")
            }
            
            // Update last used
            if err := s.repo.UpdateLastUsed(ctx, apiKey.ID); err != nil {
                // Log error but don't fail
            }
            
            return apiKey, nil
        }
    }
    
    return nil, fmt.Errorf("invalid api key")
}

// Revoke revokes an API key
func (s *APIKeyService) Revoke(ctx context.Context, id uuid.UUID) error {
    return s.repo.Revoke(ctx, id)
}
```

### Authentication Middleware

```go
// internal/api/middleware/auth.go

package middleware

import (
    "net/http"
    "strings"
    
    "github.com/gin-gonic/gin"
    "github.com/affirm-name/internal/auth"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(apiKeyService *auth.APIKeyService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "missing_authorization",
                    "message": "Authorization header is required",
                },
            })
            c.Abort()
            return
        }
        
        // Parse Bearer token
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "invalid_authorization",
                    "message": "Authorization header must be 'Bearer <token>'",
                },
            })
            c.Abort()
            return
        }
        
        token := parts[1]
        
        // Verify API key
        apiKey, err := apiKeyService.Verify(c.Request.Context(), token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "invalid_token",
                    "message": "Invalid or expired API key",
                },
            })
            c.Abort()
            return
        }
        
        // Set user context
        c.Set("user_id", apiKey.ID.String())
        c.Set("user_role", apiKey.Role)
        
        c.Next()
    }
}

// RequireRole creates a role-checking middleware
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("user_role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "unauthorized",
                    "message": "Authentication required",
                },
            })
            c.Abort()
            return
        }
        
        // Check if user has required role
        hasRole := false
        for _, role := range roles {
            if userRole == role {
                hasRole = true
                break
            }
        }
        
        if !hasRole {
            c.JSON(http.StatusForbidden, gin.H{
                "error": gin.H{
                    "code":    "forbidden",
                    "message": "Insufficient permissions",
                },
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### Rate Limiting

```go
// internal/api/middleware/ratelimit.go

package middleware

import (
    "net/http"
    "sync"
    "time"
    
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

// RateLimiter manages rate limiting
type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps int, burst int) *RateLimiter {
    return &RateLimiter{
        limiters: make(map[string]*rate.Limiter),
        rate:     rate.Limit(rps),
        burst:    burst,
    }
}

// GetLimiter gets or creates a limiter for a key
func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    limiter, exists := rl.limiters[key]
    if !exists {
        limiter = rate.NewLimiter(rl.rate, rl.burst)
        rl.limiters[key] = limiter
    }
    
    return limiter
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Use API key or IP as key
        key := c.GetString("user_id")
        if key == "" {
            key = c.ClientIP()
        }
        
        // Get limiter for key
        l := limiter.GetLimiter(key)
        
        // Check if allowed
        if !l.Allow() {
            // Calculate retry after
            reservation := l.Reserve()
            retryAfter := int(reservation.Delay().Seconds())
            reservation.Cancel()
            
            c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.burst))
            c.Header("X-RateLimit-Remaining", "0")
            c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(reservation.Delay()).Unix()))
            c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
            
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": gin.H{
                    "code":        "rate_limit_exceeded",
                    "message":     "Rate limit exceeded",
                    "retry_after": retryAfter,
                },
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

## Implementation Roadmap

### Phase 1: Foundation (Weeks 1-2)

**Objective:** Set up project structure, database, and basic configuration

**Tasks:**
1. Initialize Go module and project structure
2. Set up PostgreSQL database
3. Create migration files
4. Implement configuration loading
5. Set up database connection pooling
6. Implement basic logging

**Deliverables:**
- [ ] Project scaffolding complete
- [ ] Database schema created
- [ ] Migrations working
- [ ] Configuration system functional
- [ ] Can connect to database
- [ ] Logging configured

**Success Criteria:**
- `go build` succeeds
- Migrations run without errors
- Can query database
- Configuration loads from environment

### Phase 2: Core API (Weeks 3-4)

**Objective:** Build REST API framework and country management

**Tasks:**
1. Set up HTTP server (Gin framework)
2. Implement routing
3. Create country CRUD endpoints
4. Implement API key authentication
5. Add request validation
6. Implement error handling
7. Add health check endpoints

**Deliverables:**
- [ ] HTTP server running
- [ ] Country endpoints functional
- [ ] Authentication working
- [ ] Validation in place
- [ ] Error responses standardized
- [ ] Health checks operational

**Success Criteria:**
- Can create/read/update/delete countries via API
- Authentication blocks unauthorized requests
- Invalid requests return 400 with details
- Health endpoint returns 200

### Phase 3: File Upload & Storage (Week 5)

**Objective:** Implement file upload and storage abstraction

**Tasks:**
1. Create upload endpoint
2. Implement file validation
3. Build local storage implementation
4. Build S3 storage implementation
5. Create dataset CRUD operations
6. Set up job queue table

**Deliverables:**
- [ ] Upload endpoint functional
- [ ] File validation working
- [ ] Local storage operational
- [ ] S3 storage operational
- [ ] Dataset records created
- [ ] Jobs queued

**Success Criteria:**
- Can upload CSV files
- Files saved to storage
- Dataset records in database
- Job records created

### Phase 4: Parser Framework (Week 6)

**Objective:** Build extensible parser system

**Tasks:**
1. Define parser interface
2. Implement parser registry
3. Create normalizer
4. Build US SSA parser
5. Implement CSV streaming
6. Create batch inserter

**Deliverables:**
- [ ] Parser interface defined
- [ ] Registry functional
- [ ] Normalizer working
- [ ] US parser complete
- [ ] Streaming parser operational
- [ ] Batch insertion working

**Success Criteria:**
- US SSA files parse correctly
- Data normalized properly
- Can add new parsers easily
- Batch insertion efficient

### Phase 5: Background Worker (Week 7)

**Objective:** Implement asynchronous job processing

**Tasks:**
1. Create worker pool
2. Implement job polling
3. Build job processor
4. Add transaction management
5. Implement retry logic
6. Add job status updates

**Deliverables:**
- [ ] Worker pool running
- [ ] Jobs processed asynchronously
- [ ] Errors handled gracefully
- [ ] Retries working
- [ ] Status updates accurate

**Success Criteria:**
- Jobs process in background
- Failed jobs retry
- No data corruption
- Status tracking accurate

### Phase 6: Query API (Week 8)

**Objective:** Build name query endpoints

**Tasks:**
1. Create name query endpoint
2. Implement filtering
3. Add pagination
4. Implement sorting
5. Create database indexes
6. Optimize queries

**Deliverables:**
- [ ] Query endpoint functional
- [ ] Filters working
- [ ] Pagination operational
- [ ] Sorting implemented
- [ ] Indexes created
- [ ] Performance acceptable

**Success Criteria:**
- Can query names with filters
- Pagination works correctly
- Response time <100ms
- Handles large result sets

### Phase 7: Trend Analysis (Week 9)

**Objective:** Implement trend analysis endpoints

**Tasks:**
1. Create trend endpoint
2. Implement aggregation queries
3. Add rank calculation
4. Build gender probability endpoint
5. Optimize aggregations

**Deliverables:**
- [ ] Trend endpoint functional
- [ ] Aggregations working
- [ ] Ranks calculated
- [ ] Gender probability accurate

**Success Criteria:**
- Trend data correct
- Performance <500ms
- Gender probabilities accurate

### Phase 8: Additional Parsers (Week 10)

**Objective:** Add more country parsers

**Tasks:**
1. Implement UK ONS parser
2. Implement Canada parser
3. Implement Australia parser
4. Add parser tests
5. Create sample datasets

**Deliverables:**
- [ ] UK parser complete
- [ ] Canada parser complete
- [ ] Australia parser complete
- [ ] Tests passing
- [ ] Sample data available

**Success Criteria:**
- All parsers work correctly
- Easy to add new parsers
- Data quality validated

### Phase 9: Testing & Documentation (Week 11)

**Objective:** Comprehensive testing and documentation

**Tasks:**
1. Write unit tests
2. Write integration tests
3.
Create API documentation
4. Write deployment guide
5. Create operations runbook
6. Add sample data and scripts

**Deliverables:**
- [ ] Unit tests (>80% coverage)
- [ ] Integration tests passing
- [ ] API documentation complete
- [ ] Deployment guide written
- [ ] Operations runbook ready
- [ ] Sample data available

**Success Criteria:**
- All tests pass
- Documentation complete
- Can deploy to production

### Phase 10: Production Readiness (Week 12)

**Objective:** Prepare for production deployment

**Tasks:**
1. Set up Prometheus metrics
2. Create Grafana dashboards
3. Configure alert rules
4. Implement rate limiting
5. Conduct security audit
6. Perform load testing
7. Document backup procedures

**Deliverables:**
- [ ] Metrics collection working
- [ ] Dashboards created
- [ ] Alerts configured
- [ ] Rate limiting active
- [ ] Security validated
- [ ] Load tests passed
- [ ] Backup procedures documented

**Success Criteria:**
- Monitoring operational
- Performance targets met
- Security hardened
- Ready for production

---

## Testing Strategy

### Unit Tests

**Coverage Targets:**
- Parsers: 90%
- Services: 85%
- Repositories: 80%
- Handlers: 75%

**Example Test Structure:**

```go
// internal/parser/parsers/us_ssa_test.go

package parsers

import (
    "context"
    "strings"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUSSSAParser_Parse(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expected    []parser.Record
        expectError bool
    }{
        {
            name:  "valid CSV",
            input: "Emma,F,20000\nLiam,M,19000\n",
            expected: []parser.Record{
                {Name: "Emma", Gender: "F", Count: 20000},
                {Name: "Liam", Gender: "M", Count: 19000},
            },
            expectError: false,
        },
        {
            name:        "invalid gender",
            input:       "Emma,X,20000\n",
            expected:    nil,
            expectError: true,
        },
        {
            name:        "invalid count",
            input:       "Emma,F,abc\n",
            expected:    nil,
            expectError: true,
        },
        {
            name:        "empty name",
            input:       ",F,20000\n",
            expected:    nil,
            expectError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := NewUSSSAParser()
            reader := strings.NewReader(tt.input)
            
            records, errors := parser.Parse(context.Background(), reader)
            
            var result []parser.Record
            var errs []error
            
            for record := range records {
                result = append(result, record)
            }
            
            for err := range errors {
                errs = append(errs, err)
            }
            
            if tt.expectError {
                assert.NotEmpty(t, errs)
            } else {
                assert.Empty(t, errs)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}

func TestUSSSAParser_Validate(t *testing.T) {
    parser := NewUSSSAParser()
    
    t.Run("valid format", func(t *testing.T) {
        input := "Emma,F,20000\n"
        err := parser.Validate(strings.NewReader(input))
        assert.NoError(t, err)
    })
    
    t.Run("invalid format - wrong column count", func(t *testing.T) {
        input := "Emma,F\n"
        err := parser.Validate(strings.NewReader(input))
        assert.Error(t, err)
    })
    
    t.Run("invalid format - non-numeric count", func(t *testing.T) {
        input := "Emma,F,abc\n"
        err := parser.Validate(strings.NewReader(input))
        assert.Error(t, err)
    })
}
```

### Integration Tests

**Test Database Setup:**

```go
// tests/integration/setup.go

package integration

import (
    "context"
    "database/sql"
    "fmt"
    "testing"
    
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestDB(t *testing.T) (*sql.DB, func()) {
    ctx := context.Background()
    
    // Start PostgreSQL container
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_DB":       "test_affirm_name",
            "POSTGRES_USER":     "test",
            "POSTGRES_PASSWORD": "test",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections"),
    }
    
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("Failed to start container: %v", err)
    }
    
    // Get connection string
    host, err := container.Host(ctx)
    if err != nil {
        t.Fatalf("Failed to get host: %v", err)
    }
    
    port, err := container.MappedPort(ctx, "5432")
    if err != nil {
        t.Fatalf("Failed to get port: %v", err)
    }
    
    dsn := fmt.Sprintf("postgres://test:test@%s:%s/test_affirm_name?sslmode=disable", host, port.Port())
    
    // Connect to database
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    
    // Run migrations
    if err := runMigrations(db); err != nil {
        t.Fatalf("Failed to run migrations: %v", err)
    }
    
    // Return cleanup function
    cleanup := func() {
        db.Close()
        container.Terminate(ctx)
    }
    
    return db, cleanup
}
```

**End-to-End Test:**

```go
// tests/integration/upload_test.go

package integration

import (
    "bytes"
    "context"
    "io"
    "mime/multipart"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestUploadFlow(t *testing.T) {
    // Setup test database
    db, cleanup := SetupTestDB(t)
    defer cleanup()
    
    // Setup test server
    server := setupTestServer(db)
    
    // 1. Create country
    country := createTestCountry(t, server, "US", "United States")
    
    // 2. Upload file
    file := createTestFile(t, "Emma,F,20000\nLiam,M,19000\n")
    upload := uploadTestFile(t, server, country.ID, file)
    
    require.NotEmpty(t, upload.DatasetID)
    require.NotEmpty(t, upload.JobID)
    assert.Equal(t, "pending", upload.Status)
    
    // 3. Wait for job completion (or process synchronously in test)
    processJob(t, db, upload.JobID)
    
    // 4. Verify data in database
    var count int
    err := db.QueryRow(`
        SELECT COUNT(*) FROM names WHERE dataset_id = $1
    `, upload.DatasetID).Scan(&count)
    require.NoError(t, err)
    assert.Equal(t, 2, count)
    
    // 5. Query names via API
    names := queryNames(t, server, "US", 2020, "F")
    require.Len(t, names, 1)
    assert.Equal(t, "Emma", names[0].Name)
    assert.Equal(t, 20000, names[0].Count)
}
```

### Performance Tests

**Load Testing with k6:**

```javascript
// tests/load/query_test.js

import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '2m', target: 100 },   // Ramp up to 100 users
        { duration: '5m', target: 100 },   // Stay at 100 users
        { duration: '2m', target: 200 },   // Ramp up to 200 users
        { duration: '5m', target: 200 },   // Stay at 200 users
        { duration: '2m', target: 0 },     // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'],  // 95% of requests must complete below 200ms
        http_req_failed: ['rate<0.01'],    // Error rate must be below 1%
    },
};

export default function() {
    const countries = ['US', 'GB', 'CA', 'AU'];
    const years = [2015, 2016, 2017, 2018, 2019, 2020];
    const genders = ['M', 'F'];
    
    // Random query parameters
    const country = countries[Math.floor(Math.random() * countries.length)];
    const year = years[Math.floor(Math.random() * years.length)];
    const gender = genders[Math.floor(Math.random() * genders.length)];
    
    const res = http.get(
        `http://localhost:8080/v1/names?country=${country}&year=${year}&gender=${gender}&limit=100`
    );
    
    check(res, {
        'status is 200': (r) => r.status === 200,
        'response time < 200ms': (r) => r.timings.duration < 200,
        'has data': (r) => JSON.parse(r.body).data.length > 0,
    });
    
    sleep(1);
}
```

**Benchmark Tests:**

```go
// internal/parser/parsers/us_ssa_bench_test.go

package parsers

import (
    "context"
    "strings"
    "testing"
)

func BenchmarkUSSSAParser_Parse(b *testing.B) {
    parser := NewUSSSAParser()
    
    // Generate test data
    var sb strings.Builder
    for i := 0; i < 10000; i++ {
        sb.WriteString(fmt.Sprintf("Name%d,F,%d\n", i, 1000+i))
    }
    input := sb.String()
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        reader := strings.NewReader(input)
        records, errors := parser.Parse(context.Background(), reader)
        
        // Consume channels
        for range records {
        }
        for range errors {
        }
    }
}

func BenchmarkBatchInserter_Insert(b *testing.B) {
    db, cleanup := setupTestDB(b)
    defer cleanup()
    
    inserter := parser.NewBatchInserter(db, 1000)
    
    // Generate test records
    records := make(chan parser.Record, 10000)
    go func() {
        for i := 0; i < 10000; i++ {
            records <- parser.Record{
                Name:   fmt.Sprintf("Name%d", i),
                Gender: "F",
                Count:  1000 + i,
            }
        }
        close(records)
    }()
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _, err := inserter.Insert(
            context.Background(),
            uuid.New(),
            uuid.New(),
            2020,
            records,
        )
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Test Data

**Fixture Files:**

```
tests/fixtures/
├── countries/
│   └── sample_countries.sql
├── datasets/
│   ├── us_ssa_2020_valid.csv
│   ├── us_ssa_2020_invalid_gender.csv
│   ├── us_ssa_2020_invalid_count.csv
│   ├── uk_ons_2020.csv
│   └── large_dataset_100k.csv
└── expected/
    ├── us_ssa_2020_normalized.json
    └── uk_ons_2020_normalized.json
```

**Sample Data Generator:**

```go
// tests/fixtures/generator.go

package fixtures

import (
    "fmt"
    "math/rand"
    "os"
)

var names = []string{
    "Emma", "Olivia", "Ava", "Isabella", "Sophia",
    "Liam", "Noah", "Oliver", "Elijah", "William",
}

func GenerateCSV(filename string, rows int) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    for i := 0; i < rows; i++ {
        name := names[rand.Intn(len(names))]
        gender := []string{"M", "F"}[rand.Intn(2)]
        count := rand.Intn(20000) + 1000
        
        fmt.Fprintf(file, "%s,%s,%d\n", name, gender, count)
    }
    
    return nil
}
```

### Continuous Integration

**GitHub Actions Workflow:**

```yaml
# .github/workflows/ci.yml

name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
          POSTGRES_DB: test_affirm_name
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-
      
      - name: Download dependencies
        run: go mod download
      
      - name: Run migrations
        run: |
          go run cmd/migrate/main.go up
        env:
          DATABASE_URL: postgres://postgres:test@localhost:5432/test_affirm_name?sslmode=disable
      
      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.out ./...
        env:
          DATABASE_URL: postgres://postgres:test@localhost:5432/test_affirm_name?sslmode=disable
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
          flags: unittests
          name: codecov-umbrella
      
      - name: Run linter
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
  
  integration:
    runs-on: ubuntu-latest
    needs: test
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run integration tests
        run: |
          go test -v -tags=integration ./tests/integration/...
  
  build:
    runs-on: ubuntu-latest
    needs: [test, integration]
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2
      
      - name: Build Docker image
        uses: docker/build-push-action@v4
        with:
          context: .
          push: false
          tags: affirm-name:latest
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

---

## Conclusion

This implementation plan provides a comprehensive roadmap for building the baby name statistics platform. The plan includes:

✅ **Detailed Parser Architecture**: Extensible design for country-specific parsers
✅ **Upload Service Design**: Complete flow from upload to processing
✅ **Operational Requirements**: Deployment, monitoring, and logging
✅ **Security Patterns**: Authentication, authorization, and rate limiting
✅ **Phased Roadmap**: 12-week implementation schedule
✅ **Testing Strategy**: Unit, integration, and performance tests

**Key Success Factors:**
- Follow the phased approach
- Maintain test coverage
- Document as you build
- Regular code reviews
- Monitor performance metrics
- Iterate based on feedback

**Next Steps:**
1. Review and approve this plan
2. Set up development environment
3. Begin Phase 1 implementation
4. Schedule weekly progress reviews

The architecture is designed to be maintainable, performant, and easy to extend, providing a solid foundation for the platform's growth.