# Cross-Cutting Concerns

This chapter addresses system-wide concerns that affect both frontend and backend.

## Error Handling

### Standard Error Response

**JSON Format:**

```json
{
  "error": {
    "code": "invalid_params",
    "message": "Invalid year range: year_from must be <= year_to"
  }
}
```

### Error Codes

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

When rate limit is exceeded (429), include these headers in response:

```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1605564000
Retry-After: 60
```

**Field Descriptions:**
- `X-RateLimit-Limit`: Maximum requests allowed per window
- `X-RateLimit-Remaining`: Requests remaining in current window
- `X-RateLimit-Reset`: Unix timestamp when limit resets
- `Retry-After`: Seconds to wait before retrying

### Frontend Error Display

**Inline Errors:**
- Show error message near relevant UI element (e.g., filter bar).
- Use clear, actionable language.
- Provide suggestions for fixing the error.

**Toast Notifications:**
- For transient errors (e.g., network timeout).
- Auto-dismiss after 5 seconds.
- Allow manual dismissal.

**Error Boundary:**
- Catch React errors and show fallback UI.
- Log errors to error tracking service (e.g., Sentry).
- Provide "Reload Page" button.

**Retry Mechanism:**
- Provide "Retry" button for failed API calls.
- Use exponential backoff for automatic retries.
- Limit retry attempts (max 3).

### Backend Error Logging

**Structured Logging:**
- Use JSON format for logs.
- Include context: request ID, user, timestamp, params.
- Include stack traces for debugging.

**Log Levels:**
- **DEBUG**: Detailed information for debugging.
- **INFO**: General informational messages.
- **WARN**: Warning messages (recoverable issues).
- **ERROR**: Error messages (unrecoverable issues).

**Error Monitoring:**
- Integrate with error tracking service (e.g., Sentry, Rollbar).
- Alert on critical errors.
- Track error rates and trends.

---

## Performance

### Backend Considerations

**Indexing:**
1. **Composite index** on `(country_id, year, name, gender)` for filtering.
2. **GIN trigram index** on `name` for glob matching.
3. **Index** on `dataset_id` for dataset queries.
4. **Unique index** on `countries.code`.

**Query Optimization:**
- Use CTEs (Common Table Expressions) for complex queries.
- Use window functions for popularity computation.
- Limit result set size (max `page_size` = 100).
- Use `EXPLAIN ANALYZE` to profile queries.

**Caching:**
- Cache `/api/meta/years` and `/api/meta/countries` responses.
- Use HTTP cache headers (e.g., `Cache-Control: max-age=3600`).
- Consider Redis for session/query caching.

**Rate Limiting:**
- Implement rate limiting to prevent abuse.
- Use token bucket or sliding window algorithm.
- Return 429 (Too Many Requests) when limit exceeded.

**Connection Pooling:**
- Use database connection pool.
- Set appropriate pool size (e.g., 10–20 connections).
- Monitor connection usage.

### Frontend Considerations

**Debouncing:**
- Debounce text inputs (name glob) by 300–500ms.
- Debounce sliders (year range, gender balance) by 300–500ms.
- Use `lodash.debounce` or custom hook.

**Pagination:**
- Limit table rows per page (default 50, max 100).
- Use virtual scrolling if rendering large lists.
- Consider infinite scroll for mobile.

**React Query Caching:**
- Cache API responses with appropriate `staleTime`.
- Use `keepPreviousData` option for smooth pagination.
- Prefetch next page on pagination.

**Code Splitting:**
- Lazy load chart library (Recharts) on detail page.
- Lazy load pages with `React.lazy` and `Suspense`.
- Use dynamic imports for large dependencies.

**Asset Optimization:**
- Optimize images (WebP format, responsive sizes).
- Minify CSS and JavaScript.
- Use CDN for static assets.

**Bundle Size:**
- Monitor bundle size with `webpack-bundle-analyzer`.
- Tree-shake unused code.
- Target bundle size: <500KB (initial), <200KB (per route).

---

## Accessibility

### Keyboard Navigation

**Requirements:**
- All interactive elements (buttons, links, inputs) must be keyboard accessible.
- Tab order must be logical and predictable.
- Focus indicators must be clearly visible.

**Sliders:**
- Support arrow keys for adjustment.
- Support Page Up/Down for larger increments.
- Support Home/End for min/max values.

**Table:**
- Rows must be focusable and navigable with Tab/Shift+Tab.
- Support keyboard shortcuts (e.g., Enter to open detail page).

### Screen Reader Support

**Semantic HTML:**
- Use `<table>`, `<th>`, `<td>` for tables.
- Use `<nav>`, `<main>`, `<section>` for layout.
- Use `<label>` for form inputs.

