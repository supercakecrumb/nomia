# Architecture Refinement Plan

**Created:** 2025-11-15  
**Status:** In Progress  
**Priority:** Critical for Development Start

This document addresses all feedback from the architecture review and provides a systematic plan for refinements before development begins.

---

## Executive Summary

**Overall Assessment:** A- (Excellent with refinements needed)

**Critical Issues:** 8 items (must fix before development)  
**Important Issues:** 5 items (address week 1)  
**Moderate Issues:** 5 items (address weeks 2-3)

**Timeline:**
- **Critical Fixes:** Complete before development kickoff (Days 1-2)
- **Week 1 Fixes:** Complete during foundation phase (Days 3-7)
- **Week 2-3 Fixes:** Complete during core features phase (Days 8-21)

---

## Critical Issues (Days 1-2, Blocking)

### 1. Missing Fixture Files ⚠️ BLOCKING

**Issue:** No `/spec-examples/*.json` files exist. Frontend and backend cannot start parallel development without them.

**Impact:** Blocks all development work.

**Solution:**

Create four fixture files with comprehensive edge cases:

**Files to Create:**
1. `spec-examples/meta-years.json`
2. `spec-examples/countries.json`
3. `spec-examples/names-list.json`
4. `spec-examples/name-detail.json`

**Requirements:**
- Include edge cases: empty results, boundary years, gender extremes (0, 50, 100)
- Cover multiple countries per name
- Show various popularity ranks and cumulative shares
- Include names with different presence periods (–end, start–, –)

**Owner:** Architecture team  
**Due:** Before kickoff (Day 1)  
**Status:** 🔴 Not Started

---

### 2. Popularity Filter Logic Unclear

**Issue:** When multiple popularity filters (`min_count`, `top_n`, `coverage_percent`) are provided, precedence is undefined.

**Current State:** Contract says "only one should be active" but doesn't specify what happens if multiple are sent.

**Solution:**

**Define Precedence Order:**
1. `coverage_percent` (highest priority)
2. `top_n`
3. `min_count` (lowest priority)

**Backend Behavior:**
- If `coverage_percent` is provided and > 0, use it (ignore others).
- Else if `top_n` is provided and > 0, use it (ignore `min_count`).
- Else if `min_count` is provided and > 0, use it.
- Else apply no popularity filter (return all names).

