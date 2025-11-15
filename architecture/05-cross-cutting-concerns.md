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

| Code | HTTP Status | Meaning |
|------|-------------|---------|
| `invalid_params` | 400 | Invalid query parameters |
| `not_found` | 404 | Resource not found |
| `internal_error` | 500 | Server error |
| `database_error` | 500 | Database query failed |
| `parse_error` | 500 | Dataset parsing failed |

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

## Internationalization Awareness

### Current Scope

**UI Language:**
- UI is in English.
- Focus on English-speaking users initially.

**Data Support:**
- Data includes names from multiple languages and scripts (Latin, Cyrillic, Arabic, etc.).
- System is UTF-8 safe throughout.

### Future Considerations

**UI Localization:**
- Translate UI strings to multiple languages (Spanish, French, German, etc.).
- Use i18n library (e.g., `react-i18next`).
- Support language switching in UI.

**Name Normalization:**
- Handle diacritics and accents (e.g., "José" vs "Jose").
- Use Unicode normalization (NFC or NFD).
- Consider collation for sorting.

**Collation:**
- Use locale-aware sorting for name lists.
- PostgreSQL: Use collation (e.g., `COLLATE "en_US"`).
- Frontend: Use `Intl.Collator` for client-side sorting.

**RTL Support:**
- Support right-to-left languages (Arabic, Hebrew).
- Use CSS logical properties (`margin-inline-start` instead of `margin-left`).
- Test with RTL layouts.

### Architecture Implications

**UTF-8 Encoding:**
- Use UTF-8 encoding throughout (database, API, frontend).
- Store names in their original form (no normalization).
- Display names with proper character rendering.

**Database:**
- Use UTF-8 encoding for database (PostgreSQL: `UTF8`).
- Use appropriate collation for locale-aware sorting.
- Test with names from multiple scripts.

**API:**
- Ensure JSON responses are UTF-8 encoded.
- Handle URL encoding for query parameters with special characters.

**Frontend:**
- Set `<meta charset="utf-8">` in HTML.
- Use Unicode-aware string operations.
- Test with names from multiple languages and scripts.

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