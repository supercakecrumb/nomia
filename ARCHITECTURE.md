# Affirm Name – Architecture Plan

**Version:** 1.0  
**Last Updated:** 2025-11-15  
**Status:** Scaffolding / Carcass Definition

---

## About This Document

This architecture plan defines the **scaffolding** for Affirm Name, a data-driven web application helping trans, nonbinary, and other users explore given names across countries, languages, and decades.

The plan enables **parallel development** by establishing a stable contract that frontend and backend teams can implement against independently from day one.

---

## Table of Contents

### [00. Overview](architecture/00-overview.md)
Introduces the purpose, goals, key principles, and system context of the architecture.

**Key Topics:**
- Purpose and goals
- Architectural principles (Contract-First, Mock-First, Separation of Concerns)
- **Technology stack with pinned versions** (Go 1.25.4, PostgreSQL 18.1, React 19.2.0, React Router 7.9.6, and more)
- System context diagram

---

### [01. Shared Contract](architecture/01-shared-contract.md)
Defines the API surface, data semantics, and terminology—the core of the carcass.

**Key Topics:**
- Terminology: Gender Balance Axis, Popularity Metrics, Presence Period, Name Glob Filter
- API Endpoints: `/api/meta/years`, `/api/meta/countries`, `/api/names`, `/api/names/trend`
- Query parameters and response formats
- JSON fixtures for parallel development

---

### [02. Backend Carcass](architecture/02-backend-carcass.md)
Describes the backend skeleton: database schema, HTTP handlers, ingestion flow, and mock/real mode.

**Key Topics:**
- Database schema: `countries`, `name_datasets`, `names` tables
- HTTP handlers and routing structure
- Filter and popularity pipeline (6 stages)
- Dataset ingestion: upload endpoint, worker, parser abstraction
- Mock vs real mode configuration

---

### [03. Frontend Carcass](architecture/03-frontend-carcass.md)
Outlines the frontend structure: routing, state management, API client, and UI components.

**Key Topics:**
- Routing and page structure
- Global filter state model
- API client layer with fixture/real mode switching
- UI component skeleton: filter bar, names table, detail page, charts
- Mocking strategy (MSW vs direct imports)

---

### [04. Integration Strategy](architecture/04-integration-strategy.md)
Explains how frontend and backend teams work in parallel and integrate their work.

**Key Topics:**
- Backend workflow: fixture mode → ingestion → real queries → optimization
- Frontend workflow: fixture mode → interactions → visualizations → real API
- Contract validation and continuous integration
- Parallel development timeline (8 weeks)
- Risk mitigation strategies

---

### [05. Cross-Cutting Concerns](architecture/05-cross-cutting-concerns.md)
Addresses system-wide concerns: errors, performance, accessibility, configuration, i18n, security.

**Key Topics:**
- Error handling (standard format, codes, frontend/backend strategies)
- Performance (indexing, caching, debouncing, optimization)
- Accessibility (keyboard nav, screen readers, WCAG compliance)
- Configuration (environment variables, config files)
- Internationalization awareness
- Security considerations

---

### [06. Development Workflow](architecture/06-development-workflow.md)
Provides practical guidance on development, testing, deployment, and maintenance.

**Key Topics:**
- Initial setup (backend and frontend)
- Development phases (4 phases, 8 weeks)
- Testing strategy (unit, integration, E2E, accessibility)
- CI/CD pipeline
- Code review guidelines
- Deployment (backend and frontend platforms)
- Monitoring and troubleshooting

---

## Quick Navigation

### For Backend Developers
1. Start with [Shared Contract](architecture/01-shared-contract.md) to understand the API.
2. Read [Backend Carcass](architecture/02-backend-carcass.md) for implementation details.
3. Review [Development Workflow](architecture/06-development-workflow.md) for setup and testing.

### For Frontend Developers
1. Start with [Shared Contract](architecture/01-shared-contract.md) to understand the API.
2. Read [Frontend Carcass](architecture/03-frontend-carcass.md) for implementation details.
3. Review [Development Workflow](architecture/06-development-workflow.md) for setup and testing.

### For Project Managers
1. Read [Overview](architecture/00-overview.md) for high-level understanding.
2. Review [Integration Strategy](architecture/04-integration-strategy.md) for timeline and risks.
3. Check [Development Workflow](architecture/06-development-workflow.md) for phases and milestones.

---

## Summary

This architecture plan provides:

1. **Shared Contract**: Stable API specification that both teams implement against.
2. **Backend Carcass**: Database schema, handler structure, ingestion skeleton, mock/real mode.
3. **Frontend Carcass**: Routing, state management, API client, UI component skeleton.
4. **Integration Strategy**: How teams work in parallel and integrate their work.
5. **Cross-Cutting Concerns**: Error handling, performance, accessibility, configuration, i18n, security.
6. **Development Workflow**: Practical guidance on development, testing, and deployment.

**Key Benefit:** Frontend and backend teams can start development immediately using fixtures, then gradually replace mocks with real implementations while maintaining the same contract.

---

## Next Steps

1. **Review and Approve**: Both teams review and approve this architecture plan.
2. **Create Fixtures**: Create JSON fixture files in `/spec-examples/` based on contract.
3. **Set Up Projects**: Initialize backend (Go) and frontend (React/TypeScript) projects.
4. **Begin Phase 1**: Both teams start development in fixture mode.
5. **Regular Syncs**: Schedule weekly syncs to review progress and address integration issues.

---

**Questions or feedback?** Contact the architecture team or open an issue in the project repository.