**Frontend Behavior:**
- Track which filter was last changed by user (`popularityDriver`).
- Send only the active filter to backend.
- Derive other two from API response metadata (see Issue #2 in teammate feedback).

**Documentation Updates:**
- Add to `architecture/01-shared-contract.md` under "Popularity Metrics"
- Add to `architecture/03-frontend-carcass.md` under "Popularity Filter Trio"

**Owner:** Architecture team  
**Due:** Day 1  
**Status:** 🔴 Not Started

---

### 3. Database Index Specifications Incomplete

**Issue:** Index descriptions are conceptual. Need exact SQL with column order for performance.

**Solution:**

Create `migrations/001_initial_schema.sql` with:

```sql
-- Enable pg_trgm extension
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Countries table
CREATE TABLE countries (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    data_source_name VARCHAR(255) NOT NULL,
    data_source_url TEXT NOT NULL,
    data_source_description TEXT,
    data_source_requires_manual_download BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_countries_code ON countries(code);

-- Name datasets table
CREATE TABLE name_datasets (
    id SERIAL PRIMARY KEY,
    country_id INTEGER NOT NULL REFERENCES countries(id),
    source_file_name VARCHAR(255) NOT NULL,
    source_url TEXT,
    year_from INTEGER,
    year_to INTEGER,
    file_type VARCHAR(50) NOT NULL,
    storage_path TEXT NOT NULL,
    parser_version VARCHAR(50),
    parse_status VARCHAR(50) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT NOW(),
    parsed_at TIMESTAMP,
    error_message TEXT
);

CREATE INDEX idx_datasets_country ON name_datasets(country_id);
CREATE INDEX idx_datasets_status ON name_datasets(parse_status);

-- Names table (core fact table)
CREATE TABLE names (
    id BIGSERIAL PRIMARY KEY,
    country_id INTEGER NOT NULL REFERENCES countries(id),
    dataset_id INTEGER NOT NULL REFERENCES name_datasets(id),
    year INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    gender CHAR(1) NOT NULL CHECK (gender IN ('M', 'F', 'U')),
    count INTEGER NOT NULL CHECK (count > 0)
);

-- Critical indexes for filter pipeline performance
-- Stage 1: Year + Country + Name filtering
CREATE INDEX idx_names_filter ON names(country_id, year, name, gender);

-- Stage 1: Glob matching on name
CREATE INDEX idx_names_name_trgm ON names USING GIN (name gin_trgm_ops);

-- Dataset queries
CREATE INDEX idx_names_dataset ON names(dataset_id);

-- Optional: Covering index for aggregation without table lookup
-- CREATE INDEX idx_names_aggregation ON names(country_id, year, name, gender) INCLUDE (count);
```

**Documentation Updates:**
- Add migration file to repository
- Update `architecture/02-backend-carcass.md` with exact SQL
- Add performance testing criteria (p95 < 500ms for 1M+ rows)

**Owner:** Backend lead  
**Due:** Day 1  
**Status:** 🔴 Not Started

---

### 4. Gender Axis Handling for Nonbinary/Unknown Data ⚠️ CRITICAL

**Issue:** Gender balance formula only considers male/female counts, but schema stores `gender='U'` rows. Unclear how unknown/nonbinary records are handled.

**Current Formula:**
```
gender_balance = 100 × (male_count / (male_count + female_count))
```

**Problems:**
- Ignores `gender='U'` rows entirely
- Divide-by-zero when both male_count and female_count are 0
- Not inclusive of nonbinary data

**Solution:**

**Option A: Expose Unknown Count Separately (Recommended)**

1. **Extend API Response:**
```json
{
  "name": "Alex",
  "total_count": 125430,
  "female_count": 62715,
  "male_count": 62715,
  "unknown_count": 0,
  "gender_balance": 50.0,
  "has_unknown_data": false
}
```

2. **Gender Balance Calculation:**
```
If (male_count + female_count) == 0:
    gender_balance = NULL  // No binary gender data
    has_unknown_data = true
Else:
    gender_balance = 100 × (male_count / (male_count + female_count))
    has_unknown_data = (unknown_count > 0)
```

3. **UI Display:**
- Show gender balance bar only when `gender_balance` is not NULL
- When `has_unknown_data` is true, show badge: "Includes nonbinary/unknown data"
- In detail view, show breakdown: "Female: 45% | Male: 55% | Unknown: 2%"

**Option B: Three-Way Split (Alternative)**

Not recommended for initial release. Adds complexity to axis visualization.

**Decision:** Use Option A (expose unknown_count, keep binary axis simple).

**Documentation Updates:**
- Update `architecture/01-shared-contract.md` - Gender Balance Axis section
- Update `architecture/02-backend-carcass.md` - Aggregation Stage logic
- Update `architecture/03-frontend-carcass.md` - Gender Balance Column display

**Owner:** Architecture team + Frontend lead  
**Due:** Day 2  
**Status:** 🔴 Not Started

---

### 5. Popularity Summary Metadata Missing

**Issue:** Frontend needs to derive inactive popularity filter values from API response, but paginated results don't include enough metadata to compute correct cutoffs.

**Current Problem:**
- API returns per-name ranks within a page
- Frontend can't infer global cutoffs (e.g., "what min_count gives same results as top_n=1000?")
- Pagination changes reported thresholds

**Solution:**

Add `popularity_summary` to `/api/names` response metadata:

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
  "names": [...]
}
```

**Field Semantics:**
- `population_total`: Total count of all people in filtered dataset (before popularity cut)
- `active_driver`: Which filter was applied ("min_count" | "top_n" | "coverage_percent" | null)
- `active_value`: Value of the active filter
- `derived_min_count`: Minimum count that would give same results
- `derived_top_n`: Top N rank that would give same results
- `derived_coverage_percent`: Coverage percentage achieved

**Frontend Usage:**
- Display all three values in UI
- Highlight the active driver
- When user changes one, send to backend, receive updated derived values

**Documentation Updates:**
- Update `architecture/01-shared-contract.md` - GET /api/names response
- Update `architecture/03-frontend-carcass.md` - Popularity Filter Trio logic

**Owner:** Backend lead + Frontend lead  
**Due:** Day 2  
**Status:** 🔴 Not Started

---

### 6. Pagination Parameters Not Synced

**Issue:** API exposes `page_size`, filter state tracks `pageSize`, but URL mapping omits it. Shareable links won't preserve page size.

**Solution:**

1. **Add to URL Mapping:**
```typescript
| URL Param | State Field | Type | Default |
|-----------|-------------|------|---------|
| `page_size` | `pageSize` | number | 50 |
```

2. **Validation:**
- Min: 10
- Max: 100
- Default: 50
- If invalid, use default and log warning

3. **Update Documentation:**
- `architecture/01-shared-contract.md` - add validation rules
- `architecture/03-frontend-carcass.md` - add to URL mapping table

**Owner:** Frontend lead  
**Due:** Day 2  
**Status:** 🔴 Not Started

---

### 7. Dataset Upload Security Missing

**Issue:** POST `/api/datasets/upload` has no authentication, file size limits, or deduplication.

**Solution:**

**Authentication:**
- Use JWT tokens for admin authentication
- Require `Authorization: Bearer <token>` header
- Return 401 if missing/invalid

**File Validation:**
- Max file size: 100MB
- Allowed types: CSV, TSV, XLSX
- Reject with 413 (Payload Too Large) if exceeds limit

**Deduplication:**
- Compute SHA-256 checksum of uploaded file
- Check if `name_datasets` has row with same:
  - `country_id`
  - `year_from` and `year_to` (if provided)
  - File checksum
- If duplicate found, return 409 (Conflict) with existing `dataset_id`

**Idempotency:**
- Worker uses `dataset_id` as idempotency key
- If worker crashes, re-running with same `dataset_id` is safe
- Check `parse_status` before starting: if "parsed", skip; if "uploaded", process; if "failed", retry

**Audit Logging:**
- Log all upload attempts with:
  - User ID
  - Timestamp
  - File name
  - File size
  - Checksum
  - Success/failure
  - Error message (if failed)

**Documentation Updates:**
- Update `architecture/02-backend-carcass.md` - Dataset Ingestion section
- Update `architecture/05-cross-cutting-concerns.md` - Security section
- Add new section in `02-backend-carcass.md` for authentication

**Owner:** Backend lead  
**Due:** Day 2  
**Status:** 🔴 Not Started

---

### 8. Error Handling HTTP Status Codes Incomplete

**Issue:** Missing status codes for common scenarios: rate limiting, validation errors, conflicts.

**Solution:**

**Complete Error Code Table:**

| Code | HTTP Status | Meaning | Example |
|------|-------------|---------|---------|
| `invalid_params` | 400 | Invalid query parameters | `year_from > year_to` |
| `invalid_glob` | 400 | Invalid glob pattern | `name_glob="[invalid"` |
| `validation_error` | 400 | Request validation failed | Missing required field |
| `unauthorized` | 401 | Missing or invalid auth token | No `Authorization` header |
| `forbidden` | 403 | Insufficient permissions | Non-admin accessing upload |
| `not_found` | 404 | Resource not found | Name doesn't exist |
| `conflict` | 409 | Duplicate resource | Dataset already uploaded |
| `payload_too_large` | 413 | File too large | Upload > 100MB |
| `too_many_requests` | 429 | Rate limit exceeded | > 100 req/min |
| `internal_error` | 500 | Server error | Unexpected exception |
| `database_error` | 500 | Database query failed | Connection lost |
| `parse_error` | 500 | Dataset parsing failed | Invalid CSV format |

**Rate Limiting Headers:**
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1605564000
Retry-After: 60
```

