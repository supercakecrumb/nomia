# Shared Contract

The shared contract is the **core of the carcass**. It defines the API surface, data semantics, and terminology that both frontend and backend implement against.

## Terminology & Core Concepts

### Gender Balance Axis (0–100)

A single numeric axis representing the ratio between female and male usage:

- **0** = 100% female, 0% male (fully female)
- **50** = 50% female, 50% male (perfectly balanced/unisex)
- **100** = 0% female, 100% male (fully male)

**Calculation:**
```
If (male_count + female_count) == 0:
    gender_balance = NULL  // No binary gender data
Else:
    gender_balance = 100 × (male_count / (male_count + female_count))
```

Where:
- `female_count` = total occurrences recorded as female
- `male_count` = total occurrences recorded as male
- `unknown_count` = total occurrences recorded as unknown/nonbinary

**Rationale:** This spectrum model avoids binary gender assumptions and makes it easy to filter for unisex names (values near 50) or names with specific gender associations.

**Handling Unknown/Nonbinary Data:**

The API exposes three gender count fields:
- `female_count` - count of female occurrences
- `male_count` - count of male occurrences
- `unknown_count` - count of unknown/nonbinary occurrences

Additional fields:
- `gender_balance` - 0–100 axis value (NULL if no binary gender data)
- `has_unknown_data` - boolean, true if unknown_count > 0

**Display Rules:**
- Show gender balance bar only when `gender_balance` is not NULL
- When `has_unknown_data` is true, show indicator (e.g., badge "Includes nonbinary/unknown data")
- In detail view, show full breakdown: "Female: 45% | Male: 55% | Unknown: 2%"

### Popularity Metrics

Within a filtered set of names (after applying year, country, gender balance, and glob filters), popularity is computed as follows:

1. **Sort** all names by `total_count` (descending).
2. For each name, compute:
   - **rank**: Position in the sorted list (1 = most frequent).
   - **cumulative_count**: Sum of counts from rank 1 up to and including this name.
   - **cumulative_share**: `cumulative_count / total_count_in_filtered_set` (expressed as 0–1 or percentage).

These metrics enable three user-friendly popularity filters:

- **Min Total Count**: Numeric threshold (e.g., "at least 500 people").
- **Top N**: Keep only names with rank ≤ N (e.g., "top 1000 names").
- **Coverage Percentile**: Keep names while `cumulative_share ≤ threshold` (e.g., "top 95% of people").

**Important:** These three filters are different expressions of the same underlying cut in the popularity distribution. The frontend treats only one as the "driver" at any time and derives the other two from API responses.

**Precedence When Multiple Filters Provided:**

If the backend receives multiple popularity filters, it applies them in this priority order:

1. **`coverage_percent`** (highest priority) - If provided and > 0, use it (ignore others)
2. **`top_n`** - If provided and > 0 (and no coverage_percent), use it (ignore min_count)
3. **`min_count`** (lowest priority) - If provided and > 0 (and no other filters)
4. **None** - If no filters provided or all are 0/null, return all names (no popularity filter)

**Frontend Behavior:**
- Track which filter was last changed by user (`popularityDriver` in state)
- Send only the active filter to backend in API request
- Derive inactive filter values from `popularity_summary` in API response (see GET /api/names response below)

### Presence Period

For each name, define:

- **name_start**: Earliest year in which the name appears in the filtered data.
- **name_end**: Latest year in which the name appears in the filtered data.

Also define global bounds:

- **db_start**: Global earliest year across all ingested datasets (from `/api/meta/years`).
- **db_end**: Global latest year across all ingested datasets.

**Display Semantics:**

| Condition | Display Format | Example |
|-----------|----------------|---------|
| `name_start > db_start` AND `name_end < db_end` | `start–end` | 1975–2010 |
| `name_start == db_start` AND `name_end < db_end` | `–end` | –2010 |
| `name_start > db_start` AND `name_end == db_end` | `start–` | 1995– |
| `name_start == db_start` AND `name_end == db_end` | `–` | – |

The backend exposes raw values (`name_start`, `name_end`, `db_start`, `db_end`); the frontend applies these formatting rules.

### Name Glob Filter

A glob-based pattern filter for name matching:

- **Query Parameter**: `name_glob`
- **Pattern Syntax**:
  - `*` matches any sequence of characters (including empty).
  - `?` matches any single character.
- **Matching**:
  - Case-insensitive.
  - Implicitly anchored (pattern must match the entire name).
