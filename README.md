# Baby Name Statistics Platform

A web application that aggregates and visualizes baby name statistics from multiple countries.

**Current Version:** 0.1.0

[![CI](https://github.com/yourusername/affirm-name/workflows/CI/badge.svg)](https://github.com/yourusername/affirm-name/actions/workflows/ci.yml)
[![Docker](https://github.com/yourusername/affirm-name/workflows/Docker/badge.svg)](https://github.com/yourusername/affirm-name/actions/workflows/docker.yml)
[![Security](https://github.com/yourusername/affirm-name/workflows/Security/badge.svg)](https://github.com/yourusername/affirm-name/actions/workflows/security.yml)
[![codecov](https://codecov.io/gh/yourusername/affirm-name/branch/main/graph/badge.svg)](https://codecov.io/gh/yourusername/affirm-name)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/affirm-name)](https://goreportcard.com/report/github.com/yourusername/affirm-name)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Changelog](https://img.shields.io/badge/changelog-Keep%20a%20Changelog-blue)](CHANGELOG.md)
[![Semantic Versioning](https://img.shields.io/badge/semver-2.0.0-blue)](https://semver.org/)

## Implementation Status

### Phase 1: Foundation ✅
- ✅ Go project structure with proper module organization
- ✅ Configuration system with environment variable support
- ✅ Database connection layer with pgx connection pooling
- ✅ Migration system for database schema management
- ✅ Structured logging with standard library slog
- ✅ Comprehensive test coverage

### Phase 2: REST API ✅
- ✅ REST API framework with Gin
- ✅ Country management endpoints
- ✅ Dataset management endpoints
- ✅ Job management endpoints
- ✅ Middleware (logging, recovery, CORS)

### Phase 3: File Upload & Storage ✅
- ✅ File upload endpoint
- ✅ Local storage implementation
- ✅ S3 storage implementation
- ✅ Upload service with validation

### Phase 4: Parser Framework ✅
- ✅ Parser interface and registry
- ✅ US SSA parser implementation
- ✅ Name normalization
- ✅ Batch insertion with staging tables
- ✅ Parser service

### Phase 5: Background Worker ✅
- ✅ Worker pool with configurable concurrency
- ✅ Job processor with retry logic
- ✅ Exponential backoff for retries
- ✅ Graceful shutdown
- ✅ Integration tests

## Project Structure

```
affirm-name/
├── cmd/
│   ├── api/              # REST API server
│   ├── migrate/          # Database migration tool
│   └── worker/           # Background worker
├── internal/
│   ├── api/              # API handlers and middleware
│   ├── config/           # Configuration management
│   ├── database/         # Database connection layer
│   ├── logging/          # Logging setup
│   ├── model/            # Domain models
│   ├── parser/           # Parser framework and implementations
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   ├── storage/          # File storage (local/S3)
│   └── worker/           # Worker pool and processor
├── migrations/           # SQL migration files
├── tests/                # Integration tests
├── docs/                 # Architecture and API documentation
├── go.mod
└── README.md
```

## Prerequisites

### For Local Development
- Go 1.24 or higher
- PostgreSQL 15 or higher

### For Docker Deployment
- Docker 20.10 or higher
- Docker Compose 2.0 or higher

## Configuration

The application is configured via environment variables:

### Database Configuration
- `DATABASE_URL` (required) - PostgreSQL connection string
- `DATABASE_MAX_CONNECTIONS` (default: 100) - Maximum number of connections
- `DATABASE_MAX_IDLE` (default: 10) - Maximum idle connections
- `DATABASE_CONN_MAX_LIFETIME` (default: 1h) - Maximum connection lifetime
- `DATABASE_CONN_MAX_IDLE_TIME` (default: 10m) - Maximum idle time

### Storage Configuration
- `STORAGE_TYPE` (default: local) - Storage type: "local" or "s3"
- `STORAGE_PATH` (default: ./uploads) - Path for local storage
- `S3_BUCKET` - S3 bucket name (required if STORAGE_TYPE=s3)
- `S3_REGION` (default: us-east-1) - S3 region
- `S3_ENDPOINT` - S3 endpoint URL

### Server Configuration
- `SERVER_PORT` (default: 8080) - HTTP server port
- `SERVER_READ_TIMEOUT` (default: 30s) - Read timeout
- `SERVER_WRITE_TIMEOUT` (default: 30s) - Write timeout
- `SERVER_IDLE_TIMEOUT` (default: 120s) - Idle timeout

### Worker Configuration
- `WORKER_CONCURRENCY` (default: 4) - Number of concurrent workers
- `WORKER_POLL_INTERVAL` (default: 5s) - Job polling interval
- `WORKER_MAX_RETRIES` (default: 3) - Maximum retry attempts

### Logging Configuration
- `LOG_LEVEL` (default: info) - Log level: debug, info, warn, error
- `LOG_FORMAT` (default: json) - Log format: json or text

## Getting Started

### 1. Install Dependencies

```bash
go mod download
```

### 2. Set Up Database

Create a PostgreSQL database:

```bash
createdb affirm_name
```

Set the database URL:

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/affirm_name?sslmode=disable"
```

### 3. Run Migrations

Apply database migrations:

```bash
go run cmd/migrate/main.go up
```

Check migration status:

```bash
go run cmd/migrate/main.go status
```

Rollback last migration:

```bash
go run cmd/migrate/main.go down
```

### 4. Build the Project

```bash
# Build API server
go build -o api ./cmd/api

# Build worker
go build -o worker ./cmd/worker

# Build migration tool
go build -o migrate ./cmd/migrate
```

### 5. Run the Application

Start the API server:

```bash
./api
```

In a separate terminal, start the worker:

```bash
./worker
```

The API will be available at `http://localhost:8080`.

### 6. Run Tests

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Database Schema

The initial schema includes the following tables:

- **countries** - Country metadata and data sources
- **datasets** - Uploaded files and processing status
- **names** - Normalized baby name statistics
- **jobs** - Asynchronous job queue
- **api_keys** - API authentication keys

See [`docs/database-schema.md`](docs/database-schema.md) for detailed schema documentation.

## API Usage Examples

### Upload a Dataset

```bash
# Upload US SSA dataset
curl -X POST http://localhost:8080/v1/datasets/upload \
  -F "file=@yob2023.txt" \
  -F "country_id=1" \
  -F "year=2023" \
  -F "source_type=us_ssa"
```

### Query Names

```bash
# List all names
curl http://localhost:8080/v1/names

# Search for a specific name
curl "http://localhost:8080/v1/names/search?name=Emma&country_id=1"

# Get details for a specific name
curl http://localhost:8080/v1/names/Emma
```

### Check Job Status

```bash
# List all jobs
curl http://localhost:8080/v1/jobs

# Get specific job
curl http://localhost:8080/v1/jobs/{job_id}
```

## Versioning

This project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html) and maintains a [changelog](CHANGELOG.md) following the [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format.

### Version Format

- **MAJOR** version for incompatible API changes
- **MINOR** version for new functionality in a backwards compatible manner
- **PATCH** version for backwards compatible bug fixes

### Changelog Management

We use [Changie](https://changie.dev/) to manage our changelog. See [CONTRIBUTING.md](CONTRIBUTING.md) for details on creating changelog entries.

#### Quick Start

```bash
# Install Changie
go install github.com/miniscruff/changie@latest

# Create a new changelog entry
make changelog-new

# Batch changes for a release
make changelog-batch VERSION=0.2.0

# Merge changes into CHANGELOG.md
make changelog-merge
```

### Release Process

To create a new release:

```bash
# Patch release (0.1.0 -> 0.1.1)
make release VERSION=patch

# Minor release (0.1.0 -> 0.2.0)
make release VERSION=minor

# Major release (0.1.0 -> 1.0.0)
make release VERSION=major
```

The release script will:
1. Batch all unreleased changes
2. Update CHANGELOG.md
3. Create a git commit
4. Create a git tag
5. Optionally create a GitHub release

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed contribution guidelines.

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/config -v

# Run integration tests (requires database)
go test ./tests/... -v
```

### Code Style

This project follows standard Go conventions:
- Use `gofmt` for formatting
- Follow effective Go guidelines
- Write tests for new functionality
- Document exported functions and types
For detailed contribution guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

### Makefile Targets

```bash
# Build
make build              # Build all binaries
make clean              # Clean build artifacts

# Testing
make test               # Run tests with coverage
make lint               # Run linters
make fmt                # Format code
make check              # Run all checks (fmt, lint, test)

# Docker
make docker-build       # Build Docker images
make docker-up          # Start containers
make docker-down        # Stop containers

# Database
make migrate-up         # Run migrations
make migrate-down       # Rollback migrations

# Changelog
make changelog-new      # Create new changelog entry
make changelog-batch    # Batch unreleased changes
make changelog-merge    # Merge into CHANGELOG.md
make release            # Create new release

# Dependencies
make install-deps       # Install development tools
```


### Docker Development Workflow

1. Make code changes
2. Rebuild the affected service:
   ```bash
   docker-compose build api
   docker-compose up -d api
   ```
3. View logs to verify changes:
   ```bash
   docker-compose logs -f api
   ```

### Adding New Parsers

To add support for a new country's data format:

1. Create a new parser in `internal/parser/parsers/`
2. Implement the `Parser` interface
3. Register the parser in `internal/parser/registry.go`
4. Add tests for the new parser
5. Update documentation

See [`docs/worker.md`](docs/worker.md) for detailed parser implementation guide.

## Next Steps (Phase 2+)

The following features will be implemented in subsequent phases:

- [ ] REST API framework with Gin
- [ ] File upload and storage
- [ ] Parser framework for country-specific formats
- [ ] Background worker pool
- [ ] Query API endpoints
- [ ] Trend analysis
- [ ] Additional country parsers

See [`docs/implementation-plan.md`](docs/implementation-plan.md) for the complete roadmap.

## CI/CD and Deployment

### GitHub Actions Workflows

This project uses GitHub Actions for continuous integration and deployment:

#### CI Workflow (`.github/workflows/ci.yml`)
Runs on every push and pull request to `main` and `develop` branches:
- **Testing**: Runs tests on Go 1.21 and 1.22 with coverage reporting
- **Linting**: Runs golangci-lint to ensure code quality
- **Building**: Builds all binaries (api, worker, migrate)
- **Security**: Runs Gosec security scanner
- **Integration**: Runs integration tests with PostgreSQL

#### Docker Workflow (`.github/workflows/docker.yml`)
Builds and pushes Docker images on push to `main` and version tags:
- **Multi-platform builds**: Supports linux/amd64 and linux/arm64
- **Container registry**: Pushes to GitHub Container Registry (ghcr.io)
- **Security scanning**: Runs Trivy vulnerability scanner
- **Image testing**: Tests Docker images before pushing

#### Release Workflow (`.github/workflows/release.yml`)
Triggered on version tags (v*.*.*):
- **Changelog generation**: Uses Changie to generate release notes
- **Binary builds**: Builds binaries for multiple platforms (Linux, macOS, Windows)
- **Docker images**: Builds and pushes versioned Docker images
- **GitHub releases**: Creates GitHub releases with artifacts and checksums

#### Lint Workflow (`.github/workflows/lint.yml`)
Runs comprehensive linting on pull requests:
- **Go linting**: golangci-lint, gofmt, goimports, go vet, staticcheck
- **Spelling**: Checks for common misspellings
- **Markdown**: Lints markdown files
- **YAML**: Validates YAML syntax
- **Dockerfile**: Lints Dockerfile with hadolint
- **Dependencies**: Checks for outdated and unused dependencies

#### Security Workflow (`.github/workflows/security.yml`)
Runs daily security scans:
- **Gosec**: Go security scanner
- **Govulncheck**: Go vulnerability database checker
- **Trivy**: Container and filesystem vulnerability scanner
- **CodeQL**: GitHub's semantic code analysis
- **Secret scanning**: Detects exposed secrets with TruffleHog
- **Dependency review**: Reviews dependency changes in PRs
- **Automated issues**: Creates issues for security problems

### Dependabot Configuration

Dependabot automatically creates pull requests for:
- **Go modules**: Weekly updates for dependencies
- **GitHub Actions**: Weekly updates for workflow actions
- **Docker images**: Weekly updates for base images

### Docker Images

Docker images are available at `ghcr.io/yourusername/affirm-name-{service}`:

```bash
# Pull latest images
docker pull ghcr.io/yourusername/affirm-name-api:latest
docker pull ghcr.io/yourusername/affirm-name-worker:latest
docker pull ghcr.io/yourusername/affirm-name-migrate:latest

# Pull specific version
docker pull ghcr.io/yourusername/affirm-name-api:v0.1.0
```

### Deployment

#### Using Docker Compose

```bash
# Pull and start all services
docker-compose pull
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

#### Using Kubernetes

```bash
# Apply manifests (create these based on your needs)
kubectl apply -f k8s/

# Check deployment status
kubectl get pods
kubectl get services
```

#### Manual Deployment

1. Download binaries from [GitHub Releases](https://github.com/yourusername/affirm-name/releases)
2. Extract the archive for your platform
3. Set environment variables (see Configuration section)
4. Run migrations: `./migrate up`
5. Start services: `./api` and `./worker`

### Release Process

To create a new release:

1. **Create changelog entries** for your changes:
   ```bash
   changie new
   ```

2. **Batch changes** for the release:
   ```bash
   changie batch v0.2.0
   changie merge
   ```

3. **Commit and tag**:
   ```bash
   git add .
   git commit -m "chore: release v0.2.0"
   git tag v0.2.0
   git push origin main --tags
   ```

4. **GitHub Actions will automatically**:
   - Build binaries for all platforms
   - Create Docker images with version tags
   - Generate release notes from changelog
   - Create a GitHub release with artifacts

### Monitoring CI/CD

- **Workflow runs**: Check [Actions tab](https://github.com/yourusername/affirm-name/actions)
- **Security alerts**: Check [Security tab](https://github.com/yourusername/affirm-name/security)
- **Dependabot PRs**: Check [Pull requests](https://github.com/yourusername/affirm-name/pulls)
- **Code coverage**: Check [Codecov](https://codecov.io/gh/yourusername/affirm-name)

## Architecture

For detailed architecture documentation, see:
- [`docs/architecture.md`](docs/architecture.md) - System architecture
- [`docs/database-schema.md`](docs/database-schema.md) - Database design
- [`docs/api-specification.md`](docs/api-specification.md) - API documentation
- [`docs/implementation-plan.md`](docs/implementation-plan.md) - Implementation roadmap

## License

[Add your license here]