**Documentation Updates:**
- Update `architecture/05-cross-cutting-concerns.md` - Error Handling section

**Owner:** Backend lead  
**Due:** Day 2  
**Status:** 🔴 Not Started

---

## Important Issues (Week 1, Days 3-7)

### 9. Name Normalization Strategy Missing

**Issue:** No strategy for handling name variants (José vs Jose, café vs cafe).

**Solution:**

**Storage Strategy:**
- Store names in their original form (no normalization)
- Apply Unicode NFC normalization only for display consistency

**Search Strategy:**
- Glob filter searches original form (case-insensitive)
- Consider adding optional "normalized search" toggle in future

**Sorting Strategy:**
- Use PostgreSQL collation for locale-aware sorting
- Default: `COLLATE "en_US.UTF-8"`
- Allow per-country collation in future

**API Changes:**
- No changes to current contract
- Document that names are stored and searched as-is

**Documentation Updates:**
- Add new section to `architecture/01-shared-contract.md` - "Name Normalization"
- Update `architecture/05-cross-cutting-concerns.md` - i18n section

**Owner:** Backend lead  
**Due:** Week 1  
**Status:** 🟡 Pending

---

### 10. Gender Field Mapping Guidance Missing

**Issue:** Parsers need guidance on mapping source data to M/F/U.

