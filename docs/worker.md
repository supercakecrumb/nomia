# Background Worker Documentation

## Overview

The background worker is responsible for asynchronously processing jobs, primarily parsing uploaded dataset files and inserting name records into the database.

## Architecture

The worker system consists of three main components:

### 1. Worker Pool (`internal/worker/pool.go`)

The worker pool manages multiple concurrent worker goroutines that poll for and process jobs.

**Key Features:**
- Configurable concurrency (number of worker goroutines)
- Configurable poll interval
- Graceful shutdown with WaitGroup
- Context-based cancellation

**Configuration:**
```env
WORKER_CONCURRENCY=4        # Number of concurrent workers
WORKER_POLL_INTERVAL=5s     # How often to poll for jobs
```

### 2. Job Processor (`internal/worker/processor.go`)

The processor handles the actual job processing logic, including:
- Locking jobs atomically
- Processing different job types
- Handling errors and retries
- Updating job and dataset status

**Supported Job Types:**
- `parse_dataset` - Parse and insert data from uploaded file
- `reprocess_dataset` - Delete existing data and reprocess

### 3. Worker Command (`cmd/worker/main.go`)

The main entry point that:
- Loads configuration
- Connects to database
- Initializes storage
- Creates processor and pool
- Handles graceful shutdown

## Job Processing Flow

```
1. Worker polls for available jobs
2. Lock next job atomically (using FOR UPDATE SKIP LOCKED)
3. Update dataset status to 'processing'
4. Load file from storage
5. Get appropriate parser for country
6. Parse and validate data
7. Insert records in batches
8. On success:
   - Update dataset status to 'completed'
   - Update job status to 'completed'
   - Store processing metrics
9. On failure:
   - Check retry attempts
   - If < max attempts: schedule retry with backoff
   - If >= max attempts: mark as 'failed'
```

## Retry Logic

The worker implements exponential backoff for transient errors:

**Backoff Schedule:**
- Attempt 1: 1 minute
- Attempt 2: 5 minutes
- Attempt 3: 15 minutes

**Configuration:**
```env
WORKER_MAX_RETRIES=3        # Maximum retry attempts
```

**Retryable Errors:**
- Database connection errors
- Storage/network errors
- Temporary file system errors

**Non-Retryable Errors:**
- Invalid file format
- Missing parser for country
- Invalid filename format
- Dataset not found

## Running the Worker

### Prerequisites

1. Database must be running and migrated
2. Storage must be configured (local or S3)
3. Environment variables must be set

### Start Worker

```bash
# Using binary
./worker

# Using go run
go run cmd/worker/main.go

# With custom config
DATABASE_URL=postgres://... WORKER_CONCURRENCY=8 ./worker
```

### Graceful Shutdown

The worker handles `SIGINT` and `SIGTERM` signals for graceful shutdown:

1. Stop accepting new jobs
2. Wait for in-progress jobs to complete
3. Close database connections
4. Exit cleanly

```bash
# Send shutdown signal
kill -SIGTERM <worker_pid>

# Or use Ctrl+C
```

## Monitoring

### Logging

The worker logs structured information about:
- Worker startup/shutdown
- Job locking and processing
- Success/failure with details
- Retry scheduling

**Log Levels:**
- `INFO` - Normal operations
- `WARN` - Recoverable issues
- `ERROR` - Processing failures

**Example Logs:**
```json
{
  "time": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "msg": "Locked job for processing",
  "worker_id": "worker-1",
  "job_id": "123e4567-e89b-12d3-a456-426614174000",
  "job_type": "parse_dataset",
  "dataset_id": "456e7890-e89b-12d3-a456-426614174000",
  "attempt": 1
}
```

### Metrics (Future Enhancement)

Recommended metrics to track:
- Jobs processed per minute
- Average processing time
- Success/failure rate
- Queue depth
- Worker utilization

## Deployment

### Single Worker

