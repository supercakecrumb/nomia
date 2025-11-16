# Logging Configuration Guide

The Affirm Name backend uses [Zap](https://github.com/uber-go/zap) for structured logging with configurable log levels.

## Quick Start

### Set Log Level via Environment Variable

```bash
# Debug mode (most verbose) - shows all logs including debug messages
LOG_LEVEL=debug go run cmd/server/main.go

# Info mode (default) - normal operation logs
LOG_LEVEL=info go run cmd/server/main.go

# Warn mode - only warnings and errors
LOG_LEVEL=warn go run cmd/server/main.go

# Error mode (least verbose) - only errors
LOG_LEVEL=error go run cmd/server/main.go
```

### Set in .env File

```env
# Add to backend/.env
LOG_LEVEL=debug
```

## Log Levels

| Level | Use Case | Output |
|-------|----------|--------|
| `debug` | Development, troubleshooting | Everything including debug messages |
| `info` | Production (default) | Normal operation logs, requests, etc. |
| `warn` | Production (quiet) | Only warnings and errors |
| `error` | Production (minimal) | Only errors |

## Log Format

### Development Mode (colored, human-readable)

```
2025-11-16 20:34:22.863	INFO	server/main.go:82	Database connected successfully
2025-11-16 20:34:22.863	INFO	server/main.go:107	Server starting	{"port": "8083", "fixture_mode": false}
2025-11-16 20:34:25.123	INFO	middleware/logging.go:40	HTTP request	{"method": "GET", "path": "/api/names", "status": 200, "duration": "2.1s"}
```

### Format Details

- **Time**: `2025-11-16 20:34:22.863` (UTC, milliseconds, no timezone)
- **Level**: `INFO`, `DEBUG`, `WARN`, `ERROR` (color-coded in terminal)
- **Caller**: `server/main.go:82` (file and line number)
- **Message**: Human-readable description
- **Fields**: Structured JSON fields (port, method, status, etc.)

## Usage Examples

### Development with Debug Logs

```bash
cd backend
LOG_LEVEL=debug go run cmd/server/main.go
```

**Output will include:**
- All INFO logs
- All DEBUG logs (if you add them in code)
- Request details
- Database queries (if logged)
- Function entry/exit points (if logged)

### Production with Info Logs

```bash
cd backend
LOG_LEVEL=info go run cmd/server/main.go
```

**Output includes:**
- Server startup
- HTTP requests
- Database connection status
- Errors and warnings

### Production with Minimal Logs

```bash
cd backend
LOG_LEVEL=error go run cmd/server/main.go
```

**Output includes:**
- Only errors and fatal issues
- No request logging (unless errors occur)

## Adding Debug Logs to Code

### In Handlers

```go
import "go.uber.org/zap"

func SomeHandler(cfg *config.Config, logger *zap.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Debug log for development
        logger.Debug("Processing request",
            zap.String("path", r.URL.Path),
            zap.String("method", r.Method),
        )
        
        // ... handler logic
        
        logger.Info("Request completed",
            zap.Int("status", 200),
        )
    }
}
```

### In Database Queries

```go
func (db *DB) SomeQuery(ctx context.Context, params *Params) (*Result, error) {
    logger := zap.L() // Get global logger
    
    logger.Debug("Executing query",
        zap.String("query_type", "names_list"),
        zap.Any("params", params),
    )
    
    // ... query execution
    
    logger.Debug("Query completed",
        zap.Int("results", len(results)),
        zap.Duration("duration", time.Since(start)),
    )
    
    return results, nil
}
```

## Request Logging

All HTTP requests are automatically logged by the middleware at [`internal/middleware/logging.go`](internal/middleware/logging.go):

```go
logger.Info("HTTP request",
    zap.String("method", r.Method),
    zap.String("path", r.URL.Path),
    zap.String("query", r.URL.RawQuery),
    zap.Int("status", wrapped.statusCode),
    zap.Duration("duration", duration),
)
```

**Logged for each request:**
- HTTP method (GET, POST, etc.)
- Request path (`/api/names`)
- Query parameters (`page=1&page_size=10`)
- Response status code (200, 404, 500, etc.)
- Request duration (in human-readable format: µs, ms, s)

## Log Output Formats

### Console Output (Development)

Colored, human-readable format perfect for local development.

### JSON Output (Production)

To switch to JSON format for production logging systems (e.g., ELK stack), modify [`cmd/server/main.go`](cmd/server/main.go):

```go
// Change from NewDevelopmentConfig to NewProductionConfig
loggerConfig := zap.NewProductionConfig()
```

JSON format example:
```json
{"level":"info","ts":"2025-11-16T20:34:22.863Z","caller":"server/main.go:82","msg":"Database connected"}
```

## Configuration Reference

### Environment Variable

```bash
LOG_LEVEL=debug|info|warn|error
```

### Code Configuration

See [`cmd/server/main.go`](cmd/server/main.go) `initLogger()` function:

```go
func initLogger(logLevel string) (*zap.Logger, error) {
    var level zapcore.Level
    switch logLevel {
    case "debug":
        level = zapcore.DebugLevel
    case "info":
        level = zapcore.InfoLevel
    case "warn":
        level = zapcore.WarnLevel
    case "error":
        level = zapcore.ErrorLevel
    default:
        level = zapcore.InfoLevel
    }
    
    loggerConfig := zap.NewDevelopmentConfig()
    loggerConfig.Level = zap.NewAtomicLevelAt(level)
    // ... time and format configuration
}
```

## Troubleshooting

### No Logs Appearing

Check log level - it might be set too high:
```bash
LOG_LEVEL=debug go run cmd/server/main.go
```

### Too Many Logs

Reduce log level:
```bash
LOG_LEVEL=warn go run cmd/server/main.go
```

### Timestamp Issues

Timestamps are in UTC by default. To use local time, modify the `EncodeTime` function in `initLogger()`.

## Performance Considerations

### Log Level Performance Impact

| Level | Impact | Recommended For |
|-------|--------|-----------------|
| debug | High (many allocations) | Local development only |
| info | Medium | Development, staging |
| warn | Low | Production |
| error | Minimal | High-traffic production |

### Best Practices

1. **Use appropriate levels**: Don't log at INFO in hot paths
2. **Structure your logs**: Use `zap.String()`, `zap.Int()` etc. for fields
3. **Avoid string concatenation**: Let Zap handle it
4. **Sample high-frequency logs**: Log 1 in 100 for very frequent events
5. **Use log levels correctly**:
   - `DEBUG`: Detailed information for debugging
   - `INFO`: General informational messages
   - `WARN`: Warning messages (non-critical issues)
   - `ERROR`: Error messages (recoverable errors)
   - `FATAL`: Fatal errors (causes program exit)

## Examples

### Debug Mode for Development

```bash
cd backend
LOG_LEVEL=debug FIXTURE_MODE=true go run cmd/server/main.go
```

**Use when:**
- Debugging issues
- Understanding request flow
- Troubleshooting database queries
- Local development

### Info Mode for Staging

```bash
cd backend
LOG_LEVEL=info FIXTURE_MODE=false go run cmd/server/main.go
```

**Use when:**
- Staging environment
- Pre-production testing
- Monitoring request patterns

### Warn Mode for Production

```bash
cd backend
LOG_LEVEL=warn go run cmd/server/main.go
```

**Use when:**
- High-traffic production
- Cost optimization (less log storage)
- Mature stable service

## Integration with Logging Services

### Structured Logging

All logs include structured fields making them easy to query:

```go
logger.Info("HTTP request",
    zap.String("method", "GET"),        // Filter by method
    zap.String("path", "/api/names"),   // Filter by endpoint
    zap.Int("status", 200),             // Filter by status
    zap.Duration("duration", 2*time.Second), // Filter by duration
)
```

### Query Examples (if using log aggregation service)

```
# Find slow requests
status:200 AND duration>1s

# Find errors
level:error AND path:/api/names

# Find specific endpoint
path:/api/meta/years
```

## Resources

- [Zap Documentation](https://pkg.go.dev/go.uber.org/zap)
- [Zap Best Practices](https://github.com/uber-go/zap#faq)
- [Structured Logging](https://www.honeycomb.io/blog/structured-logging-and-your-team)