- **Examples**:
  - `alex*` → matches Alex, Alexander, Alexis, Alexandra, etc.
  - `*сан*` → matches any name containing "сан" (Cyrillic).
  - `a?ex` → matches Alex, Apex, etc.

**Backend Implementation Strategy:**

- Use case-insensitive matching (SQL `ILIKE` or equivalent).
- Add a **pg_trgm GIN index** on the `name` column in PostgreSQL to optimize glob-like patterns.
- The glob condition is applied alongside other filters (year, country, gender balance).
- Popularity metrics (rank, cumulative share) are computed only over the set of names that pass all filters, including `name_glob`.

## API Endpoints

### 1. GET /api/meta/years

**Purpose:** Returns the global minimum and maximum year across all ingested datasets.

**Query Parameters:** None.

**Response:**
```json
{
  "min_year": 1880,
  "max_year": 2023
}
```

**Semantics:**
- `min_year`: Earliest year with data in the system.
- `max_year`: Latest year with data in the system.
- Used by frontend to set default year range filter bounds.

---

### 2. GET /api/meta/countries

**Purpose:** Returns all countries known to the system with their metadata.

**Query Parameters:** None.

**Response:**
```json
{
  "countries": [
    {
      "code": "US",
      "name": "United States",
      "data_source_name": "Social Security Administration",
      "data_source_url": "https://www.ssa.gov/oact/babynames/",
      "data_source_description": "Annual baby name data from SSA records.",
      "data_source_requires_manual_download": true
    },
    {
      "code": "UK",
      "name": "United Kingdom",
      "data_source_name": "Office for National Statistics",
      "data_source_url": "https://www.ons.gov.uk/",
      "data_source_description": null,
      "data_source_requires_manual_download": true
    }
  ]
}
```

**Field Semantics:**
- `code`: ISO-like country code (e.g., "US", "UK", "SE").
- `name`: Full country name.
- `data_source_name`: Name of the statistical agency or source.
- `data_source_url`: Canonical URL for the data source.
- `data_source_description`: Optional free-text explanation.
- `data_source_requires_manual_download`: Boolean indicating if datasets must be manually downloaded.

**Usage:**
- Populates country filter dropdown.
- Displays data provenance in UI (tooltips, detail pages).

---

### 3. GET /api/names

**Purpose:** Core exploration endpoint. Returns a paginated list of names with metrics, filtered and sorted according to query parameters.

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `year_from` | integer | No | `db_start` | Lower bound of year range (inclusive). |
| `year_to` | integer | No | `db_end` | Upper bound of year range (inclusive). |
| `countries` | string | No | all | Comma-separated list of country codes (e.g., "US,UK,SE"). Union semantics: include names from any of these countries. |
| `gender_balance_min` | integer | No | 0 | Minimum gender balance (0–100). |
| `gender_balance_max` | integer | No | 100 | Maximum gender balance (0–100). |
| `min_count` | integer | No | 0 | Minimum total count threshold. |
| `top_n` | integer | No | null | Keep only names with rank ≤ N. |
| `coverage_percent` | float | No | null | Keep names while cumulative_share ≤ threshold (0–100). |
| `name_glob` | string | No | empty | Glob pattern for name matching (case-insensitive). |
| `sort_key` | string | No | "popularity" | Sort field: "popularity", "total_count", "name", "gender_balance", "countries". |
| `sort_order` | string | No | "asc" | Sort order: "asc" or "desc". |
| `page` | integer | No | 1 | Page number (1-based). |
| `page_size` | integer | No | 50 | Number of results per page (min 10, max 100). |

**Parameter Validation:**
- `year_from` must be <= `year_to`
- `gender_balance_min` must be <= `gender_balance_max`
- `page` must be >= 1 and <= 100 (pagination limit)
- `page_size` must be >= 10 and <= 100

**Filter Interaction:**
- All filters are combined with logical AND.
- Popularity metrics (rank, cumulative_share) are computed after applying all other filters.
- Only one of `min_count`, `top_n`, or `coverage_percent` should be "active" at a time (frontend responsibility to manage this).

**Sorting:**
- **Tie-breaking rules** (stable sort):
  1. Primary: specified `sort_key`.
  2. Secondary: `total_count` (desc).
  3. Tertiary: `name` (asc, lexicographic).

