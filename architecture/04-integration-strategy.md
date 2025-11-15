# Integration Strategy

The carcass enables **parallel development** by defining a stable contract that both teams implement against.

## Backend Team Workflow

### Phase 1: Fixture Mode

- Set up database schema (migrations).
- Implement HTTP handlers that return fixture JSON.
- Validate handler responses against fixtures.
- Deploy backend in fixture mode for frontend to test against.

### Phase 2: Ingestion Skeleton

- Implement dataset upload endpoint.
- Implement parser stubs.
- Implement background worker skeleton.
- Test with sample datasets.

### Phase 3: Real Queries

- Implement database queries for each endpoint.
- Compute filters, aggregations, and popularity metrics.
- Switch handlers from fixture mode to real mode.
- Validate responses still match contract.

### Phase 4: Optimization

- Add indexes for performance.
- Implement precomputed aggregates if needed.
- Optimize query plans.

---

## Frontend Team Workflow

### Phase 1: Fixture Mode

- Set up routing and page structure.
- Implement filter state management.
- Build API client that returns fixtures.
- Develop UI components with static data.

### Phase 2: Interactions

- Implement filter interactions (sliders, dropdowns, etc.).
- Sync filter state with URL.
- Implement debouncing and API call triggers.
- Test with fixture data.

### Phase 3: Charts & Visualizations

- Implement charts on detail page.
- Implement gender balance visualizations.
- Test with fixture data.

### Phase 4: Real API Integration

- Switch API client from fixture mode to real mode.
- Update `VITE_API_MODE` environment variable.
- Test against real backend.
- Fix any contract mismatches.

---

## Integration Points

### Contract Validation

**Backend:**
- Use JSON schema or OpenAPI spec to validate API responses.
- Validate real responses against fixtures.
- Run validation in CI pipeline.

**Frontend:**
- Validate API responses against TypeScript types.
- Use type guards for runtime validation.
- Run type checking in CI pipeline.

**Shared:**
- Use fixtures as shared reference.
- Any changes to contract must update fixtures first.
- Both teams review and approve fixture changes.

### Continuous Integration

**Contract Tests:**
- Backend runs tests that validate responses match contract.
- Frontend runs tests that validate types match contract.
- Fail build if contract is violated.

**Integration Tests:**
- E2E tests that exercise full user flows.
- Run against real backend (not fixtures).
- Validate that frontend and backend work together correctly.

### Communication Protocol

**Contract Changes:**
1. Propose change to contract (API endpoint, query params, response shape).
2. Update fixtures to reflect proposed change.
3. Both teams review and approve change.
4. Backend team implements change first.
5. Frontend team updates to use new contract.

**Versioning:**
- Use semantic versioning for API (v1, v2, etc.).
- Maintain backward compatibility when possible.
- Document breaking changes clearly.

---

## Parallel Development Timeline

### Week 1–2: Foundation

**Backend:**
- Database schema and migrations
- HTTP handlers returning fixtures
- Basic routing and error handling

**Frontend:**
- Project setup and routing
- API client with fixture imports
- Basic page structure and layout

**Collaboration:**
- Review and approve fixtures
- Establish CI/CD pipeline
- Set up contract validation

### Week 3–4: Core Features

**Backend:**
- Ingestion endpoint and worker skeleton
- Parser stubs for SSA and ONS
- Real query implementation (begin)

**Frontend:**
- Filter bar components
- Filter state management
- URL synchronization
- Names table with skeleton

**Collaboration:**
- Test frontend against fixture-mode backend
- Review API responses
- Adjust contract if needed

### Week 5–6: Visualization & Data

**Backend:**
- Complete real query implementation
- Popularity computation
- Switch to real mode
- Optimize queries

**Frontend:**
- Name detail page
- Charts and visualizations
- Gender balance displays
- Polish interactions

**Collaboration:**
- Test frontend against real-mode backend
- Fix integration issues
- Performance testing

### Week 7–8: Integration & Polish

**Backend:**
- Performance optimization
- Additional indexes
- Error handling improvements
- Documentation

**Frontend:**
- Switch to real API mode
- Accessibility improvements
- Error handling
- Loading states and polish

**Collaboration:**
- E2E testing
- User acceptance testing
- Bug fixes
- Deployment preparation

---

## Risk Mitigation

### Contract Drift

**Risk:** Frontend and backend interpret contract differently.

**Mitigation:**
- Use fixtures as single source of truth.
- Validate responses against fixtures in tests.
- Use TypeScript types generated from contract.
- Regular sync meetings to review contract.

### Integration Delays

**Risk:** Backend not ready when frontend needs it.

**Mitigation:**
- Frontend develops against fixtures first.
- Backend deploys in fixture mode early.
- Clear milestones for switching to real mode.

### Performance Issues

**Risk:** Real queries are too slow.

**Mitigation:**
- Design indexes upfront.
- Set performance budgets (e.g., <500ms response).
- Test with realistic data volumes early.
- Profile and optimize queries.

### Scope Creep

**Risk:** Contract changes frequently during development.

**Mitigation:**
- Lock contract early (within first 1–2 weeks).
- Require both teams to approve changes.
- Document breaking changes clearly.
- Use versioning for major changes.

---

[← Previous: Frontend Carcass](03-frontend-carcass.md) | [Next: Cross-Cutting Concerns →](05-cross-cutting-concerns.md)