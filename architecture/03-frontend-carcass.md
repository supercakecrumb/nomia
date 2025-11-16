# Frontend Carcass

The frontend carcass defines the page structure, routing, state management, API client layer, and UI component skeleton.

## Technology Stack

For complete version specifications and rationale, see the **[Technology Stack section in Overview](00-overview.md#technology-stack)**.

**Quick Reference:**
- **Node.js**: **24.11.1** (LTS)
- **Bundler**: Vite **7.2.2**
- **Language**: TypeScript **5.9.3**
- **Framework**: React **19.2.0** + react-dom **19.2.0**
- **Routing**: react-router **7.9.6** + react-router-dom **7.9.6**
- **Data Fetching**: @tanstack/react-query **5.90.9**
- **Styling**: tailwindcss **4.1.17**
- **Charts**: tremor **3.18.7**
- **State Management**: React Context + hooks (for global UI state)
- **Internationalization**: i18next **23.18.3** + react-i18next **16.0.5** (see [`05-cross-cutting-concerns.md`](05-cross-cutting-concerns.md#internationalization-i18n) for details)

## Routing & Page Structure

**Routes:**

| Route | Component | Purpose |
|-------|-----------|---------|
| `/` | MainPage | Motivation and entry point |
| `/names` | NamesExplorerPage | Name exploration table with filters |
| `/name/:name` | NameDetailPage | Detailed view for a single name |

**Query Parameter Strategy:**

Filter state is synchronized with URL query parameters to enable:
- Shareable links with specific filter configurations.
- Browser back/forward navigation.
- Bookmarking of searches.

**Example URL:**
```
/names?year_from=1980&year_to=2020&countries=US,UK&gender_balance_min=40&gender_balance_max=60&name_glob=alex*&sort_key=total_count&sort_order=desc&page=1
```

**Query Parameter Mapping:**

| URL Param | State Field | Type |
|-----------|-------------|------|
| `year_from` | `yearFrom` | number |
| `year_to` | `yearTo` | number |
| `countries` | `countries` | string[] |
| `gender_balance_min` | `genderBalanceMin` | number |
| `gender_balance_max` | `genderBalanceMax` | number |
| `min_count` | `minCount` | number |
| `top_n` | `topN` | number |
| `coverage_percent` | `coveragePercent` | number |
| `name_glob` | `nameGlob` | string |
| `sort_key` | `sortKey` | string |
| `sort_order` | `sortOrder` | "asc" \| "desc" |
| `page` | `page` | number |
| `page_size` | `pageSize` | number |

**Validation:**
- `page_size`: Min 10, Max 100, Default 50
- `page`: Min 1, Max 100 (pagination limit)

---

## Global UI State & Filters

**Filter State Model:**

```typescript
interface FilterState {
  // Year range
  yearFrom: number;
  yearTo: number;
  
  // Countries
  countries: string[]; // Array of country codes
  
  // Gender balance
  genderBalanceMin: number; // 0-100
  genderBalanceMax: number; // 0-100
  
  // Popularity (only one is "driver" at a time)
  minCount: number | null;
  topN: number | null;
  coveragePercent: number | null;
  popularityDriver: "minCount" | "topN" | "coverage" | null;
  
  // Name glob
  nameGlob: string;
  
  // Sorting
  sortKey: "popularity" | "total_count" | "name" | "gender_balance" | "countries";
  sortOrder: "asc" | "desc";
  
  // Pagination
  page: number;
  pageSize: number;
  
  // UI settings (not in URL)
  genderTintEnabled: boolean;
}
```

**State Management:**

- Use React Context or a lightweight state management solution.
- Provide a `FilterContext` that wraps the app.
- Expose:
  - `filterState` – current filter values.
  - `updateFilter(key, value)` – update a single filter.
  - `resetFilters()` – reset to defaults.
  - `applyFilters()` – trigger API call (debounced).

**URL Synchronization:**

- On filter change:
  1. Update local state.
  2. Debounce (300–500ms).
  3. Update URL query params.
  4. Trigger API call.
- On page load or URL change:
  1. Parse query params.
  2. Initialize filter state.
  3. Trigger API call.

**Debouncing:**

- Use a debounce hook or utility to avoid excessive API calls.
- Debounce applies to text inputs (name glob) and sliders (year range, gender balance).
- Dropdowns and checkboxes trigger immediately.

---

## API Client Layer

**Structure:**

Create an API client module at `/src/api/` with functions for each endpoint:

```typescript
// src/api/client.ts

export async function fetchMetaYears(): Promise<MetaYearsResponse> {
  // In mock mode: return fixture
  // In real mode: fetch from /api/meta/years
}

export async function fetchCountries(): Promise<CountriesResponse> {
  // In mock mode: return fixture
  // In real mode: fetch from /api/meta/countries
}

export async function fetchNames(params: NamesQueryParams): Promise<NamesListResponse> {
  // In mock mode: return fixture
  // In real mode: fetch from /api/names with query params
}

export async function fetchNameTrend(params: NameTrendQueryParams): Promise<NameTrendResponse> {
  // In mock mode: return fixture
  // In real mode: fetch from /api/names/trend with query params
}
```

**Mock vs Real Mode:**

Use an environment variable:

```
VITE_API_MODE=mock   # Use fixtures
VITE_API_MODE=real   # Use real API
```

**Mock Mode Implementation:**

```typescript
// src/api/fixtures.ts

import metaYearsFixture from '../../spec-examples/meta-years.json';
import countriesFixture from '../../spec-examples/countries.json';
import namesListFixture from '../../spec-examples/names-list.json';
import nameDetailFixture from '../../spec-examples/name-detail.json';

export const fixtures = {
  metaYears: metaYearsFixture,
  countries: countriesFixture,
  namesList: namesListFixture,
  nameDetail: nameDetailFixture,
};
```

```typescript
// src/api/client.ts

import { fixtures } from './fixtures';

const API_MODE = import.meta.env.VITE_API_MODE || 'mock';
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export async function fetchMetaYears(): Promise<MetaYearsResponse> {
  if (API_MODE === 'mock') {
    // Simulate network delay
    await delay(200);
    return fixtures.metaYears;
  }
  
  const response = await fetch(`${API_BASE_URL}/api/meta/years`);
  if (!response.ok) throw new Error('Failed to fetch meta years');
  return response.json();
}
```

**Integration with React Query:**

```typescript
// src/hooks/useMetaYears.ts

import { useQuery } from '@tanstack/react-query';
import { fetchMetaYears } from '../api/client';

export function useMetaYears() {
  return useQuery({
    queryKey: ['meta', 'years'],
    queryFn: fetchMetaYears,
    staleTime: Infinity, // Meta data rarely changes
  });
}
```

**Benefits:**

- Frontend can develop with fixtures immediately.
- Switching to real API only requires changing `VITE_API_MODE` environment variable.
- No changes to UI components.

---

## UI Component Skeleton

### 1. App Shell / Layout

**Responsibilities:**
- Global header with logo/name and navigation.
- Main content area (React Router outlet).
- Optional global settings (gender tint toggle, theme).

**Structure:**
```
<AppShell>
  <Header>
    <Logo />
    <Navigation />
    <LanguageSwitcher /> {/* i18n language selector */}
    <SettingsMenu />
  </Header>
  <MainContent>
    <Outlet /> {/* React Router renders pages here */}
  </MainContent>
</AppShell>
```

**i18n Integration:**
- Wrap app with i18next provider in `main.tsx`
- All text content uses translation keys (see [`05-cross-cutting-concerns.md`](05-cross-cutting-concerns.md#internationalization-i18n))
- Language preference persisted in localStorage
- Automatic language detection from browser on first visit

---

### 2. Main / Motivation Page (`/`)

**Responsibilities:**
- Explain the mission and purpose of Affirm Name.
- Highlight trans/nonbinary focus.
- Provide clear CTA to `/names`.

**Components:**
- `HeroSection` – headline and tagline.
- `MissionStatement` – why the project exists.
- `FeaturesOverview` – key capabilities (filters, charts, data sources).
- `CTAButton` – "Explore Names" → navigates to `/names`.

**Data Requirements:**
- None (static content).

---

### 3. Filter Bar (on `/names`)

**Responsibilities:**
- Display and manage all filter controls.
- Sync with global filter state.
- Trigger debounced API calls on change.

**Sub-Components:**

#### Year Range Filter
- Dual-range slider (e.g., using a library like `rc-slider`).
- Numeric input fields for precise control.
- Displays `yearFrom` and `yearTo`.
- Bounds set by `/api/meta/years` response.

#### Countries Filter
- Multi-select dropdown with checkboxes.
- Shows country codes and names.
- Displays data source info in tooltip (source name, URL).
- Fetches options from `/api/meta/countries`.

#### Gender Balance Filter
- Range slider on 0–100 axis.
- Labels: "More Female" (0) ← → "More Male" (100).
- Shows current bounds (e.g., "40–60").
- Visual indicator for "unisex" range (around 50).

#### Popularity Filter Trio
- Three input controls:
  - **Min Count**: Numeric input (e.g., "at least 500 people").
  - **Top N**: Numeric input (e.g., "top 1000 names").
  - **Coverage Percent**: Numeric input or slider (e.g., "top 95%").
- Only one is "active driver" at a time.
- When user changes one, the other two are derived from API response.
- Display logic:
  - Show all three values.
  - Highlight the active driver.
  - Update derived values after API response.

#### Name Glob Filter
- Text input with placeholder: "e.g., alex* or *сан*".
- Help tooltip explaining glob syntax (`*` and `?`).
- Debounced input (300–500ms).

#### Sort Controls
- Dropdown for `sort_key`:
  - Options: "Popularity", "Total Count", "Name", "Gender Balance", "Countries".
- Toggle button for `sort_order`: "Ascending" / "Descending".

**Layout:**
- Horizontal bar at top of page (sticky or fixed).
- Collapsible on mobile.
- Clear visual grouping of related filters.

---

### 4. Names Table

**Responsibilities:**
- Display paginated list of names.
- Show key metrics per name.
- Enable navigation to detail page.
- Handle loading, empty, and error states.

**Columns:**

| Column | Content | Interaction |
|--------|---------|-------------|
| **Name** | Name string | Clickable → navigate to `/name/:name` |
| **Popularity** | Rank + percentile label (e.g., "Top 5%") | Sortable |
| **Total Count** | Formatted number (e.g., "125,430") | Sortable |
| **Gender Balance** | Visual bar with dot + percentages | Sortable |
| **Presence Period** | Formatted period (e.g., "1975–2010") | Non-sortable |
| **Countries** | Flags (max 3) + "+N" if more | Sortable by count |

**Gender Balance Column:**
- Horizontal bar representing 0–100 axis.
- Dot positioned at `gender_balance` value.
- Left side: female percentage (e.g., "45%").
- Right side: male percentage (e.g., "55%").
- Optional: background tint based on `genderTintEnabled` setting.

**Row Interaction:**
- Click row or name → navigate to `/name/:name?<current_filters>`.
- Preserve current filters in URL for context.

**States:**
- **Loading**: Skeleton rows (shimmer effect).
- **Empty**: Message "No names match your filters. Try adjusting your criteria."
- **Error**: Error message with retry button.

**Pagination:**
- Controls at bottom: "Previous", "Next", page number input.
- Display: "Showing 1–50 of 1,523 names".

---

### 5. Name Detail Page (`/name/:name`)

**Responsibilities:**
- Display detailed information for a single name.
- Show time-series and country-level charts.
- Allow filtering by year range and countries (subset of main filters).
- Provide back link to `/names` with preserved filters.

**Components:**

#### Summary Header
- **Name**: Large, prominent display.
- **Total Count**: Across selected filters.
- **Gender Balance**: Visual bar + percentages.
- **Presence Period**: Formatted period.
- **Countries**: List of country codes/names.

#### Filter Bar (Subset)
- Year range slider.
- Countries multi-select.
- Apply button (or auto-apply with debounce).

#### Charts

**Popularity Over Time Chart:**
- Line chart: X-axis = year, Y-axis = total count.
- Shows trend of name usage over time.
- Optionally: separate lines for female/male counts.

**Gender Distribution Over Time Chart:**
- Stacked area chart or line chart.
- X-axis = year, Y-axis = count or percentage.
- Two series: female count, male count.
- Shows how gender balance shifts over time.

**Country Breakdown Chart:**
- Horizontal bar chart.
- X-axis = total count, Y-axis = country name.
- Shows which countries contribute most to this name's usage.

**Chart Library:**
- Use Recharts or similar declarative library.
- Ensure charts are responsive and accessible.

#### Back Link
- Button or link: "← Back to Name Explorer".
- Preserves filters from URL query params.

**Data Source:**
- Fetches from `/api/names/trend?name=<name>&year_from=<yearFrom>&year_to=<yearTo>&countries=<countries>`.

---

### 6. Settings / Preferences

**Responsibilities:**
- Toggle `genderTintEnabled` setting.
- Stored in local state or localStorage.
- Affects table row background tint only.

**UI:**
- Checkbox or toggle switch.
- Label: "Enable gender tint in table rows".
- Accessible via settings menu in header or as inline control near table.

---

## Mocking Strategy

**Option 1: Mock Service Worker (MSW)**

- Use MSW to intercept `/api/*` requests in the browser.
- Return fixture data for each endpoint.
- Enables realistic network behavior (delays, errors).

**Setup:**
```typescript
// src/mocks/handlers.ts

import { http, HttpResponse } from 'msw';
import { fixtures } from '../api/fixtures';

export const handlers = [
  http.get('/api/meta/years', () => {
    return HttpResponse.json(fixtures.metaYears);
  }),
  
  http.get('/api/meta/countries', () => {
    return HttpResponse.json(fixtures.countries);
  }),
  
  http.get('/api/names', () => {
    return HttpResponse.json(fixtures.namesList);
  }),
  
  http.get('/api/names/trend', () => {
    return HttpResponse.json(fixtures.nameDetail);
  }),
];
```

**Option 2: Direct Fixture Import**

- API client directly imports and returns fixture JSON.
- Simpler setup, no additional libraries.
- Less realistic (no network delay, no error simulation).

**Recommendation:**
- Start with Option 2 (direct import) for simplicity.
- Migrate to Option 1 (MSW) if realistic network behavior is needed.

---

## Internationalization Components

See comprehensive i18n implementation in [`05-cross-cutting-concerns.md`](05-cross-cutting-concerns.md#internationalization-i18n).

### Language Switcher

**Location:** App header (see [App Shell](#1-app-shell--layout))

**Component:** `LanguageSwitcher`

**Responsibilities:**
- Display available languages (English, Russian)
- Indicate current language
- Switch language on click
- Persist preference to localStorage

**Visual Design:**
- Flag emojis + language labels
- Active state highlighting
- Accessible with keyboard navigation
- Mobile-friendly dropdown on small screens

### Translation Integration

**All UI Components Must:**
1. Import `useTranslation` hook from `react-i18next`
2. Use translation keys for all user-facing text
3. Support dynamic content with interpolation
4. Handle pluralization for counts

**Example:**
```typescript
import { useTranslation } from 'react-i18next';

function FilterBar() {
  const { t } = useTranslation('filters');
  
  return (
    <label>{t('year_range.label')}</label>
  );
}
```

### Translation File Organization

**Namespaces:**
- `common`: Navigation, buttons, general UI
- `filters`: All filter-related strings
- `pages`: Page-specific content
- `errors`: Error messages
- `validation`: Form validation messages

**File Location:** `public/locales/{lang}/{namespace}.json`

**Supported Languages:** English (`en`), Russian (`ru`)

---

[← Previous: Backend Carcass](02-backend-carcass.md) | [Next: Integration Strategy →](04-integration-strategy.md)