# Testing Guide

This document describes the testing strategy and how to run tests for the Affirm Name backend.

## Test Structure

```
backend/
├── internal/
│   ├── db/
│   │   └── queries_test.go       # Parameter parsing and validation tests
│   └── handlers/
│       └── params_test.go        # Handler parameter tests
└── TESTING.md                    # This file
```

## Running Tests

### Run All Tests
```bash
cd backend
go test ./...
```

### Run Tests with Coverage
```bash
cd backend
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View coverage in browser
```

### Run Tests with Race Detector
```bash
cd backend
go test -race ./...
```

### Run Specific Package Tests
```bash
cd backend
go test ./internal/handlers/... -v
go test ./internal/db/... -v
```

### Run Specific Test
```bash
cd backend
go test ./internal/db -run TestParseNamesListParams -v
```

## Test Categories

### 1. Unit Tests (`*_test.go`)
Tests for individual functions and methods without external dependencies.

**Current Coverage:**
- ✅ Parameter parsing (`ParseNamesListParams`, `ParseNameTrendParams`)
- ✅ Parameter validation (year ranges, gender balance, pagination)
- ✅ Popularity filter logic (`GetActivePopularityFilter`)

**Examples:**
```go
func TestParseNamesListParams(t *testing.T) {
    // Test default parameters
    // Test custom parameters
    // Test validation errors
}
```

### 2. Integration Tests (Future)
Tests that interact with the database using a test database.

**Planned Coverage:**
- Database query functions (`GetYearRange`, `GetCountries`, `GetNamesList`, `GetNameTrend`)
- Database connections and connection pooling
- Transaction handling

**Example Setup:**
```go
func TestGetNamesList_Integration(t *testing.T) {
    // Setup test database
    db, cleanup := setupTestDB(t)
    defer cleanup()
    
    // Insert test data
    insertTestData(t, db)
    
    // Run query
    result, err := db.GetNamesList(ctx, params)
    
    // Assert results
    assert.NoError(t, err)
    assert.Equal(t, expectedCount, len(result.Names))
}
```

### 3. End-to-End Tests (Future)
Tests that make HTTP requests to the API endpoints.

**Planned Coverage:**
- All 4 API endpoints
- Fixture mode responses
- Database mode responses
- Error handling (404, 400, 500)

**Example:**
```go
func TestNamesListEndpoint_E2E(t *testing.T) {
    // Start test server
    server := httptest.NewServer(handler)
    defer server.Close()
    
    // Make request
    resp, err := http.Get(server.URL + "/api/names?page=1")
    
    // Assert response
    assert.Equal(t, 200, resp.StatusCode)
}
```

## Test Data

### Fixture Data
Located in [`../spec-examples/`](../spec-examples/):
- [`meta-years.json`](../spec-examples/meta-years.json)
- [`countries.json`](../spec-examples/countries.json)
- [`names-list.json`](../spec-examples/names-list.json)
- [`name-detail.json`](../spec-examples/name-detail.json)

### Database Test Data (Future)
Will be located in `testdata/`:
- SQL seed files for test scenarios
- CSV files for bulk insert tests

## CI/CD Integration

Tests run automatically on:
- Every push to `main` or `develop` branches
- Every pull request
- GitHub Actions workflow: [`.github/workflows/test.yml`](../.github/workflows/test.yml)

### CI Steps:
1. ✅ Setup PostgreSQL test database
2. ✅ Run database migrations
3. ✅ Execute all tests with race detector
4. ✅ Check code formatting (`gofmt`)
5. ✅ Run static analysis (`go vet`)
6. ✅ Build application
7. ✅ Upload coverage report
8. ✅ Run linter (`golangci-lint`)

## Coverage Goals

| Package | Current | Target |
|---------|---------|--------|
| handlers | ~50% | 80% |
| db | ~30% | 70% |
| config | 0% | 60% |
| middleware | 0% | 60% |

## Writing New Tests

### Test Naming Convention
- Test files: `*_test.go`
- Test functions: `Test<FunctionName>(t *testing.T)`
- Subtests: `t.Run("descriptive_name", func(t *testing.T) { ... })`

### Table-Driven Tests
Use table-driven tests for multiple test cases:

```go
tests := []struct {
    name    string
    input   string
    want    int
    wantErr bool
}{
    {"valid input", "123", 123, false},
    {"invalid input", "abc", 0, true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := ParseInt(tt.input)
        if tt.wantErr {
            assert.Error(t, err)
            return
        }
        assert.NoError(t, err)
        assert.Equal(t, tt.want, got)
    })
}
```

### Test Helpers
Common test utilities should be in `testing.go` or package-specific helpers.

## Debugging Tests

### View Test Output
```bash
go test -v ./...
```

### Run Single Test with Details
```bash
go test -v -run TestSpecificFunction ./internal/handlers
```

### Debug with Delve
```bash
dlv test ./internal/handlers -- -test.run TestSpecificFunction
```

## Mocking (Future)

For database and external service mocking, we'll use:
- [`gomock`](https://github.com/golang/mock) for interface mocking
- [`httptest`](https://pkg.go.dev/net/http/httptest) for HTTP testing
- [`testcontainers-go`](https://github.com/testcontainers/testcontainers-go) for integration tests

## Performance Testing (Future)

### Benchmark Tests
```go
func BenchmarkGetNamesList(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Run function
    }
}
```

### Run Benchmarks
```bash
go test -bench=. -benchmem ./...
```

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Testing Best Practices](https://golang.org/doc/effective_go#testing)