**Solution:**

**Parser Guidelines:**

1. **US SSA:** Maps directly (M, F) → (M, F)
2. **UK ONS:** Maps directly (Male, Female) → (M, F)
3. **Generic CSV:**
   - Recognize variations: "M", "Male", "m", "Männlich" → M
   - Recognize variations: "F", "Female", "f", "Weiblich" → F
   - Unknown/other/blank → U

**When to Use 'U':**
- Source data explicitly marks as "Unknown", "Other", "Nonbinary"
- Source data has blank/null gender field
- Source data has unrecognized gender value

**Quality Assurance:**
- Log warnings when mapping to 'U'
- Include count of U-mapped rows in dataset stats
- Allow admin review before finalizing ingestion

**Documentation Updates:**
- Add to `architecture/02-backend-carcass.md` - Parser Abstraction section

**Owner:** Backend lead  
**Due:** Week 1  
**Status:** 🟡 Pending

---

### 11. Pagination Performance Won't Scale

**Issue:** OFFSET-based pagination performs poorly beyond page 100 (requires scanning all previous rows).

**Solution:**

**Short-term (Phase 1-2):**
- Keep OFFSET-based pagination
- Add warning in docs: "Limit pagination to first 100 pages"
- Limit `page` parameter: max value = 100

**Long-term (Phase 3-4):**
- Implement cursor-based pagination
- Use keyset pagination with composite key: `(sort_key_value, name)`
- API changes:
  - Add `cursor` parameter (base64-encoded)
  - Remove `page` parameter
  - Keep `page_size` parameter

**Cursor Format:**
```json
{
  "cursor_type": "keyset",
  "sort_key": "total_count",
  "sort_order": "desc",
  "last_value": 125430,
  "last_name": "Alex"
}
```

**Documentation Updates:**
- Add to `architecture/01-shared-contract.md` - Pagination section
- Add to `architecture/02-backend-carcass.md` - Performance considerations

**Owner:** Backend lead  
**Due:** Week 1 (short-term), Week 6 (long-term)  
**Status:** 🟡 Pending

---

### 12. Upgrade Fixture/Mocking Approach

**Issue:** Direct fixture imports ignore query params. Phase 2 filter testing requires param-aware mocking.

**Solution:**

**Option A: Mock Service Worker (Recommended)**

1. **Install MSW:**
```bash
npm install msw --save-dev
```

