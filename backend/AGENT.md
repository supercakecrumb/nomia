# AI Agent Guide for Affirm Name Backend

This document helps AI agents quickly understand the Affirm Name backend project and work effectively with it.

## 🎯 Project Overview

**What is Affirm Name?**
A REST API that serves historical baby name data from multiple countries, enabling users to explore naming trends, gender balance, and popularity over time.

**Current Status:** Production-ready with 144 years of US data (1880-2024), 102K unique names, 370M+ occurrences

**Tech Stack:**
- **Language**: Go 1.21+
- **Database**: PostgreSQL 16 (Docker)
- **Router**: Chi v5.2.3
- **Config**: Viper
- **Logging**: Zap
- **Testing**: Go test + GitHub Actions

## 📁 Project Structure

```
affirm-name-backend/
├── backend/                    # Main Go application
│   ├── cmd/
│   │   ├── server/main.go     # HTTP server entry point
│   │   └── import/main.go     # Data import tool
│   ├── internal/
│   │   ├── config/            # Viper configuration
│   │   ├── db/                # Database layer (queries, connection)
│   │   ├── handlers/          # HTTP handlers
│   │   └── middleware/        # HTTP middleware (logging, CORS)
│   ├── scripts/               # Shell scripts (import, download, test)
│   ├── .env                   # Local configuration (gitignored)
│   └── Makefile              # Development commands
├── migrations/                # SQL migration files
├── spec-examples/             # JSON fixtures for frontend dev
├── architecture/              # Architecture documentation
├── docker-compose.yml         # Development database
└── docker-compose.prod.yml    # Production stack
```

## 🗂️ Key Files and Their Purposes

### Critical Files (Read These First)

| File | Purpose | When to Modify |
|------|---------|----------------|
| [`ARCHITECTURE.md`](ARCHITECTURE.md) | System design overview | Never (reference only) |
| [`architecture/02-backend-carcass.md`](architecture/02-backend-carcass.md) | Backend implementation guide | Never (reference only) |
| [`backend/internal/db/queries.go`](backend/internal/db/queries.go) | All database queries | Adding new queries or fixing bugs |
| [`backend/internal/handlers/*.go`](backend/internal/handlers/) | HTTP request handlers | Adding new endpoints or modifying responses |
| [`backend/cmd/server/main.go`](backend/cmd/server/main.go) | Server initialization | Adding middleware or routes |
| [`migrations/*.sql`](migrations/) | Database schema | Adding tables or columns |

### Configuration Files

| File | Purpose |
|------|---------|
| [`backend/.env`](backend/.env) | Local config (gitignored, create from .env.example) |
| [`backend/.env.example`](backend/.env.example) | Configuration template |
| [`backend/internal/config/config.go`](backend/internal/config/config.go) | Config struct and loading logic |

### Documentation Files

| File | Lines | Purpose |
|------|-------|---------|
| [`backend/DATABASE.md`](backend/DATABASE.md) | 76 | Database setup and management |
| [`backend/TESTING.md`](backend/TESTING.md) | 229 | Testing guide and best practices |
| [`backend/DATA_IMPORT.md`](backend/DATA_IMPORT.md) | 402 | Data import procedures |
| [`DEPLOYMENT.md`](DEPLOYMENT.md) | 248 | Production deployment guide |
| [`backend/LOGGING.md`](backend/LOGGING.md) | 237 | Logging configuration |

## 🚀 Quick Start for Agents

### 1. Understand Current State

```bash
# Check what's running
docker-compose ps

# Check database contents
docker-compose exec postgres psql -U postgres -d affirm_name -c "
SELECT 
    COUNT(DISTINCT year) as years,
    COUNT(DISTINCT name) as unique_names,
    COUNT(*) as total_records 
FROM names;"

# Check server status (if running)
curl http://localhost:8080/health
```

### 2. Start Development

```bash
# Start database
docker-compose up -d

# Run in fixture mode (for frontend dev)
cd backend && make dev

# Run in database mode (with real data)
cd backend && make prod

# Run tests
cd backend && make test
```

### 3. Common Development Tasks

**Run server:**
```bash
cd backend && go run cmd/server/main.go
```

**Run tests:**
```bash
cd backend && go test ./...
```

**Import data:**
```bash
bash backend/scripts/import-us-data.sh all
```

## 🗄️ Database Schema

### Tables