**ARIA Labels:**
- Provide `aria-label` for icon buttons and controls.
- Use `aria-describedby` for help text.
- Use `aria-live` regions for dynamic content.

**State Announcements:**
- Announce filter changes (e.g., "Year range updated to 1980–2020").
- Announce loading states (e.g., "Loading names...").
- Announce errors (e.g., "Error loading data").

### Visual Accessibility

**Color Contrast:**
- Ensure sufficient contrast (WCAG AA minimum: 4.5:1 for text, 3:1 for UI components).
- Test with contrast checker tools.
- Avoid color as sole indicator of meaning.

**Text Labels:**
- Provide text labels for all visual indicators.
- Don't rely solely on color (e.g., gender balance bar should have percentages).

**Focus Indicators:**
- Provide clear focus indicators (e.g., outline, border).
- Don't remove default focus styles without replacement.

### Chart Accessibility

**Descriptions:**
- Provide `<figcaption>` or `aria-label` to describe chart purpose.
- Provide text summary of key insights.

**Data Tables:**
- Consider providing data tables as alternative to charts.
- Allow users to toggle between chart and table views.

**Interactive Charts:**
- Ensure chart tooltips are keyboard accessible.
- Provide keyboard shortcuts for navigation.

---

## Configuration

### Backend Environment Variables

| Variable | Purpose | Example |
|----------|---------|---------|
| `FIXTURE_MODE` | Enable fixture mode | `true` / `false` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@localhost/affirm_name` |
| `PORT` | HTTP server port | `8080` |
| `STORAGE_PATH` | Dataset file storage path | `/var/data/datasets` |
| `LOG_LEVEL` | Logging verbosity | `debug` / `info` / `warn` / `error` |
| `CORS_ORIGINS` | Allowed CORS origins | `http://localhost:5173,https://affirm-name.com` |
| `RATE_LIMIT_REQUESTS` | Max requests per minute | `100` |
| `RATE_LIMIT_WINDOW` | Rate limit window (seconds) | `60` |

### Frontend Environment Variables

| Variable | Purpose | Example |
|----------|---------|---------|
| `VITE_API_MODE` | API mode (mock/real) | `mock` / `real` |
| `VITE_API_BASE_URL` | Backend API base URL | `http://localhost:8080` |
| `VITE_DEFAULT_LANGUAGE` | Default UI language | `en` |
| `VITE_SUPPORTED_LANGUAGES` | Comma-separated language codes | `en,ru` |
| `VITE_ENABLE_ANALYTICS` | Enable analytics | `true` / `false` |
| `VITE_SENTRY_DSN` | Sentry error tracking DSN | `https://...@sentry.io/...` |

### Configuration Files

**Backend:**
- Use `config.yaml` or environment variables.
- Support multiple environments (development, staging, production).
- Validate configuration on startup.

**Frontend:**
- Use `.env` files (`.env.development`, `.env.production`).
- Use Vite's environment variable support.
- Never commit sensitive values (API keys, secrets).

---

## Internationalization (i18n)

### Overview

The application supports **multi-language UI** to serve a global audience, with initial support for **English** and **Russian**, and extensibility for additional languages. The data layer already supports names from multiple languages and scripts.

### Supported Languages (Phase 1)

- **English (en)**: Default language
- **Russian (ru)**: Primary additional language

**Future Languages (Phase 2+):**
- Spanish (es), German (de), French (fr), Ukrainian (uk), etc.

---

### Frontend Implementation

#### Tech Stack

