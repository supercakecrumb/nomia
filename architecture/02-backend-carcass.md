# Backend Carcass

The backend carcass defines the database schema, HTTP handler structure, ingestion skeleton, and how to start with fixtures and transition to real implementations.

## Technology Stack

For complete version specifications and rationale, see the **[Technology Stack section in Overview](00-overview.md#technology-stack)**.

**Quick Reference:**
- **Language**: Go **1.25.4**
- **Database**: PostgreSQL **18.1**
- **HTTP**: `net/http` + `github.com/go-chi/chi/v5@v5.2.3`
- **Database Driver**: `github.com/jackc/pgx/v5@v5.7.6` (primary)
- **Migrations**: `github.com/golang-migrate/migrate/v4@v4.19.0`
- **Configuration**: `github.com/spf13/viper@v1.21.0`
- **Logging**: `log/slog` over sugaring of `go.uber.org/zap@v1.27.0`

## Database Schema

### Table: `countries`

Stores metadata about countries and their data sources.

**Columns:**
- `id` (serial, primary key)
- `code` (varchar(10), unique, not null) – e.g., "US", "UK"
- `name` (varchar(255), not null) – e.g., "United States"
- `data_source_name` (varchar(255), not null)
- `data_source_url` (text, not null)
- `data_source_description` (text, nullable)
- `data_source_requires_manual_download` (boolean, default true)
- `created_at` (timestamp, default now)
- `updated_at` (timestamp, default now)

**Indexes:**
- Primary key on `id`.
- Unique index on `code`.

**Purpose:**
- Drives country filter options in UI.
- Provides data provenance information.
- Links to datasets and names.

---

### Table: `name_datasets`

Represents an uploaded dataset file from a given source.

**Columns:**
- `id` (serial, primary key)
- `country_id` (integer, foreign key → `countries.id`, not null)
- `source_file_name` (varchar(255), not null) – original filename
- `source_url` (text, nullable) – original download URL
- `year_from` (integer, nullable) – earliest year covered by this dataset
- `year_to` (integer, nullable) – latest year covered by this dataset
- `file_type` (varchar(50), not null) – e.g., "csv", "tsv", "xlsx"
- `storage_path` (text, not null) – local path or object storage key
- `parser_version` (varchar(50), nullable) – version of parser used
- `parse_status` (varchar(50), not null) – "uploaded", "parsed", "failed"
- `uploaded_at` (timestamp, default now)
- `parsed_at` (timestamp, nullable)
- `error_message` (text, nullable) – if parse_status = "failed"

**Indexes:**
- Primary key on `id`.
- Index on `country_id`.
- Index on `parse_status`.

**Purpose:**
- Trace origin of each batch of data.
- Support re-ingestion, debugging, and auditing.
- Enable dataset management UI (future).

---

### Table: `names`

Core fact table storing atomic name records.

**Columns:**
- `id` (bigserial, primary key)
- `country_id` (integer, foreign key → `countries.id`, not null)
- `dataset_id` (integer, foreign key → `name_datasets.id`, not null)
- `year` (integer, not null)
- `name` (varchar(255), not null)
- `gender` (char(1), not null) – 'M', 'F', or 'U' (unknown)
- `count` (integer, not null)

**Indexes:**
- Primary key on `id`.
- Composite index on `(country_id, year, name, gender)` – for filtering and aggregation.
- GIN index on `name` using `pg_trgm` extension – for glob matching:
  ```sql
  CREATE EXTENSION IF NOT EXISTS pg_trgm;
  CREATE INDEX idx_names_name_trgm ON names USING GIN (name gin_trgm_ops);
  ```
- Index on `dataset_id` – for dataset-level queries.

**Purpose:**
- Raw material for all aggregation queries.
- Enables flexible filtering by country, year, name, and gender.

---

### Optional: Precomputed Aggregates

For performance, consider a materialized view or aggregate table:

**Table: `name_stats` (optional)**

Precomputed aggregates for common queries.

**Columns:**
- `name` (varchar(255), primary key)
- `total_count` (bigint)
- `female_count` (bigint)
- `male_count` (bigint)
- `gender_balance` (numeric)
- `earliest_year` (integer)
- `latest_year` (integer)
- `country_codes` (text[]) – array of country codes

**Refresh Strategy:**
- Rebuild after each dataset ingestion.
- Or use PostgreSQL materialized views with `REFRESH MATERIALIZED VIEW`.

**Note:** This is optional and can be added later for optimization. Start with direct queries on `names` table.

---

## HTTP Handlers & Routing

**Routing Structure:**

```
/api/meta/years       → MetaYearsHandler
/api/meta/countries   → MetaCountriesHandler
/api/names            → NamesListHandler
/api/names/trend      → NameTrendHandler
/api/datasets/upload  → DatasetUploadHandler (admin/ingestion)
```

**Handler Skeleton (Conceptual):**

Each handler follows this pattern:

1. **Parse and validate query parameters.**
2. **Log request details** (for debugging).
3. **In fixture mode**: Return contents of corresponding JSON file from `/spec-examples/`.
4. **In real mode**: Execute database query, compute metrics, return JSON response.
5. **Handle errors** with standard error response format.

**Example: NamesListHandler (Pseudocode)**

```
func NamesListHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Parse query params
    params := parseNamesListParams(r.URL.Query())
    
    // 2. Validate params
    if err := validateParams(params); err != nil {
        writeError(w, 400, "invalid_params", err.Error())
        return
    }
    
    // 3. Check mode
    if config.FixtureMode {
        // Return fixture
        data := loadFixture("spec-examples/names-list.json")
        writeJSON(w, 200, data)
        return
    }
    
    // 4. Execute real query
    names, meta, err := queryNames(params)
    if err != nil {
        writeError(w, 500, "query_failed", err.Error())
        return
    }
    
    // 5. Return response
    response := map[string]interface{}{
        "meta": meta,
        "names": names,
    }
    writeJSON(w, 200, response)
}
```

**Handler Organization:**

- Group handlers by feature:
  - `handlers/meta.go` – meta endpoints
  - `handlers/names.go` – names list and trend
  - `handlers/datasets.go` – ingestion endpoints
- Shared utilities:
  - `handlers/params.go` – parameter parsing and validation
  - `handlers/response.go` – JSON response helpers
  - `handlers/errors.go` – error response helpers

---

## Filter & Popularity Pipeline (Conceptual)

The `/api/names` endpoint follows this logical pipeline:

### Stage 1: Filter

Apply filters to the `names` table:

1. **Year range**: `year >= year_from AND year <= year_to`
2. **Countries**: `country_id IN (selected_country_ids)` (union semantics)
3. **Name glob**: `name ILIKE pattern` (case-insensitive, using trigram index)
4. **Optional raw min_count**: Applied after aggregation (see Stage 2)

### Stage 2: Aggregation

For all rows passing Stage 1 filters:

1. **Group by name** (and optionally country if needed for country list).
2. **Compute per-name metrics**:
   - `total_count` = SUM(count)
   - `female_count` = SUM(count WHERE gender = 'F')
   - `male_count` = SUM(count WHERE gender = 'M')
   - `gender_balance` = 100 × (male_count / (male_count + female_count))
   - `name_start` = MIN(year)
   - `name_end` = MAX(year)
   - `countries` = ARRAY_AGG(DISTINCT country_code)

### Stage 3: Gender Balance Filter

Apply gender balance filter:

- `gender_balance >= gender_balance_min AND gender_balance <= gender_balance_max`

### Stage 4: Popularity Computation

1. **Sort** aggregated names by `total_count` DESC.
2. **Compute**:
   - `rank` = row number in sorted list
   - `cumulative_count` = running sum of `total_count`
   - `cumulative_share` = `cumulative_count / total_count_in_filtered_set`

### Stage 5: Popularity Filter

Apply effective popularity cut (one of):

- `total_count >= min_count`
- `rank <= top_n`
- `cumulative_share <= coverage_percent / 100`

### Stage 6: Sorting & Pagination

1. **Sort** by specified `sort_key` and `sort_order`.
2. **Apply tie-breaking rules**:
   - Secondary: `total_count` DESC
   - Tertiary: `name` ASC
3. **Paginate**: `LIMIT page_size OFFSET (page - 1) * page_size`

**Implementation Notes:**

- Stages 1–3 can be done in a single SQL query with CTEs (Common Table Expressions).
- Stage 4 (popularity computation) may require window functions or application-level processing.
- Stage 5 (popularity filter) is applied after computing rank/cumulative_share.
- Stage 6 (sorting/pagination) is the final SQL step.

---

## Dataset Ingestion Carcass

### Upload Endpoint Skeleton

**POST /api/datasets/upload**

**Purpose:** Accept a dataset file upload and queue it for parsing.

**Request:**
- Multipart form data with:
  - `file` (file upload)
  - `country_id` (integer)
  - `source_url` (optional string)
  - `notes` (optional string)

**Response:**
```json
{
  "dataset_id": 123,
  "status": "uploaded",
  "message": "Dataset uploaded successfully. Parsing will begin shortly."
}
```

**Handler Logic:**

1. Validate `country_id` exists in `countries` table.
2. Save uploaded file to storage (local filesystem or object storage).
3. Insert row into `name_datasets`:
   - `country_id`, `source_file_name`, `source_url`, `file_type`, `storage_path`
   - `parse_status` = "uploaded"
4. Enqueue background job: `parse_dataset(dataset_id)`.
5. Return response with `dataset_id`.

---

### Worker Skeleton

**Background Job: `parse_dataset(dataset_id)`**

**Purpose:** Parse an uploaded dataset and populate the `names` table.

**Logic:**

1. Fetch `name_datasets` row by `dataset_id`.
2. Determine which parser to use:
   - Based on `country_id` (e.g., US → SSA parser, UK → ONS parser).
   - Or based on `file_type` (e.g., CSV → generic CSV parser).
3. Load file from `storage_path`.
4. Parse file:
   - Extract rows with (year, name, gender, count).
   - Validate data (e.g., year is numeric, count is positive).
5. Insert rows into `names` table:
   - Batch insert for performance.
   - Set `country_id`, `dataset_id`, `year`, `name`, `gender`, `count`.
6. On success:
   - Update `name_datasets`: `parse_status` = "parsed", `parsed_at` = now.
7. On failure:
   - Update `name_datasets`: `parse_status` = "failed", `error_message` = error details.
   - Log error for debugging.

**Worker Implementation:**

- Use a job queue (e.g., Redis + worker pool, or simple Go channels).
- Or use a cron job that polls for `parse_status = 'uploaded'` rows.

---

### Parser Abstraction

Define a parser interface:

```
type DatasetParser interface {
    Parse(filePath string) ([]NameRecord, error)
}

type NameRecord struct {
    Year   int
    Name   string
    Gender string // "M", "F", or "U"
    Count  int
}
```

**Concrete Parsers (Stubs for Now):**

- `SSAParser` – parses US SSA CSV format.
- `ONSParser` – parses UK ONS format.
- `GenericCSVParser` – parses generic CSV with configurable columns.

**Parser Selection Logic:**

```
func selectParser(countryID int, fileType string) DatasetParser {
    switch countryID {
    case 1: // US
        return &SSAParser{}
    case 2: // UK
        return &ONSParser{}
    default:
        return &GenericCSVParser{}
    }
}
```

**Note:** Parsers are stubs initially. Implement them incrementally as datasets are added.

---

## Mock vs Real Mode

**Configuration:**

Use an environment variable or config flag:

```
FIXTURE_MODE=true   # Use fixtures
FIXTURE_MODE=false  # Use real database
```

**Handler Behavior:**

```
if config.FixtureMode {
    // Load and return fixture JSON
    data := loadFixture("spec-examples/names-list.json")
    writeJSON(w, 200, data)
} else {
    // Execute real database query
    names, meta, err := queryNames(params)
    // ...
}
```

**Benefits:**

- Backend can start by implementing handlers that return fixtures.
- Frontend can develop against fixtures without waiting for backend.
- Later, backend switches to real mode by setting `FIXTURE_MODE=false`.
- Contract remains unchanged.

---

[← Previous: Shared Contract](01-shared-contract.md) | [Next: Frontend Carcass →](03-frontend-carcass.md)