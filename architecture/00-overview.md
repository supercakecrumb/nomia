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