**Response:**
```json
{
  "meta": {
    "page": 1,
    "page_size": 50,
    "total_count": 1523,
    "total_pages": 31,
    "db_start": 1880,
    "db_end": 2023,
    "popularity_summary": {
      "population_total": 1530245,
      "active_driver": "coverage_percent",
      "active_value": 95.0,
      "derived_min_count": 487,
      "derived_top_n": 1200,
      "derived_coverage_percent": 95.0
    }
  },
  "names": [
    {
      "name": "Alex",
      "total_count": 125430,
      "female_count": 62715,
      "male_count": 62715,
      "gender_balance": 50.0,
      "rank": 1,
      "cumulative_share": 0.082,
      "name_start": 1920,
      "name_end": 2023,
      "countries": ["US", "UK", "CA"]
    },
    {
      "name": "Jordan",
      "total_count": 98234,
      "female_count": 45000,
      "male_count": 53234,
      "gender_balance": 54.2,
      "rank": 2,
      "cumulative_share": 0.146,
      "name_start": 1880,
      "name_end": 2023,
      "countries": ["US", "UK"]
    }
  ]
}
```

**Field Semantics:**
- `meta.total_count`: Total number of names matching filters (before pagination).
- `meta.total_pages`: Total pages available.
- `meta.db_start`, `meta.db_end`: Global year bounds (for presence period formatting).
- `name`: The given name.
- `total_count`: Sum of occurrences across selected countries and years.
- `female_count`, `male_count`: Gender-specific counts.
- `gender_balance`: 0–100 axis value.
- `rank`: Position in popularity ranking (1 = most popular).
- `cumulative_share`: Fraction of total population covered up to this name (0–1).
- `name_start`, `name_end`: Earliest and latest year this name appears in filtered data.
- `countries`: Array of country codes where this name appears.

---

### 4. GET /api/names/trend

**Purpose:** Returns detailed time-series and country-level data for a single name.

**Query Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `name` | string | Yes | - | The name to retrieve details for. |
| `year_from` | integer | No | `db_start` | Lower bound of year range. |
| `year_to` | integer | No | `db_end` | Upper bound of year range. |
| `countries` | string | No | all | Comma-separated list of country codes. |

**Response:**
```json
{
  "name": "Alex",
  "meta": {
    "db_start": 1880,
    "db_end": 2023
  },
  "summary": {
    "total_count": 125430,
    "female_count": 62715,
    "male_count": 62715,
    "gender_balance": 50.0,
    "name_start": 1920,
    "name_end": 2023,
    "countries": ["US", "UK", "CA"]
  },
  "time_series": [
    {
      "year": 1920,
      "total_count": 450,
      "female_count": 200,
      "male_count": 250,
      "gender_balance": 55.6
    },
    {
      "year": 1921,
      "total_count": 480,
      "female_count": 240,
      "male_count": 240,
      "gender_balance": 50.0
    }
  ],
  "by_country": [
    {
      "country_code": "US",
      "country_name": "United States",
      "total_count": 98000,
      "female_count": 49000,
      "male_count": 49000,
      "gender_balance": 50.0
    },
    {
      "country_code": "UK",
      "country_name": "United Kingdom",
      "total_count": 20000,
      "female_count": 10000,
      "male_count": 10000,
      "gender_balance": 50.0
    }
  ]
}
```

**Field Semantics:**
- `summary`: Aggregated metrics for the name across all selected years and countries.
- `time_series`: Year-by-year breakdown (only years with data).
- `by_country`: Country-level breakdown.

---

## JSON Fixtures

To enable parallel development, the contract is exemplified by JSON fixture files stored in `/spec-examples/`:

| File | Purpose |
|------|---------|
| `meta-years.json` | Example response for `/api/meta/years`. |
| `countries.json` | Example response for `/api/meta/countries`. |
| `names-list.json` | Example response for `/api/names` (with various filter scenarios). |
| `name-detail.json` | Example response for `/api/names/trend`. |

**Usage:**
- **Backend**: Validates real responses against these fixtures to ensure contract compliance.
- **Frontend**: Uses these fixtures directly during development (before backend is ready).

**Fixture Requirements:**
- Must include all fields defined in the API contract.
- Should cover edge cases:
  - Empty results.
  - Names with presence periods at boundaries (–end, start–, –).
  - Names with gender balance at extremes (0, 25, 40, 50, 54, 100).
  - Names with unknown_count > 0 (has_unknown_data = true, gender_balance = NULL).
  - Multiple countries per name.
  - Various popularity ranks and cumulative shares.
  - Popularity summary metadata with all three derived values.

**See:** [`architecture/FIXTURE-SPECIFICATIONS.md`](FIXTURE-SPECIFICATIONS.md) for complete fixture specifications with JSON examples.

---

[← Previous: Overview](00-overview.md) | [Next: Backend Carcass →](02-backend-carcass.md)