2. **Create Handlers:**
```typescript
// src/mocks/handlers.ts
import { http, HttpResponse } from 'msw';

export const handlers = [
  http.get('/api/names', ({ request }) => {
    const url = new URL(request.url);
    const yearFrom = url.searchParams.get('year_from');
    const nameGlob = url.searchParams.get('name_glob');
    
    // Filter fixture data based on params
    let names = namesListFixture.names;
    
    if (nameGlob) {
      const pattern = globToRegex(nameGlob);
      names = names.filter(n => pattern.test(n.name));
    }
    
    return HttpResponse.json({
      meta: { ...namesListFixture.meta, total_count: names.length },
      names: names.slice(0, 50)
    });
  }),
];
```

3. **Start Worker:**
```typescript
// src/mocks/browser.ts
import { setupWorker } from 'msw/browser';
import { handlers } from './handlers';

export const worker = setupWorker(...handlers);
```

4. **Initialize in Dev:**
```typescript
// src/main.tsx
if (import.meta.env.VITE_API_MODE === 'mock') {
  const { worker } = await import('./mocks/browser');
  await worker.start();
}
```

**Option B: Lightweight Mock Server**

- Create Express server that parses query params
- Run on `localhost:8081`
- Frontend points to mock server URL

**Decision:** Use Option A (MSW) for Phase 2+.

**Documentation Updates:**
- Update `architecture/03-frontend-carcass.md` - Mocking Strategy section
- Update `architecture/06-development-workflow.md` - Phase 2 setup

**Owner:** Frontend lead  
**Due:** Week 1 (before Phase 2)  
**Status:** 🟡 Pending

---

### 13. Workflow Milestones Need Verifiable Exit Criteria

**Issue:** Phase checklists are descriptive, not verifiable. No clear definition of done.

**Solution:**

**Phase 1: Foundation - Exit Criteria:**
- [ ] Backend serves fixture responses for all 4 endpoints
- [ ] Fixture responses validate against JSON schema
- [ ] Frontend navigates between all 3 pages without errors
- [ ] Fixtures cover edge cases: empty results, boundary years, gender extremes
- [ ] CI pipeline runs linter and tests successfully
- [ ] Health check endpoint returns 200

**Phase 2: Core Features - Exit Criteria:**
- [ ] Filter interactions update URL query params correctly
- [ ] Debouncing works: text inputs trigger API after 500ms idle
- [ ] Table displays paginated results with loading skeletons
- [ ] Dataset upload returns 201 with dataset_id
- [ ] Worker skeleton processes test dataset without errors
- [ ] MSW handlers parse query params and filter fixture data

**Phase 3: Visualization & Data - Exit Criteria:**
- [ ] `/api/names` returns real data from database
- [ ] Query performance: p95 < 500ms for 1M rows
- [ ] Charts render correctly with real data
- [ ] Gender balance visualization shows correct percentages
- [ ] Presence period formatting handles all 4 cases (–end, start–, –, start–end)
- [ ] Accessibility: keyboard navigation works, ARIA labels present

**Phase 4: Integration & Polish - Exit Criteria:**
- [ ] Frontend works correctly against real backend
- [ ] E2E tests pass: search, filter, view details, pagination
- [ ] Accessibility tests pass (axe-core, manual screen reader testing)
- [ ] Performance budget met: Lighthouse score > 90
- [ ] Bundle size: initial < 500KB, per-route < 200KB
- [ ] Error handling: all error codes tested, retry mechanism works
- [ ] Rate limiting: returns 429 after limit exceeded

**Documentation Updates:**
- Update `architecture/06-development-workflow.md` - add exit criteria for each phase

**Owner:** Project manager + Tech leads  
**Due:** Week 1  
**Status:** 🟡 Pending

---

## Moderate Issues (Weeks 2-3, Days 8-21)

### 14. CORS Configuration Too Permissive

**Issue:** CORS section mentions whitelist but doesn't specify production configuration.

**Solution:**

**Development:**
```go
cors.AllowOrigins([]string{"http://localhost:5173", "http://localhost:3000"})
```

**Staging:**
```go
cors.AllowOrigins([]string{"https://staging.affirm-name.com"})
```

