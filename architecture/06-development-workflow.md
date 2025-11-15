# Development Workflow

This chapter outlines the practical development workflow, testing strategies, and deployment considerations.

## Initial Setup

### Backend Setup

1. **Project Structure:**
```
affirm-name-backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── models/
│   ├── parsers/
│   └── db/
├── migrations/
├── spec-examples/
├── config/
└── go.mod
```

2. **Database Setup:**
   - Install PostgreSQL 14+.
   - Create database: `affirm_name`.
   - Run migrations to create schema.
   - Enable `pg_trgm` extension.

3. **Environment Configuration:**
   - Copy `.env.example` to `.env`.
   - Set `FIXTURE_MODE=true` initially.
   - Set `DATABASE_URL`.

4. **Run Server:**
   - `go run cmd/server/main.go`.
   - Server starts in fixture mode.
   - Test endpoints with curl or Postman.

### Frontend Setup

1. **Project Structure:**
```
affirm-name-frontend/
├── src/
│   ├── api/
│   ├── components/
│   ├── pages/
│   ├── hooks/
│   ├── contexts/
│   └── types/
├── spec-examples/
├── public/
└── package.json
```

2. **Install Dependencies:**
   - `npm install` or `yarn install`.
   - Install React, TypeScript, React Router, TanStack Query, Tailwind.

3. **Environment Configuration:**
   - Copy `.env.example` to `.env.development`.
   - Set `VITE_API_MODE=mock` initially.

4. **Run Development Server:**
   - `npm run dev` or `yarn dev`.
   - App starts at `http://localhost:5173`.
   - UI loads with fixture data.

---

## Development Phases

### Phase 1: Foundation (Weeks 1–2)

**Backend Tasks:**
- [ ] Create database schema and migrations.
- [ ] Implement fixture handlers for all endpoints.
- [ ] Set up error handling and logging.
- [ ] Deploy locally with `FIXTURE_MODE=true`.

**Frontend Tasks:**
- [ ] Set up React + TypeScript project.
- [ ] Create routing structure (`/`, `/names`, `/name/:name`).
- [ ] Implement API client with fixture imports.
- [ ] Build basic page layouts and navigation.

**Validation:**
- Backend serves fixture responses correctly.
- Frontend navigates between pages.
- Fixtures match contract specification.

---

### Phase 2: Core Features (Weeks 3–4)

**Backend Tasks:**
- [ ] Implement dataset upload endpoint.
- [ ] Create parser interface and stubs.
- [ ] Implement background worker skeleton.
- [ ] Begin real query implementation.

**Frontend Tasks:**
- [ ] Build filter bar components (year, countries, gender, popularity, glob).
- [ ] Implement filter state management (Context + hooks).
- [ ] Sync filters with URL query parameters.
- [ ] Build names table with skeleton rows.

**Validation:**
- Upload endpoint accepts files and returns dataset ID.
- Filters update URL and trigger API calls.
- Table displays fixture data correctly.

---

### Phase 3: Visualization & Data (Weeks 5–6)

**Backend Tasks:**
- [ ] Complete real query implementation for `/api/names`.
- [ ] Implement popularity computation pipeline.
- [ ] Switch handlers to real mode (`FIXTURE_MODE=false`).
- [ ] Optimize queries with indexes.

**Frontend Tasks:**
- [ ] Build name detail page layout.
- [ ] Implement charts (popularity over time, gender distribution, country breakdown).
- [ ] Build gender balance visualization for table.
- [ ] Implement presence period formatting.

**Validation:**
- Real queries return correct, performant results.
- Charts display data accurately.
- Gender balance visualization is clear and intuitive.

---

### Phase 4: Integration & Polish (Weeks 7–8)

**Backend Tasks:**
- [ ] Add indexes for performance optimization.
- [ ] Implement comprehensive error handling.
- [ ] Add rate limiting.
- [ ] Write API documentation.