See [`00-overview.md`](00-overview.md#technology-stack) for versions:
- **i18next**: Core i18n framework
- **react-i18next**: React integration
- **i18next-browser-languagedetector**: Automatic language detection
- **i18next-http-backend**: Load translations from server

#### Translation File Structure

```
frontend/
  public/
    locales/
      en/
        common.json           # Common UI strings (navigation, buttons)
        filters.json          # Filter labels and placeholders
        pages.json            # Page-specific content
        errors.json           # Error messages
        validation.json       # Form validation messages
      ru/
        common.json
        filters.json
        pages.json
        errors.json
        validation.json
```

**File Organization Principles:**
- **Namespace by feature**: Separate translation files by UI concern
- **Flat structure**: Avoid deep nesting in JSON (max 2 levels)
- **Key naming**: Use `snake_case` for keys (e.g., `year_range_filter`, `apply_filters_button`)

**Example: `en/filters.json`**
```json
{
  "year_range": {
    "label": "Year Range",
    "from_placeholder": "From",
    "to_placeholder": "To",
    "help_text": "Filter names by the years they were popular"
  },
  "countries": {
    "label": "Countries",
    "placeholder": "Select countries...",
    "select_all": "Select All",
    "clear_all": "Clear All"
  },
  "gender_balance": {
    "label": "Gender Balance",
    "more_female": "More Female",
    "more_male": "More Male",
    "unisex_range": "Unisex (40-60%)"
  },
  "name_glob": {
    "label": "Name Pattern",
    "placeholder": "e.g., alex* or *сан*",
    "help_text": "Use * for any characters, ? for single character"
  }
}
```

#### i18next Configuration

**Setup (`src/i18n/config.ts`):**

```typescript
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import Backend from 'i18next-http-backend';

i18n
  .use(Backend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: 'en',
    supportedLngs: ['en', 'ru'],
    defaultNS: 'common',
    ns: ['common', 'filters', 'pages', 'errors', 'validation'],
    
    backend: {
      loadPath: '/locales/{{lng}}/{{ns}}.json',
    },
    
    detection: {
      order: ['localStorage', 'navigator', 'htmlTag'],
      caches: ['localStorage'],
      lookupLocalStorage: 'preferredLanguage',
    },
    
    interpolation: {
      escapeValue: false, // React already escapes
    },
    
    react: {
      useSuspense: true,
    },
  });

export default i18n;
```

#### Usage in Components

**Basic Usage:**
```typescript
import { useTranslation } from 'react-i18next';

function FilterBar() {
  const { t } = useTranslation('filters');
  
  return (
    <div>
      <label>{t('year_range.label')}</label>
      <input placeholder={t('year_range.from_placeholder')} />
      <span className="help-text">{t('year_range.help_text')}</span>
    </div>
  );
}
```

**With Interpolation:**
```typescript
const { t } = useTranslation('common');

// Translation: "Showing {{from}}-{{to}} of {{total}} names"
<p>{t('pagination.showing', { from: 1, to: 50, total: 1523 })}</p>
```

**With Pluralization:**
```typescript
// Translation keys:
// "results_count": "{{count}} result",
// "results_count_plural": "{{count}} results"
<p>{t('results_count', { count: names.length })}</p>
```

#### Language Switcher Component

**Location:** Add to app header (see [`03-frontend-carcass.md`](03-frontend-carcass.md#1-app-shell--layout))

**Component (`src/components/LanguageSwitcher.tsx`):**
```typescript
import { useTranslation } from 'react-i18next';

const LANGUAGES = [
  { code: 'en', label: 'English', flag: '🇬🇧' },
  { code: 'ru', label: 'Русский', flag: '🇷🇺' },
];

export function LanguageSwitcher() {
  const { i18n } = useTranslation();
  
  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };
  
  return (
    <div className="language-switcher">
      {LANGUAGES.map(lang => (
        <button
          key={lang.code}
          onClick={() => changeLanguage(lang.code)}
          className={i18n.language === lang.code ? 'active' : ''}
          aria-label={`Switch to ${lang.label}`}
        >
          <span className="flag">{lang.flag}</span>
          <span className="label">{lang.label}</span>
        </button>
      ))}
    </div>
  );
}
```

#### URL Strategy

**Approach: Query Parameter** (simplest, chosen for Phase 1)

**Format:** `/names?lang=ru&year_from=1980...`

**Implementation:**
- Language preference stored in `localStorage`
- Optional `lang` query param overrides stored preference
- Sync language with URL on language change

**Alternative Approaches (Future Consideration):**
- **Path-based**: `/ru/names`, `/en/names` (better for SEO, requires route duplication)
- **Subdomain**: `ru.affirm-name.com` (enterprise, requires infrastructure)

#### Locale-Aware Formatting

**Numbers:**
```typescript
// Use Intl.NumberFormat for locale-aware formatting
const formatter = new Intl.NumberFormat(i18n.language);
<span>{formatter.format(125430)}</span>
// English: "125,430"
// Russian: "125 430"
```

**Dates:**
```typescript
const dateFormatter = new Intl.DateTimeFormat(i18n.language, {
  year: 'numeric',
  month: 'long',
});
<span>{dateFormatter.format(new Date())}</span>
// English: "November 2025"
// Russian: "ноябрь 2025 г."
```

---

### Backend Considerations

#### Error Messages

**Approach:** Return error codes only, frontend translates

**Current Format:**
```json
{
  "error": {
    "code": "invalid_params",
    "message": "Invalid year range: year_from must be <= year_to"
  }
}
```

**i18n-Ready Format (Optional Enhancement):**
```json
{
  "error": {
    "code": "invalid_year_range",
    "params": {
      "year_from": 2020,
      "year_to": 1980
    }
  }
}
```

Frontend translation:
```typescript
// en/errors.json: "invalid_year_range": "Invalid year range: {{year_from}} must be <= {{year_to}}"
// ru/errors.json: "invalid_year_range": "Неверный диапазон лет: {{year_from}} должен быть <= {{year_to}}"

const errorMessage = t(`errors.${error.code}`, error.params);
```

#### API Response Headers

**Optional: Add Content-Language header**
```go
w.Header().Set("Content-Language", "en")
```

**Note:** Not strictly necessary if UI language is client-side only.

---

### Data Layer Considerations

#### UTF-8 Encoding

**Database (PostgreSQL):**
- Encoding: `UTF8` (already configured)
- Collation: `en_US.UTF-8` or `C` (depending on sorting needs)
- Store names in original form (no normalization)

**SQL Example:**
```sql
CREATE DATABASE affirm_name
  ENCODING 'UTF8'
  LC_COLLATE 'en_US.UTF-8'
  LC_CTYPE 'en_US.UTF-8';
```

#### Name Sorting

**Challenge:** Names from different languages sort differently

**Approach 1: Server-Side (Current)**
```sql
-- Sort with database collation
SELECT name FROM names ORDER BY name COLLATE "en_US.UTF-8";
```

**Approach 2: Client-Side**
```typescript
// Use Intl.Collator for locale-aware sorting
const collator = new Intl.Collator(i18n.language, {
  sensitivity: 'base',
  numeric: true,
});

names.sort((a, b) => collator.compare(a.name, b.name));
```

**Recommendation:** Use server-side sorting initially. If users need locale-specific sorting (e.g., Russian users want Cyrillic sorted correctly), consider client-side sorting or multiple database collations.

#### Name Glob Matching

**Challenge:** Case-insensitive matching across scripts

**Current Implementation:**
```sql
WHERE name ILIKE 'alex%'
```

**Works for:** Latin, Cyrillic, most scripts
**Note:** PostgreSQL's `ILIKE` is Unicode-aware in UTF8 databases.

**Test Cases:**
- Latin: `alex*` matches "Alex", "Alexander", "ALEXANDRA"
- Cyrillic: `сан*` matches "Саня", "САНДРА", "Санька"
- Mixed: `*сан*` matches names containing "сан" in any case

---

### Content That Needs Translation

#### Essential (Phase 1)

**Navigation & Layout:**
- Logo/site name (if translated)
- Navigation menu items
- Footer links

**Main Page:**
- Headline and tagline
- Mission statement
- Feature descriptions
- Call-to-action buttons

**Filter Bar:**
- All filter labels
- Placeholders
- Help text
- Validation messages

**Names Table:**
- Column headers
- Pagination controls ("Previous", "Next", "Showing X-Y of Z")
- Empty state message
- Loading state text

**Name Detail Page:**
- Section headers
- Chart titles and labels
- Back button text

**Error Messages:**
- All error codes
- Retry button text
- Error boundary fallback

#### Optional (Phase 2)

**Meta/SEO:**
- Page titles (`<title>`)
- Meta descriptions
- Open Graph tags

**Country Names:**
- Translate country names in filter dropdown
- Requires country metadata with translations

**Data Source Descriptions:**
- Translate data source descriptions in tooltips

---

### Translation Workflow

#### 1. Development

**Initial Setup:**
1. Developer creates English strings in `en/*.json`
2. Keys use descriptive `snake_case` names
3. Commit to repository

**Adding New Strings:**
1. Add to English file first
2. Use `TODO` value in Russian file as placeholder
3. Mark for translation review

#### 2. Translation

**Tools:**
- Use translation management platform (e.g., Lokalise, Crowdin, POEditor)
- Or manual JSON editing for initial phase

**Process:**
1. Export English JSON files
2. Translators create Russian translations
3. Import translated JSON files
4. Review and test in-context

**Quality Checks:**
- Length verification (ensure UI doesn't break with longer text)
- Context verification (translators see screenshots)
- Technical term consistency (glossary of terms)

#### 3. Testing

**Key Test Cases:**
1. Switch languages and verify all UI updates
2. Check text wrapping and layout
3. Test with long/short translations
4. Verify placeholders and dynamic content
5. Test error messages in both languages

**Browser Testing:**
- Test in Chrome, Firefox, Safari
- Test on mobile devices
- Verify font rendering (especially Cyrillic)

---

### Future Enhancements

#### RTL Support (Arabic, Hebrew)

**CSS Approach:**
```css
/* Use logical properties */
.filter-bar {
  margin-inline-start: 1rem; /* Instead of margin-left */
  padding-inline: 1rem 2rem; /* Instead of padding-left/right */
}

/* RTL-specific overrides */
[dir="rtl"] .icon {
  transform: scaleX(-1); /* Flip icons */
}
```

**HTML:**
```html
<html lang="ar" dir="rtl">
```

#### Additional Languages

**To Add a New Language:**

1. **Update supported languages:**
```typescript
// src/i18n/config.ts
supportedLngs: ['en', 'ru', 'es', 'de'],
```

2. **Create translation files:**
```
public/locales/es/
  common.json
  filters.json
  ...
```

3. **Update LanguageSwitcher:**
```typescript
const LANGUAGES = [
  { code: 'en', label: 'English', flag: '🇬🇧' },
  { code: 'ru', label: 'Русский', flag: '🇷🇺' },
  { code: 'es', label: 'Español', flag: '🇪🇸' },
];
```

4. **Test thoroughly**

#### Country/Data Source Translations

**Approach:** Add translation keys to country metadata

**Database Schema Addition:**
```sql
-- Add to countries table (future)
ALTER TABLE countries ADD COLUMN translations JSONB;

-- Example value:
{
  "en": {
    "name": "United States",
    "data_source_description": "Social Security Administration baby names"
  },
  "ru": {
    "name": "Соединённые Штаты",
    "data_source_description": "Данные о детских именах от Администрации социального обеспечения"
  }
}
```

**API Response:**
```json
{
  "code": "US",
  "name_translations": {
    "en": "United States",
    "ru": "Соединённые Штаты"
  }
}
```

**Frontend Usage:**
```typescript
const countryName = country.name_translations[i18n.language] || country.name_translations.en;
```

---

### Environment Variables

**Add to [`05-cross-cutting-concerns.md`](05-cross-cutting-concerns.md#frontend-environment-variables):**

| Variable | Purpose | Example |
|----------|---------|---------|
| `VITE_DEFAULT_LANGUAGE` | Default UI language | `en` |
| `VITE_SUPPORTED_LANGUAGES` | Comma-separated language codes | `en,ru` |
| `VITE_TRANSLATION_BACKEND` | Translation file location | `/locales` |

---

### Testing Strategy

#### Unit Tests

**Test translation keys exist:**
```typescript
import en from '../public/locales/en/common.json';
import ru from '../public/locales/ru/common.json';

test('Russian translations have all English keys', () => {
  const enKeys = Object.keys(en);
  const ruKeys = Object.keys(ru);
  
  enKeys.forEach(key => {
    expect(ruKeys).toContain(key);
  });
});
```

#### Integration Tests

**Test language switching:**
```typescript
test('changes UI language when switching', async () => {
  render(<App />);
  
  const switcher = screen.getByLabelText(/switch to русский/i);
  await userEvent.click(switcher);
  
  await waitFor(() => {
    expect(screen.getByText(/применить фильтры/i)).toBeInTheDocument();
  });
});
```

#### Manual QA Checklist

- [ ] All UI text is translated
- [ ] No English text appears in Russian mode
- [ ] Numbers and dates are formatted correctly per locale
- [ ] Error messages appear in selected language
- [ ] Language preference persists across sessions
- [ ] Layout doesn't break with longer translations
- [ ] Tooltips and help text are translated
- [ ] Placeholder text is translated
- [ ] Accessibility labels are translated

---

---

## Security Considerations

### Backend

**Input Validation:**
- Validate all query parameters.
- Sanitize user inputs to prevent SQL injection.
- Use parameterized queries.

**Authentication & Authorization:**
- Not required for public read-only endpoints.
- Required for dataset upload endpoint (admin only).
- Use JWT or session-based auth.

**Rate Limiting:**
- Prevent abuse and DDoS attacks.
- Apply to all endpoints.
- More strict limits for write endpoints.

**CORS:**
- Configure allowed origins.
- Use whitelist approach.
- Validate `Origin` header.

### Frontend

**XSS Prevention:**
- React escapes content by default.
- Avoid `dangerouslySetInnerHTML`.
- Sanitize user inputs if needed.

**HTTPS:**
- Use HTTPS in production.
- Set `Secure` flag on cookies.
- Use HSTS header.

**Content Security Policy:**
- Set CSP headers to prevent XSS.
- Whitelist allowed sources.

**Dependency Management:**
- Keep dependencies up to date.
- Use `npm audit` to check for vulnerabilities.
- Monitor security advisories.

---

[← Previous: Integration Strategy](04-integration-strategy.md) | [Next: Development Workflow →](06-development-workflow.md)