For development or low-volume environments:

```bash
./worker
```

### Multiple Workers

For high-volume environments, run multiple worker processes:

```bash
# Terminal 1
WORKER_CONCURRENCY=4 ./worker

# Terminal 2
WORKER_CONCURRENCY=4 ./worker

# Terminal 3
WORKER_CONCURRENCY=4 ./worker
```

This provides 12 total concurrent workers across 3 processes.

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o worker cmd/worker/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/worker .
CMD ["./worker"]
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: affirm-name-worker
spec:
  replicas: 3
  selector:
    matchLabels:
      app: affirm-name-worker
  template:
    metadata:
      labels:
        app: affirm-name-worker
    spec:
      containers:
      - name: worker
        image: affirm-name-worker:latest
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: affirm-name-secrets
              key: database-url
        - name: WORKER_CONCURRENCY
          value: "4"
        - name: WORKER_POLL_INTERVAL
          value: "5s"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

## Troubleshooting

### Worker Not Processing Jobs

**Check:**
1. Worker is running: `ps aux | grep worker`
2. Database connection: Check logs for connection errors
3. Jobs in queue: Query `jobs` table for queued jobs
4. Worker configuration: Verify `WORKER_CONCURRENCY` and `WORKER_POLL_INTERVAL`

### Jobs Failing Repeatedly

**Check:**
1. Job error messages in `jobs.last_error`
2. Dataset error messages in `datasets.error_message`
3. File exists in storage
4. Parser available for country code
5. File format is valid

### High Memory Usage

**Solutions:**
1. Reduce `WORKER_CONCURRENCY`
2. Implement batch size limits in parser
3. Add memory limits in deployment
4. Monitor for memory leaks

### Slow Processing

**Solutions:**
1. Increase `WORKER_CONCURRENCY`
2. Optimize database batch inserts
3. Add database indexes
4. Use faster storage (SSD, S3 with CDN)
5. Profile code for bottlenecks

## Testing

### Unit Tests

```bash
# Test processor logic
go test ./internal/worker/...
```

### Integration Tests

```bash
# Run integration tests (requires database)
go test ./tests/... -v

# Skip integration tests
go test ./tests/... -short
```

### Manual Testing

1. Start API server
2. Upload a dataset file
3. Start worker
4. Monitor logs for processing
5. Verify data in database

```bash
# Terminal 1: API
./api

# Terminal 2: Worker
./worker

# Terminal 3: Upload file
curl -X POST http://localhost:8080/v1/datasets/upload \
  -F "file=@yob2023.txt" \
  -F "country_id=<country-uuid>"
```

## Performance Tuning

### Database

- Increase connection pool size: `DATABASE_MAX_CONNECTIONS=200`
- Tune batch insert size in parser
- Add indexes on frequently queried columns

### Worker

- Increase concurrency: `WORKER_CONCURRENCY=8`
- Decrease poll interval: `WORKER_POLL_INTERVAL=1s`
- Run multiple worker processes

### Storage

- Use local SSD for better I/O
- Use S3 with CloudFront for distributed access
- Enable compression for large files

## Security Considerations

1. **Database Credentials**: Store in secrets manager, not environment files
2. **Storage Access**: Use IAM roles for S3, not access keys
3. **File Validation**: Always validate file format before processing
4. **Resource Limits**: Set memory and CPU limits to prevent DoS
5. **Logging**: Don't log sensitive data (user info, file contents)

## Future Enhancements

1. **Priority Queues**: Process high-priority jobs first
2. **Dead Letter Queue**: Move permanently failed jobs to separate queue
3. **Metrics Dashboard**: Real-time monitoring of worker performance
4. **Auto-scaling**: Scale workers based on queue depth
5. **Job Scheduling**: Schedule jobs for specific times
6. **Webhooks**: Notify external systems on job completion
7. **Progress Tracking**: Report progress for long-running jobs