**Frontend Tasks:**
- [ ] Switch to real API mode (`VITE_API_MODE=real`).
- [ ] Implement accessibility improvements (keyboard nav, ARIA labels).
- [ ] Add loading states and skeletons.
- [ ] Polish UI (animations, transitions, error states).

**Validation:**
- Frontend works correctly against real backend.
- All accessibility requirements met (WCAG AA).
- Performance targets achieved (<500ms API responses).

---

## Testing Strategy

### Backend Testing

**Unit Tests:**
- Test parsers with sample data files.
- Test query builders and filter logic.
- Test popularity computation.
- Target: >80% code coverage.

**Integration Tests:**
- Test handlers with real database.
- Test dataset ingestion flow.
- Test error handling.

**Contract Tests:**
- Validate real responses match fixtures.
- Use JSON schema or OpenAPI spec.
- Run in CI pipeline.

**Performance Tests:**
- Load test with realistic data volumes.
- Profile slow queries with `EXPLAIN ANALYZE`.
- Target: <500ms for `/api/names`, <200ms for `/api/meta/*`.

**Tools:**
- `testing` package (Go standard library).
- `testify` for assertions.
- `sqlmock` for database mocking.
- `k6` or `vegeta` for load testing.

### Frontend Testing

**Unit Tests:**
- Test utility functions (date formatting, presence period logic).
- Test hooks (filter state, debouncing).
- Target: >70% code coverage.

**Component Tests:**
- Test UI components in isolation.
- Use React Testing Library.
- Test interactions (clicks, typing, selections).

**Integration Tests:**
- Test page flows with mocked API.
- Use Mock Service Worker (MSW).
- Test filter interactions and API calls.

**E2E Tests:**
- Test full user flows with real backend.
- Use Playwright or Cypress.
- Test: search for name, view details, apply filters.

**Accessibility Tests:**
- Use `axe-core` or `jest-axe`.
- Test keyboard navigation manually.
- Test with screen reader (NVDA, JAWS, VoiceOver).

**Tools:**
- `vitest` for unit and component tests.
- `@testing-library/react` for component tests.
- `MSW` for API mocking.
- `Playwright` for E2E tests.
- `axe-core` for accessibility tests.

---

## CI/CD Pipeline

### Backend CI

**On Pull Request:**
1. Run linter (`golangci-lint`).
2. Run unit tests.
3. Run integration tests.
4. Run contract tests (validate responses match fixtures).
5. Build Docker image.

**On Merge to Main:**
1. Run all tests.
2. Build and push Docker image to registry.
3. Deploy to staging environment.
4. Run smoke tests.

### Frontend CI

**On Pull Request:**
1. Run linter (`eslint`).
2. Run type checker (`tsc`).
3. Run unit and component tests.
4. Run accessibility tests.
5. Build production bundle.

**On Merge to Main:**
1. Run all tests.
2. Build production bundle.
3. Deploy to staging environment (Vercel/Netlify preview).
4. Run E2E tests against staging.

### Deployment Triggers

**Staging:**
- Automatic on merge to `main`.
- Always uses `FIXTURE_MODE=false` (real data).

**Production:**
- Manual approval required.
- Tag release with version number.
- Deploy backend first, then frontend.
- Run smoke tests after deployment.

---

## Code Review Guidelines

### Backend Reviews

**Focus Areas:**
- SQL queries are optimized and use proper indexes.
- Error handling is comprehensive.
- Responses match contract specification.
- No SQL injection vulnerabilities.
- Proper logging for debugging.

**Checklist:**
- [ ] Tests pass.
- [ ] Responses match fixtures.
- [ ] Queries are performant.
- [ ] Errors are handled gracefully.
- [ ] Code is documented.

### Frontend Reviews

**Focus Areas:**
- Components are accessible (keyboard nav, ARIA labels).
- State management is clear and predictable.
- API calls are debounced appropriately.
- Loading and error states are handled.
- UI is responsive and polished.