**`countries`** - Country metadata
- `id`, `code` (US, UK, SE), `name`, `data_source_*` fields

**`name_datasets`** - Tracks imported files
- `id`, `country_id`, `year_from`, `year_to`, `parse_status`

**`names`** - Main fact table (2M+ records)
- `id`, `country_id`, `dataset_id`, `year`, `name`, `gender`, `count`
- **Indexes**: composite on (country_id, year, name, gender), GIN trigram on name

**`audit_log`** - Change tracking
- `id`, `table_name`, `operation`, `user_id`, `changes`

### Gender Balance Calculation

```
gender_balance = 100 × (male_count / (male_count + female_count))

0   = 100% female
50  = perfectly neutral
100 = 100% male
NULL = no binary gender data
```

## 🔌 API Endpoints

### GET /health
**Purpose**: Health monitoring
**Returns**: `{status, timestamp, version, database}`

### GET /api/meta/years
**Purpose**: Get available year range
**Returns**: `{min_year, max_year}`

### GET /api/meta/countries
**Purpose**: List available countries
**Returns**: `{countries: [...]}`

### GET /api/names
**Purpose**: Core exploration endpoint
**Parameters**: 17 total (see [`architecture/01-shared-contract.md`](architecture/01-shared-contract.md))
- `year_from`, `year_to`, `countries`
- `gender_balance_min`, `gender_balance_max`
- `min_count`, `top_n`, `coverage_percent` (only one active)
- `name_glob` (supports `*` and `?`)
- `sort_key`, `sort_order`
- `page`, `page_size`

**Returns**: Paginated list with metadata

**Implementation**: 6-stage SQL CTE pipeline in [`backend/internal/db/queries.go`](backend/internal/db/queries.go:339-564)

### GET /api/names/trend
**Purpose**: Detailed information for a specific name
**Parameters**: `name` (required), `year_from`, `year_to`, `countries`
**Returns**: Summary, time series, country breakdown

## 🧪 Testing

### Run Tests
```bash
cd backend && make test           # All tests
cd backend && make test-cover     # With coverage
cd backend && make test-race      # With race detector
```

### Test Locations
- [`backend/internal/db/queries_test.go`](backend/internal/db/queries_test.go) - Parameter parsing tests
- [`backend/internal/handlers/params_test.go`](backend/internal/handlers/params_test.go) - Handler parameter tests

### CI/CD
- GitHub Actions: [`.github/workflows/test.yml`](.github/workflows/test.yml)
- Runs on: Push to any branch, pull requests
- Includes: Tests, linting, formatting, building

## 🐛 Common Issues and Solutions

### Issue: Port Already in Use
```bash
# Find and kill process
lsof -ti :8080 | xargs kill -9

# Or use different port
PORT=8081 go run cmd/server/main.go
```

### Issue: Database Connection Failed
```bash
# Check if PostgreSQL is running
docker-compose ps

# Start database
docker-compose up -d

# Check logs
docker-compose logs postgres
```

### Issue: Nil Pointer Dereference
**Fixed in queries.go** - Always check errors from `GetYearRange()`:
```go
yearRange, err := db.GetYearRange(ctx)
if err != nil {
    return nil, fmt.Errorf("failed to get year range: %w", err)
}
```

