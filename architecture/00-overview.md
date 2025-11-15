# Overview

## Purpose

This document defines the **architectural scaffolding** for Affirm Name—a data-driven web application helping trans, nonbinary, and other users explore given names across countries, languages, and decades to find names that fit their gender identity, cultural background, and legal realities.

## Goals

The architecture enables:

1. **Parallel Development**: Frontend and backend teams can work independently from day one.
2. **Contract-First Design**: A shared API contract and data semantics that both teams implement against.
3. **Incremental Implementation**: Start with fixtures/mocks, gradually replace with real implementations.
4. **Clear Boundaries**: Well-defined module responsibilities, data flows, and integration points.

## Key Principles

- **Contract as Source of Truth**: API endpoints, query parameters, response shapes, and semantics are defined upfront and remain stable.
- **Mock-First Development**: Both frontend and backend start by working with JSON fixtures that exemplify the contract.
- **Separation of Concerns**: Backend handles data ingestion, normalization, filtering, and aggregation; frontend handles presentation, interaction, and user experience.
- **Performance-Aware**: Design considers indexing strategies, pagination, and efficient filtering from the start.

## Technology Stack

This section defines the specific versions and tools used in the Affirm Name project as of 2025-11-16. These versions represent a **modern but conservative** approach: current stable releases that align with ecosystem directions while keeping the stack lean for future maintainability.

### Backend Stack (Go)

#### Core
- **Language / Toolchain**
  - Go: **1.25.4**
  
- **Database**
  - PostgreSQL: **18.1** (latest stable, released 2025-11)

#### HTTP & Routing
- **Standard Library**: `net/http` (no versioning, tied to Go 1.25.4)
- **Router**: `github.com/go-chi/chi/v5` @ **v5.2.3**

#### Database Access
- **Primary Driver**: `github.com/jackc/pgx/v5` @ **v5.7.6**
  - Use native pgx APIs for most operations
  - If you need `database/sql`, use `github.com/jackc/pgx/v5/stdlib` as the adapter
  - **Note**: Avoid `lib/pq` for new code; it's in maintenance mode

#### Migrations
- **Migration Tool**: `github.com/golang-migrate/migrate/v4` @ **v4.19.0**

#### SQL Codegen (Optional but Recommended)
- **Code Generator**: `github.com/sqlc-dev/sqlc` @ **v1.30.0**

#### Configuration
- **Config Library**: `github.com/spf13/viper` @ **v1.21.0**
- **Pattern**: Environment variables as the source of truth; Viper parses/merges into a typed Config struct validated on startup

#### Logging
- **Structured Logging**: `log/slog` over sugaring of `go.uber.org/zap` @ **v1.27.0**
- **Pattern**: Use the "production" config with context-aware loggers passed down per request

#### Observability (Direction)
- **Metrics & Tracing**: OpenTelemetry SDK for Go, exporting to your choice of backend (Prometheus / OTLP)
- **Note**: OTEL versions move quickly; use "latest stable in the OTEL Go 1.x line" rather than pinning

---

### Frontend Stack (React)

#### Runtime & Tooling
- **Node.js**: **24.11.1** (current LTS 24.x)
  - Meets React 19 and React Router 7 requirements
  - Exceeds Vite's Node minimum
  
- **Bundler / Dev Tool**: Vite **7.2.2**
- **Language**: TypeScript **5.9.3**

#### React Application Stack
- **Framework**: React **19.2.0** + react-dom **19.2.0**
  - React 19 is the current stable line
  
- **Routing**: react-router **7.9.6** + react-router-dom **7.9.6**
  - React Router v7 is the forward path, designed to bridge React 18→19
  - Can act as "library only" or "mini framework" for future flexibility
  
- **Data Fetching / Server State**: @tanstack/react-query **5.90.9**
  - Handles all server state management
  
- **Styling**: tailwindcss **4.1.17**
  - Primary styling engine
  - Optional: design system on top (Radix UI, MUI, or internal)
  
- **State Management**:
  - Global server state → TanStack Query
  - Light global UI state (filters, theme, layout) → React Context + custom hooks
  - For complex client-only state, consider Zustand or Jotai later (not required initially)
  
- **Charts**: recharts **3.4.1**

---

## System Context

```
┌─────────────────────────────────────────────────────────────────┐
│                    Public Statistical Agencies                   │
│              (SSA, ONS, SCB, etc. – CSV/XLSX files)             │
└────────────────────────┬────────────────────────────────────────┘
                         │ Manual Download
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Dataset Upload Module                        │
│              (Admin UI or CLI for file upload)                   │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Go Backend (Ingestion Layer)                   │
│         • Parse datasets                                         │
│         • Normalize to unified schema                            │
│         • Store in PostgreSQL                                    │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                         PostgreSQL                               │
│         • countries, name_datasets, names tables                 │
│         • Indexes for filtering and glob matching                │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Go Backend (REST API)                       │
│         • GET /api/meta/years                                    │
│         • GET /api/meta/countries                                │
│         • GET /api/names                                         │
│         • GET /api/names/trend                                   │
│         • POST /api/datasets/upload                              │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTP/JSON
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                   React/TypeScript Frontend                      │
│         • Main/Motivation Page (/)                               │
│         • Name Explorer Table (/names)                           │
│         • Name Detail Page (/name/:name)                         │
└─────────────────────────────────────────────────────────────────┘
```

---

[Next: Shared Contract →](01-shared-contract.md)