**Checklist:**
- [ ] Tests pass.
- [ ] Accessibility requirements met.
- [ ] No console errors or warnings.
- [ ] Types are correct.
- [ ] Code is documented.

---

## Deployment

### Backend Deployment

**Platform Options:**
- **Fly.io**: Simple, PostgreSQL included, low cost.
- **Railway**: Similar to Fly.io, easy setup.
- **AWS**: ECS + RDS, more complex but scalable.

**Steps:**
1. Build Docker image.
2. Push to container registry.
3. Deploy to platform.
4. Run migrations.
5. Set environment variables (`FIXTURE_MODE=false`, `DATABASE_URL`).
6. Verify health endpoint (`/health`).

**Database:**
- Use managed PostgreSQL (Fly.io Postgres, AWS RDS, etc.).
- Enable automatic backups.
- Set up monitoring and alerts.

### Frontend Deployment

**Platform Options:**
- **Vercel**: Optimized for Vite/React, automatic deployments.
- **Netlify**: Similar to Vercel, easy setup.
- **Cloudflare Pages**: Fast, global CDN.

**Steps:**
1. Build production bundle (`npm run build`).
2. Upload to platform (or connect GitHub repo).
3. Set environment variables (`VITE_API_MODE=real`, `VITE_API_BASE_URL`).
4. Configure redirects for client-side routing.
5. Verify deployment.

**CDN:**
- Use platform's built-in CDN.
- Configure caching headers.
- Enable gzip/brotli compression.

---

## Monitoring & Maintenance

### Monitoring

**Backend:**
- Monitor API response times (target: p95 <500ms).
- Monitor error rates (target: <1%).
- Monitor database connections and query performance.
- Set up alerts for critical errors.

**Frontend:**
- Monitor Core Web Vitals (LCP, FID, CLS).
- Monitor JavaScript errors (Sentry).
- Monitor API call success rates.
- Monitor bundle size.

**Tools:**
- **Backend**: Prometheus + Grafana, Datadog, New Relic.
- **Frontend**: Google Analytics, Sentry, Vercel Analytics.

### Maintenance Tasks

**Regular:**
- Update dependencies (monthly).
- Review and address security vulnerabilities.
- Monitor and optimize database queries.
- Review error logs and fix issues.

**Periodic:**
- Backup database (automated daily).
- Review and archive old datasets.
- Performance testing and optimization.
- Update documentation.

---

## Troubleshooting

### Common Backend Issues

**Slow Queries:**
- Check indexes are being used (`EXPLAIN ANALYZE`).
- Add missing indexes.
- Optimize query logic.
- Consider precomputed aggregates.

**High Error Rate:**
- Check logs for error patterns.
- Fix validation issues.
- Improve error handling.
- Add retries where appropriate.

**Database Connection Issues:**
- Check connection pool size.
- Increase pool size if needed.
- Check for connection leaks.
- Monitor database load.

### Common Frontend Issues

**Slow Page Load:**
- Check bundle size (use `webpack-bundle-analyzer`).
- Code split large dependencies.
- Lazy load routes and components.
- Optimize images and assets.

**API Call Failures:**
- Check network tab in browser DevTools.
- Verify API base URL is correct.
- Check CORS configuration.
- Add retry logic.

**Layout Shifts:**
- Use skeleton loaders.
- Set explicit dimensions for images.
- Avoid injecting content above fold.
- Test with Lighthouse.

---

## Next Steps

After completing the development workflow:

1. **Create fixture files** in `/spec-examples/`.
2. **Seed database** with sample data for development.
3. **Set up CI/CD** pipelines for both frontend and backend.
4. **Begin Phase 1** development.
5. **Schedule regular syncs** between frontend and backend teams.

---

[← Previous: Cross-Cutting Concerns](05-cross-cutting-concerns.md) | [Back to Overview](00-overview.md)