**Production:**
```go
cors.AllowOrigins([]string{"https://affirm-name.com", "https://www.affirm-name.com"})
```

**Configuration:**
- Use environment variable: `CORS_ORIGINS="https://affirm-name.com,https://www.affirm-name.com"`
- Validate on startup
- Log rejected requests

**Documentation Updates:**
- Update `architecture/05-cross-cutting-concerns.md` - Security section

**Owner:** Backend lead  
**Due:** Week 2  
**Status:** 🟡 Pending

---

### 15. Data Retention Strategy Missing

**Issue:** No cleanup policy for old or failed datasets.

**Solution:**

**Retention Policy:**
- Keep successful datasets indefinitely
- Keep failed datasets for 30 days (for debugging)
- Archive uploaded files after ingestion (move to cold storage)

**Cleanup Job:**
- Run daily cron job
- Delete `name_datasets` rows where:
  - `parse_status = 'failed'`
  - `uploaded_at < NOW() - INTERVAL '30 days'`
- Delete associated files from storage

**Audit:**
- Log all deletions with:
  - Dataset ID
  - Reason (retention policy)
  - Deletion timestamp

**Documentation Updates:**
- Add to `architecture/02-backend-carcass.md` - Dataset Management section

**Owner:** Backend lead  
**Due:** Week 2  
**Status:** 🟡 Pending

---

### 16. Frontend State Management Complex

**Issue:** Filter state + URL sync + debouncing adds complexity. Consider URL-as-single-source-of-truth.

**Solution:**

**Evaluation Needed:**
- Current approach: State in Context + sync to URL
- Alternative: URL is source of truth, parse on every render

**Pros of URL-as-source:**
- Simpler: no sync logic
- Shareable by default
- No state management bugs

**Cons of URL-as-source:**
- More URL parsing
- Harder to implement debouncing
- Transient UI state needs separate handling

**Decision:** Keep current approach for Phase 1-2. Re-evaluate in Phase 3 if complexity becomes issue.

**Documentation Updates:**
- Add note to `architecture/03-frontend-carcass.md` - State Management section

**Owner:** Frontend lead  
**Due:** Week 3 (evaluation)  
**Status:** 🟡 Pending

---

### 17. Chart Library Unvalidated

**Issue:** Recharts recommended but not validated for accessibility and bundle size.

**Solution:**

**Validation Criteria:**
- Accessibility: keyboard nav, ARIA labels, screen reader support
- Bundle size: < 100KB gzipped
- Features: line charts, area charts, bar charts
- TypeScript support

**Alternatives to Evaluate:**
- Recharts (current recommendation)
- Chart.js + react-chartjs-2
- Victory
- Nivo

**Testing:**
- Create proof-of-concept with each library
- Test accessibility with axe-core
- Measure bundle size impact
- Evaluate API ergonomics

**Documentation Updates:**
- Update `architecture/03-frontend-carcass.md` - Charts section with chosen library

**Owner:** Frontend lead  
**Due:** Week 2  
**Status:** 🟡 Pending

---

### 18. Internationalization Requirements Vague

**Issue:** UTF-8 handling mentioned but specific requirements unclear.

**Solution:**

**Immediate Requirements (Phase 1-3):**
- Database: UTF-8 encoding (`CREATE DATABASE affirm_name WITH ENCODING 'UTF8'`)
- API: UTF-8 Content-Type headers (`Content-Type: application/json; charset=utf-8`)
- Frontend: UTF-8 meta tag (`<meta charset="utf-8">`)
- Test with names from multiple scripts: Latin, Cyrillic, Arabic, Chinese, emoji

**Testing Checklist:**
- [ ] Store and retrieve names with diacritics (José, Müller, Zoë)
- [ ] Store and retrieve Cyrillic names (Александр, Мария)
- [ ] Store and retrieve Arabic names (محمد, فاطمة)
- [ ] Store and retrieve Chinese names (王伟, 李娜)
- [ ] Display emoji in names if present (rare but should work)
- [ ] Glob filter works with non-Latin characters