### Issue: Import Not Finding Files
- Use absolute paths: `/Users/username/Downloads/names`
- Not tilde paths: `~/Downloads/names` (won't work)

## 🎨 Code Patterns and Conventions

### Handler Pattern
```go
func EndpointName(cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Check fixture mode
        if cfg.FixtureMode {
            data, _ := LoadFixture("../spec-examples/file.json")
            WriteJSON(w, http.StatusOK, data)
            return
        }
        
        // 2. Parse and validate parameters
        params, err := ParseParams(r.URL.Query(), ...)
        if err != nil {
            http.Error(w, fmt.Sprintf("Invalid parameters: %v", err), 400)
            return
        }
        
        // 3. Query database
        result, err := cfg.DB.GetData(r.Context(), params)
        if err != nil {
            http.Error(w, fmt.Sprintf("Database error: %v", err), 500)
            return
        }
        
        // 4. Return JSON
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    }
}
```

### Database Query Pattern
```go
func (db *DB) GetSomething(ctx context.Context, params *Params) (*Result, error) {
    query := `SELECT ... FROM ... WHERE ...`
    
    rows, err := db.Pool.Query(ctx, query, param1, param2)
    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }
    defer rows.Close()
    
    var results []Item
    for rows.Next() {
        var item Item
        err := rows.Scan(&item.Field1, &item.Field2)
        if err != nil {
            return nil, fmt.Errorf("scan failed: %w", err)
        }
        results = append(results, item)
    }
    
    return &Result{Items: results}, nil
}
```

### Error Handling
- **Always** wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Return proper HTTP status codes: 400 (bad request), 404 (not found), 500 (server error)
- Log errors before returning them
- Never ignore errors (especially from database operations)

## 🔍 Important Implementation Details

### Fixture Mode vs Database Mode
Controlled by `FIXTURE_MODE` environment variable:
- `true` → Returns JSON from [`spec-examples/`](spec-examples/)
- `false` → Queries PostgreSQL database

**Both modes return identical JSON structure** (contract compliance)

### 6-Stage Filter Pipeline (GET /api/names)

Located in [`backend/internal/db/queries.go`](backend/internal/db/queries.go:343-434):

1. **Basic Filters** - Year, country, name pattern (ILIKE)
2. **Aggregation** - Group by name, sum counts, calculate gender_balance
3. **Gender Balance Filter** - Filter by male/female ratio
4. **Popularity Computation** - Calculate rank, cumulative share
5. **Popularity Filter** - Apply coverage_percent/top_n/min_count (only one active)
6. **Sorting & Pagination** - Sort by chosen key with tie-breaking

### Popularity Filter Priority
```go
if coverage_percent > 0 {
    // Highest priority
} else if top_n > 0 {
    // Second priority
} else if min_count > 0 {
    // Lowest priority
}
```

### Name Glob Patterns
Convert to SQL ILIKE:
- `*` → `%` (matches any sequence)
- `?` → `_` (matches single character)
- Example: `Alex*` → `Alex%` → matches Alexander, Alexis, etc.

## 🛠️ Adding New Features

### Adding a New API Endpoint

1. **Define handler** in `internal/handlers/`:
```go
func NewEndpoint(cfg *config.Config) http.HandlerFunc {
    // Implement handler
}
```

2. **Add database query** in `internal/db/queries.go`:
```go
func (db *DB) GetNewData(ctx context.Context, params *Params) (*Result, error) {
    // Implement query
}
```

3. **Register route** in `cmd/server/main.go`:
```go
r.Get("/api/new-endpoint", handlers.NewEndpoint(cfg))
```

4. **Add tests** in `internal/handlers/*_test.go`

5. **Update API documentation** in architecture docs

### Adding a New Country

1. **Add to** [`migrations/003_seed_all_countries.sql`](migrations/003_seed_all_countries.sql):
```sql
INSERT INTO countries (code, name, data_source_name, ...)
VALUES ('XX', 'Country Name', 'Data Source', ...);
```

2. **Document in** [`backend/data-sources.yml`](backend/data-sources.yml)

3. **Create parser** in [`backend/cmd/import/main.go`](backend/cmd/import/main.go)

4. **Import data**:
```bash
go run cmd/import/main.go -country=XX -dir=/path/to/data
```

### Adding Database Fields

1. **Create migration** in `migrations/00X_description.sql`
2. **Run migration**:
```bash
docker-compose exec postgres psql -U postgres -d affirm_name -f /docker-entrypoint-initdb.d/00X_description.sql
```
3. **Update Go structs** in `internal/db/queries.go`
4. **Update queries** to include new fields
5. **Update tests**

## 🔒 Important Constraints

### File Editing Restrictions
**Documentation Writer mode can only edit**:
- `*.md` files (all documentation)

**Other modes** may have different restrictions. Check mode capabilities before editing.

### Database Constraints
- `country_id` foreign key required for all `names` records
- `gender` must be 'M', 'F', or 'U'
- `year` should be between 1000-9999
- Duplicate records prevented by dataset tracking

### API Constraints
- `page`: 1-100
- `page_size`: 10-100
- `gender_balance_min/max`: 0-100
- `year_from` ≤ `year_to`
- Only ONE popularity filter active (coverage_percent > top_n > min_count)

## 🚨 Critical Areas (Handle with Care)

### 1. SQL ORDER BY in queries.go
**Lines 421-438**: The sorting logic is complex with multiple CASE statements. 
- **Test thoroughly** after any changes
- **Verify** with actual data, not just fixture mode
- **Check** all sort_key values: popularity, total_count, name, gender_balance

### 2. Gender Balance Calculation
**Used in multiple places** - keep consistent:
```go
100.0 * male_count / NULLIF(male_count + female_count, 0)
```
- Returns NULL if no binary gender data
- Scale is 0-100 (not 0-1)

### 3. Connection Pooling
**Important**: Always use `cfg.DB.Pool` (not individual connections)
- Pool is created in `main.go`
- Passed via `cfg` to all handlers
- Closed on server shutdown

### 4. Error Handling in queries.go
**Always check** `GetYearRange()` errors:
```go
yearRange, err := db.GetYearRange(ctx)
if err != nil {
    return nil, fmt.Errorf("failed to get year range: %w", err)
}
// Never use: yearRange, _ := db.GetYearRange(ctx)
```

## 📊 Data Flow

### Request Flow
```
HTTP Request
  ↓
Chi Router (cmd/server/main.go)
  ↓
Middleware (logging, CORS)
  ↓
Handler (internal/handlers/*.go)
  ↓
[If Fixture Mode] → Load JSON from spec-examples/
[If Database Mode] → Query Database
  ↓
Database Query (internal/db/queries.go)
  ↓
PostgreSQL (via pgxpool)
  ↓
Parse Results
  ↓
Return JSON Response
```

### Import Flow
```
Data Files (yobYYYY.txt)
  ↓
Import Tool (cmd/import/main.go)
  ↓
Parse CSV
  ↓
Create dataset record (name_datasets table)
  ↓
Batch Insert (1000 records at a time)
  ↓
Names Table
```

## 🧪 Testing Strategy

### Before Making Changes
1. Run existing tests: `cd backend && make test`
2. Check current functionality works

### After Making Changes
1. Run tests: `make test`
2. Test manually with curl or browser
3. Check logs for errors
4. Verify with both fixture and database modes

### Testing Checklist
- [ ] Unit tests pass
- [ ] Integration tests pass (if applicable)
- [ ] Manual endpoint testing
- [ ] Both fixture and database modes work
- [ ] No compilation errors
- [ ] No linter warnings
- [ ] Logs are clean

## 🎯 User's Workflow Preferences

Based on completed phases, the user prefers:

1. **Iterative Development**: Implement feature, test, commit, repeat
2. **Step-by-step Commits**: Each significant change gets its own commit
3. **Comprehensive Testing**: Test everything thoroughly before calling complete
4. **Production Quality**: Code should be production-ready, not prototypes
5. **Documentation**: Extensive documentation is expected
6. **Real Data Testing**: Always test with real database, not just fixtures

### Commit Message Style
**User prefers**: Short, one-sentence commit messages
```
feat: add feature description
fix: fix bug description
test: add tests for feature
docs: update documentation
```

## 🔨 Common Development Commands

```bash
# Development
make dev              # Start in fixture mode
make prod            # Start in database mode
make test            # Run tests
make test-cover      # Tests with coverage
make lint            # Run linter
make fmt             # Format code

# Database
make db-up           # Start PostgreSQL
make db-down         # Stop PostgreSQL
make db-reset        # Reset database (WARNING: deletes data)
make db-logs         # View database logs

# Data Import
make import-data     # Import US data
```

## 📦 Dependencies

### Go Modules

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/go-chi/chi/v5` | 5.2.3 | HTTP router |
| `github.com/go-chi/cors` | 1.2.1 | CORS middleware |
| `github.com/spf13/viper` | latest | Configuration |
| `go.uber.org/zap` | latest | Structured logging |
| `github.com/jackc/pgx/v5` | latest | PostgreSQL driver |

### Adding New Dependencies
```bash
cd backend
go get github.com/package/name@version
go mod tidy
```

## 🎨 Code Style Guidelines

### Go Style
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting (run `make fmt`)
- Prefer explicit over implicit
- Keep functions small and focused
- Use meaningful variable names

### SQL Style
- Use CTEs for complex queries
- Capitalize SQL keywords (SELECT, FROM, WHERE)
- Indent subqueries
- Comment complex logic
- Use prepared statement parameters ($1, $2, etc.)

### Error Messages
- Lowercase, no ending punctuation
- Include context: `"failed to parse year: %w"`
- Use `%w` for error wrapping
- Return proper HTTP status codes

## 🔐 Security Considerations

### Current Security Features
- CORS configured for specific frontend URL
- SQL injection prevented (parameterized queries)
- Input validation on all parameters
- Connection pooling with limits

### Future Security Enhancements
- Add rate limiting
- Add authentication/authorization
- Add request size limits
- Add SQL query timeouts
- Add HTTPS in production

## 📈 Performance Characteristics

### Query Performance
- Simple queries (meta): < 10ms
- Complex queries (names list): 100ms - 2s
- With 100M+ records, expect 2-5s for complex filters
- Indexes are crucial (already created)

### Optimization Opportunities
1. Add materialized views for name_stats
2. Add Redis caching layer
3. Add query result caching
4. Optimize CTE chains
5. Add database read replicas

## 🚢 Deployment

### Local Development
```bash
docker-compose up -d         # Start database
cd backend && make dev       # Start server (fixture mode)
```

### Production
```bash
docker-compose -f docker-compose.prod.yml up --build
```

### GitHub Actions
- **Test workflow**: Runs on all pushes/PRs
- **Deploy workflow**: Runs on push to main or releases
- **Builds**: Docker images pushed to GHCR

## 📝 When to Update Documentation

### Update AGENT.md when:
- Adding major features
- Changing architecture
- Adding new common tasks
- Discovering important gotchas

### Update architecture docs when:
- Changing API contract
- Modifying database schema
- Altering system design

### Update specific .md files when:
- DATABASE.md: Database changes
- TESTING.md: New test patterns
- DEPLOYMENT.md: Deployment process changes
- DATA_IMPORT.md: New countries or import methods
- LOGGING.md: Logging configuration changes

## 🤝 Working with the User

### Communication Style
- **Be direct and technical** (no "Great!", "Certainly!", etc.)
- **Provide specific examples** with code
- **Show actual data** in responses
- **Use markdown formatting** for readability
- **Link to files** with relative paths

### When User Says:
- "let's go further" → Continue to next logical step
- "commit message?" → Provide short, one-sentence format
- "test it" → Actually run tests and show results
- "make it work" → Debug until it actually works, don't assume

### Important:
- **Always test with real data**, not just fixtures
- **Verify changes work** before marking complete
- **Show actual output** from commands
- **Update todo lists** as you progress

## 🎓 Learning Resources

### Project Documentation
1. Start with [`ARCHITECTURE.md`](ARCHITECTURE.md)
2. Read [`architecture/02-backend-carcass.md`](architecture/02-backend-carcass.md)
3. Review [`architecture/01-shared-contract.md`](architecture/01-shared-contract.md)

### External Resources
- [Chi Router Docs](https://github.com/go-chi/chi)
- [Zap Logger Docs](https://github.com/uber-go/zap)
- [pgx Driver Docs](https://github.com/jackc/pgx)
- [Viper Config Docs](https://github.com/spf13/viper)

## ✅ Pre-flight Checklist for Agents

Before starting work:
- [ ] Read this AGENT.md file
- [ ] Check project status (what's working, what's not)
- [ ] Understand current database state
- [ ] Know which mode is active (fixture vs database)
- [ ] Review recent changes (git log)
- [ ] Check for running processes (ports in use)
- [ ] Read relevant architecture docs

Before completing work:
- [ ] All tests pass (`make test`)
- [ ] Manual testing completed
- [ ] No compilation errors
- [ ] Logs are clean
- [ ] Documentation updated if needed
- [ ] Ready for commit (provide message)

## 🎯 Current Project State

As of last update:
- **Status**: Production-ready
- **Database**: 144 years US data, 370M+ records
- **All endpoints**: Working correctly
- **Tests**: 20 passing
- **CI/CD**: Fully configured
- **Documentation**: Complete (6 guides)
- **Known issues**: None

## 📞 Quick Reference

### Important File Paths
```
backend/cmd/server/main.go           - Server entry point
backend/internal/db/queries.go       - All SQL queries
backend/internal/handlers/names.go   - Names endpoints
backend/internal/config/config.go    - Configuration
migrations/001_initial_schema.sql    - Database schema
.github/workflows/test.yml           - CI pipeline
```

### Important Commands
```bash
make test        # Run tests
make dev         # Development server
make prod        # Production server
make db-up       # Start database
make import-data # Import US data
```

### Important URLs
```
Health:     http://localhost:8080/health
Meta Years: http://localhost:8080/api/meta/years
Names:      http://localhost:8080/api/names?top_n=10
```

---

**This document should be your first stop when working on this project. Happy coding!** 🚀