**Future (Post-Launch):**
- UI localization (Spanish, French, German)
- Right-to-left (RTL) layout support
- Locale-aware number formatting

**Documentation Updates:**
- Update `architecture/05-cross-cutting-concerns.md` - i18n section

**Owner:** Backend lead + Frontend lead  
**Due:** Week 3  
**Status:** 🟡 Pending

---

## Additional Recommendations

### 19. Add OpenAPI Specification

**Purpose:** Machine-readable API contract for validation and code generation.

**Action:**
- Create `openapi.yaml` with full API spec
- Use for contract testing
- Generate TypeScript types from spec

**Owner:** Backend lead  
**Due:** Week 2  
**Status:** 🟡 Pending

---

### 20. Add System Architecture Diagram

**Purpose:** Visual reference for system components and data flow.

**Action:**
- Create Mermaid diagram showing:
  - External systems (statistical agencies)
  - Backend components (API, ingestion, database)
  - Frontend components
  - Data flow

**Owner:** Architecture team  
**Due:** Week 2  
**Status:** 🟡 Pending

---

### 21. Add Health Check Endpoints

**Purpose:** Enable monitoring and orchestration.

**Endpoints:**
- `GET /health` - Simple liveness check (returns 200 if server running)
- `GET /health/ready` - Readiness check (returns 200 if database connected)

**Owner:** Backend lead  
**Due:** Week 1  
**Status:** 🟡 Pending

---

### 22. Add Performance Budgets Section

**Purpose:** Define performance targets for monitoring.

**Metrics:**
- API response time: p50 < 200ms, p95 < 500ms, p99 < 1000ms
- Frontend: LCP < 2.5s, FID < 100ms, CLS < 0.1
- Bundle size: initial < 500KB, per-route < 200KB

**Owner:** Tech leads  
**Due:** Week 2  
**Status:** 🟡 Pending

---

## Implementation Timeline

### Pre-Development (Days 1-2)

**Day 1:**
- [ ] Create all fixture files with edge cases
- [ ] Define popularity filter precedence
- [ ] Create migration file with exact index SQL
- [ ] Document gender axis handling for unknown data

**Day 2:**
- [ ] Add popularity summary metadata to contract
- [ ] Add page_size to URL mapping
- [ ] Document dataset upload security requirements
- [ ] Complete error code table

### Week 1 (Days 3-7)

- [ ] Document name normalization strategy
- [ ] Add gender field mapping guidelines
- [ ] Document pagination performance limits
- [ ] Set up MSW for filter testing
- [ ] Add verifiable exit criteria to each phase
- [ ] Add health check endpoints

### Weeks 2-3 (Days 8-21)

- [ ] Configure CORS for all environments
- [ ] Implement data retention policy
- [ ] Evaluate frontend state management
- [ ] Validate chart library choice
- [ ] Test UTF-8 handling with multiple scripts
- [ ] Create OpenAPI specification
- [ ] Add system architecture diagram
- [ ] Add performance budgets

---

## Success Criteria

**Ready for Development Kickoff:**
- ✅ All critical issues (1-8) resolved
- ✅ Fixture files created and validated
- ✅ Migration files with indexes committed
- ✅ Contract clarifications documented

**Ready for Phase 1:**
- ✅ All Week 1 issues resolved
- ✅ Exit criteria defined for all phases
- ✅ MSW set up for frontend mocking

**Ready for Production:**
- ✅ All issues (critical + important + moderate) resolved
- ✅ Performance budgets met
- ✅ Security requirements implemented
- ✅ Accessibility requirements met

---

## Tracking

Use this checklist format in issue tracker:

```
[CRITICAL] Issue #1: Missing Fixture Files
Priority: Blocking
Owner: Architecture team
Due: Day 1
Status: 🔴 Not Started / 🟡 In Progress / 🟢 Complete
```

---

**Next Step:** Begin implementing critical fixes (Issues 1-8) before development kickoff.