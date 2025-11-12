
# REST API Specification

## Overview

This document defines the REST API for the baby name statistics platform. The API follows RESTful principles and uses JSON for request/response payloads.

**Base URL:** `https://api.affirm-name.com/v1`

**API Version:** v1

**Content Type:** `application/json`

**Authentication:** Bearer token (API Key or JWT)

---

## Table of Contents

1. [Authentication](#authentication)
2. [Common Patterns](#common-patterns)
3. [Error Handling](#error-handling)
4. [Rate Limiting](#rate-limiting)
5. [Endpoints](#endpoints)
   - [Countries](#countries)
   - [Datasets](#datasets)
   - [Names](#names)
   - [Trends](#trends)
   - [Jobs](#jobs)
   - [Health](#health)

---

## Authentication

### API Key Authentication

Include API key in the Authorization header:

```http
Authorization: Bearer <api_key>
```

**Example:**
```bash
curl -H "Authorization: Bearer ak_1234567890abcdef" \
     https://api.affirm-name.com/v1/countries
```

### Roles

- **admin**: Full access (upload, manage countries, view all data)
- **viewer**: Read-only access (query names, view trends)

---

## Common Patterns

### Pagination

All list endpoints support pagination using offset-based pagination.

**Query Parameters:**
- `limit`: Number of items per page (default: 100, max: 1000)
- `offset`: Number of items to skip (default: 0)

**Response Format:**
```json
{
  "data": [...],
  "meta": {
    "total": 1000,
    "limit": 100,
    "offset": 0,
    "has_more": true
  }
}
```

### Sorting

List endpoints support sorting using the `sort` parameter.

**Format:** `sort=field:direction`
- Direction: `asc` (ascending) or `desc` (descending)
- Default: Endpoint-specific

**Example:**
```
GET /v1/names?sort=count:desc
GET /v1/datasets?sort=uploaded_at:desc
```

### Filtering

List endpoints support filtering using query parameters.

**Example:**
```
GET /v1/names?country=US&year=2020&gender=F
GET /v1/datasets?status=completed&country=US
```

### Timestamps

All timestamps are in ISO 8601 format with UTC timezone.

**Example:** `2024-01-15T10:30:00Z`

---

## Error Handling

### Error Response Format

```json
{
  "error": {
    "code": "validation_error",
    "message": "Invalid request parameters",
    "details": [
      {
        "field": "year",
        "message": "must be between 1970 and 2030"
      }
    ]
  }
}
```

### HTTP Status Codes

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | Successful POST |
| 202 | Accepted | Async operation started |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Invalid request parameters |
| 401 | Unauthorized | Missing or invalid authentication |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource already exists |
| 413 | Payload Too Large | File too large |
| 422 | Unprocessable Entity | Validation failed |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |
| 503 | Service Unavailable | Service temporarily unavailable |

### Error Codes

| Code | Description |
|------|-------------|
| `validation_error` | Request validation failed |
| `authentication_error` | Authentication failed |
| `authorization_error` | Insufficient permissions |
| `not_found` | Resource not found |
| `conflict` | Resource conflict |
| `rate_limit_exceeded` | Too many requests |
| `internal_error` | Internal server error |
| `service_unavailable` | Service unavailable |

---

## Rate Limiting

**Limits:**
- Admin endpoints: 100 requests/minute
- Public endpoints: 1000 requests/minute per IP
- Burst: 2x rate limit

**Headers:**
```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640000000
```

**Rate Limit Exceeded Response:**
```json
{
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Rate limit exceeded",
    "retry_after": 60
  }
}
```

---

## Endpoints

### Countries

#### List Countries

Get a list of all countries.

**Endpoint:** `GET /v1/countries`

**Authentication:** Optional (public endpoint)

**Query Parameters:**
- `limit` (integer): Items per page (default: 100, max: 1000)
- `offset` (integer): Items to skip (default: 0)

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "code": "US",
      "name": "United States",
      "source_url": "https://www.ssa.gov/oact/babynames/",
      "attribution": "Social Security Administration",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "code": "GB",
      "name": "United Kingdom",
      "source_url": "https://www.ons.gov.uk/...",
      "attribution": "Office for National Statistics",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "total": 50,
    "limit": 100,
    "offset": 0,
    "has_more": false
  }
}
```

#### Get Country

Get a specific country by ID or code.

**Endpoint:** `GET /v1/countries/{id_or_code}`

**Authentication:** Optional (public endpoint)

**Path Parameters:**
- `id_or_code` (string): Country UUID or 2-letter code

**Response:** `200 OK`

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "US",
    "name": "United States",
    "source_url": "https://www.ssa.gov/oact/babynames/",
    "attribution": "Social Security Administration",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "stats": {
      "dataset_count": 51,
      "total_names": 1500000,
      "year_range": {
        "min": 1970,
        "max": 2020
      }
    }
  }
}
```

**Error Responses:**
- `404 Not Found`: Country not found

#### Create Country

Create a new country.

**Endpoint:** `POST /v1/countries`

**Authentication:** Required (admin only)

**Request Body:**

```json
{
  "code": "FR",
  "name": "France",
  "source_url": "https://www.insee.fr/",
  "attribution": "INSEE"
}
```

**Validation:**
- `code`: Required, 2 characters, uppercase, unique
- `name`: Required, 1-100 characters
- `source_url`: Optional, valid URL
- `attribution`: Optional, max 255 characters

**Response:** `201 Created`

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "code": "FR",
    "name": "France",
    "source_url": "https://www.insee.fr/",
    "attribution": "INSEE",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body
- `409 Conflict`: Country code already exists
- `422 Unprocessable Entity`: Validation failed

#### Update Country

Update an existing country.

**Endpoint:** `PATCH /v1/countries/{id}`

**Authentication:** Required (admin only)

**Path Parameters:**
- `id` (string): Country UUID

**Request Body:**

```json
{
  "name": "United States of America",
  "source_url": "https://www.ssa.gov/oact/babynames/",
  "attribution": "U.S. Social Security Administration"
}
```

**Response:** `200 OK`

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "US",
    "name": "United States of America",
    "source_url": "https://www.ssa.gov/oact/babynames/",
    "attribution": "U.S. Social Security Administration",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Error Responses:**
- `404 Not Found`: Country not found
- `422 Unprocessable Entity`: Validation failed

#### Delete Country

Delete a country (soft delete).

**Endpoint:** `DELETE /v1/countries/{id}`

**Authentication:** Required (admin only)

**Path Parameters:**
- `id` (string): Country UUID

**Response:** `204 No Content`

**Error Responses:**
- `404 Not Found`: Country not found
- `409 Conflict`: Country has associated datasets

---

### Datasets

#### List Datasets

Get a list of datasets.

**Endpoint:** `GET /v1/datasets`

**Authentication:** Required (admin only)

**Query Parameters:**
- `country` (string): Filter by country code or ID
- `status` (string): Filter by status (pending, processing, completed, failed)
- `limit` (integer): Items per page (default: 100, max: 1000)
- `offset` (integer): Items to skip (default: 0)
- `sort` (string): Sort field and direction (default: uploaded_at:desc)

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "650e8400-e29b-41d4-a716-446655440000",
      "country_id": "550e8400-e29b-41d4-a716-446655440000",
      "country_code": "US",
      "filename": "yob2020.txt",
      "file_size": 1048576,
      "status": "completed",
      "row_count": 32033,
      "uploaded_by": "admin@example.com",
      "uploaded_at": "2024-01-15T10:00:00Z",
      "processed_at": "2024-01-15T10:05:00Z",
      "processing_time_seconds": 300
    }
  ],
  "meta": {
    "total": 150,
    "limit": 100,
    "offset": 0,
    "has_more": true
  }
}
```

#### Get Dataset

Get a specific dataset by ID.

**Endpoint:** `GET /v1/datasets/{id}`

**Authentication:** Required (admin only)

**Path Parameters:**
- `id` (string): Dataset UUID

**Response:** `200 OK`

```json
{
  "data": {
    "id": "650e8400-e29b-41d4-a716-446655440000",
    "country_id": "550e8400-e29b-41d4-a716-446655440000",
    "country_code": "US",
    "country_name": "United States",
    "filename": "yob2020.txt",
    "file_path": "uploads/650e8400-e29b-41d4-a716-446655440000/original.csv",
    "file_size": 1048576,
    "status": "completed",
    "row_count": 32033,
    "error_message": null,
    "uploaded_by": "admin@example.com",
    "uploaded_at": "2024-01-15T10:00:00Z",
    "processed_at": "2024-01-15T10:05:00Z",
    "processing_time_seconds": 300,
    "job_id": "750e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Error Responses:**
- `404 Not Found`: Dataset not found

#### Upload Dataset

Upload a new dataset file.

**Endpoint:** `POST /v1/datasets/upload`

**Authentication:** Required (admin only)

**Content Type:** `multipart/form-data`

**Form Fields:**
- `file` (file): CSV file (required, max 100MB)
- `country_id` (string): Country UUID (required)
- `metadata` (JSON string): Optional metadata

**Request Example:**

```bash
curl -X POST https://api.affirm-name.com/v1/datasets/upload \
  -H "Authorization: Bearer <api_key>" \
  -F "file=@yob2020.txt" \
  -F "country_id=550e8400-e29b-41d4-a716-446655440000" \
  -F 'metadata={"year": 2020, "source": "SSA"}'
```

**Response:** `202 Accepted`

```json
{
  "data": {
    "dataset_id": "650e8400-e29b-41d4-a716-446655440000",
    "job_id": "750e8400-e29b-41d4-a716-446655440000",
    "status": "pending",
    "message": "Dataset uploaded successfully. Processing will begin shortly."
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid file or missing parameters
- `404 Not Found`: Country not found
- `413 Payload Too Large`: File exceeds 100MB
- `422 Unprocessable Entity`: Invalid file format

#### Delete Dataset

Delete a dataset (soft delete).

**Endpoint:** `DELETE /v1/datasets/{id}`

**Authentication:** Required (admin only)

**Path Parameters:**
- `id` (string): Dataset UUID

**Response:** `204 No Content`

**Error Responses:**
- `404 Not Found`: Dataset not found

#### Reprocess Dataset

Reprocess an existing dataset with updated parser.

**Endpoint:** `POST /v1/datasets/{id}/reprocess`

**Authentication:** Required (admin only)

**Path Parameters:**
- `id` (string): Dataset UUID

**Request Body:**

```json
{
  "reason": "parser_bug_fix"
}
```

**Response:** `202 Accepted`

```json
{
  "data": {
    "dataset_id": "650e8400-e29b-41d4-a716-446655440000",
    "job_id": "750e8400-e29b-41d4-a716-446655440001",
    "status": "reprocessing",
    "message": "Dataset reprocessing started"
  }
}
```

**Error Responses:**
- `404 Not Found`: Dataset not found
- `409 Conflict`: Dataset already being processed

---

### Names

#### List Names

Query names with filters.

**Endpoint:** `GET /v1/names`

**Authentication:** Optional (public endpoint, rate limited)

**Query Parameters:**
- `country` (string): Country code or ID (required)
- `year` (integer): Year (required, 1970-2030)
- `gender` (string): Gender filter (M or F, optional)
- `name` (string): Name prefix search (optional, min 2 chars)
- `min_count` (integer): Minimum count filter (optional)
- `limit` (integer): Items per page (default: 100, max: 1000)
- `offset` (integer): Items to skip (default: 0)
- `sort` (string): Sort field (count, name) and direction (default: count:desc)

**Response:** `200 OK`

```json
{
  "data": [
    {
      "name": "Emma",
      "gender": "F",
      "count": 15581,
      "year": 2020,
      "country_code": "US",
      "rank": 1
    },
    {
      "name": "Olivia",
      "gender": "F",
      "count": 17535,
      "year": 2020,
      "country_code": "US",
      "rank": 2
    }
  ],
  "meta": {
    "total": 32033,
    "limit": 100,
    "offset": 0,
    "has_more": true,
    "filters": {
      "country": "US",
      "year": 2020,
      "gender": "F"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request`: Missing required parameters
- `422 Unprocessable Entity`: Invalid parameter values

#### Search Names

Search names across all years and countries.

**Endpoint:** `GET /v1/names/search`

**Authentication:** Optional (public endpoint, rate limited)

**Query Parameters:**
- `q` (string): Search query (required, min 2 chars)
- `country` (string): Country code filter (optional)
- `gender` (string): Gender filter (M or F, optional)
- `limit` (integer): Items per page (default: 100, max: 1000)
- `offset` (integer): Items to skip (default: 0)

**Response:** `200 OK`

```json
{
  "data": [
    {
      "name": "Emma",
      "total_count": 1500000,
      "countries": ["US", "GB", "CA"],
      "year_range": {
        "min": 1970,
        "max": 2020
      },
      "gender_distribution": {
        "M": 0.5,
        "F": 99.5
      }
    }
  ],
  "meta": {
    "total": 1,
    "limit": 100,
    "offset": 0,
    "has_more": false
  }
}
```

#### Get Name Details

Get detailed information about a specific name.

**Endpoint:** `GET /v1/names/{name}`

**Authentication:** Optional (public endpoint, rate limited)

**Path Parameters:**
- `name` (string): Name to query

**Query Parameters:**
- `country` (string): Country code filter (optional)

**Response:** `200 OK`

```json
{
  "data": {
    "name": "Emma",
    "total_count": 1500000,
    "countries": [
      {
        "code": "US",
        "name": "United States",
        "count": 1200000,
        "year_range": {
          "min": 1970,
          "max": 2020
        }
      }
    ],
    "gender_distribution": {
      "M": {
        "count": 7500,
        "percentage": 0.5
      },
      "F": {
        "count": 1492500,
        "percentage": 99.5
      }
    },
    "popularity_trend": "increasing",
    "peak_year": 2020,
    "peak_count": 15581
  }
}
```

**Error Responses:**
- `404 Not Found`: Name not found

---

### Trends

#### Get Name Trend

Get trend data for a specific name over time.

**Endpoint:** `GET /v1/trends/{name}`

**Authentication:** Optional (public endpoint, rate limited)

**Path Parameters:**
- `name` (string): Name to analyze

**Query Parameters:**
- `country` (string): Country code (required)
- `start_year` (integer): Start year (optional, default: earliest available)
- `end_year` (integer): End year (optional, default: latest available)
- `gender` (string): Gender filter (M or F, optional)

**Response:** `200 OK`

```json
{
  "data": {
    "name": "Emma",
    "country_code": "US",
    "country_name": "United States",
    "year_range": {
      "start": 1970,
      "end": 2020
    },
    "trends": [
      {
        "year": 1970,
        "gender": "F",
        "count": 5000,
        "rank": 150,
        "percentage_of_total": 0.25
      },
      {
        "year": 1971,
        "gender": "F",
        "count": 5200,
        "rank": 145,
        "percentage_of_total": 0.26
      }
    ],
    "summary": {
      "total_count": 1200000,
      "average_count": 23529,
      "peak_year": 2020,
      "peak_count": 15581,
      "trend_direction": "increasing",
      "growth_rate": 2.5
    }
  }
}
```

**Error Responses:**
- `400 Bad Request`: Missing required parameters
- `404 Not Found`: No data found for name/country combination

#### Compare Names

Compare trends for multiple names.

**Endpoint:** `POST /v1/trends/compare`

**Authentication:** Optional (public endpoint, rate limited)

**Request Body:**

```json
{
  "names": ["Emma", "Olivia", "Ava"],
  "country": "US",
  "start_year": 2010,
  "end_year": 2020,
  "gender": "F"
}
```

**Validation:**
- `names`: Required, array of 2-5 names
- `country`: Required, country code
- `start_year`: Optional, integer
- `end_year`: Optional, integer
- `gender`: Optional, M or F

**Response:** `200 OK`

```json
{
  "data": {
    "country_code": "US",
    "year_range": {
      "start": 2010,
      "end": 2020
    },
    "comparisons": [
      {
        "name": "Emma",
        "trends": [
          {
            "year": 2010,
            "count": 17345,
            "rank": 1
          }
        ],
        "summary": {
          "total_count": 190000,
          "average_rank": 1.5
        }
      }
    ]
  }
}
```

#### Get Gender Probability

Get gender probability for a name.

**Endpoint:** `GET /v1/trends/{name}/gender`

**Authentication:** Optional (public endpoint, rate limited)

**Path Parameters:**
- `name` (string): Name to analyze

**Query Parameters:**
- `country` (string): Country code filter (optional, defaults to all)

**Response:** `200 OK`

```json
{
  "data": {
    "name": "Jordan",
    "male_count": 450000,
    "female_count": 150000,
    "total_count": 600000,
    "male_probability": 75.0,
    "female_probability": 25.0,
    "classification": "predominantly_male",
    "confidence": "high",
    "countries": [
      {
        "code": "US",
        "male_probability": 75.0,
        "female_probability": 25.0
      }
    ]
  }
}
```

**Classifications:**
- `strongly_male`: >90% male
- `predominantly_male`: 70-90% male
- `neutral`: 40-60% either gender
- `predominantly_female`: 70-90% female
- `strongly_female`: >90% female

**Confidence Levels:**
- `high`: >10,000 total occurrences
- `medium`: 1,000-10,000 occurrences
- `low`: <1,000 occurrences

---

### Jobs

#### Get Job Status

Get the status of a background job.

**Endpoint:** `GET /v1/jobs/{id}`

**Authentication:** Required (admin only)

**Path Parameters:**
- `id` (string): Job UUID

**Response:** `200 OK`

```json
{
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440000",
    "dataset_id": "650e8400-e29b-41d4-a716-446655440000",
    "type": "parse_dataset",
    "status": "completed",
    "attempts": 1,
    "max_attempts": 3,
    "created_at": "2024-01-15T10:00:00Z",
    "started_at": "2024-01-15T10:00:05Z",
    "completed_at": "2024-01-15T10:05:00Z",
    "processing_time_seconds": 295,
    "result": {
      "rows_processed": 32033,
      "rows_skipped": 0,
      "errors": []
    }
  }
}
```

**Status Values:**
- `queued`: Waiting to be processed
- `running`: Currently being processed
- `completed`: Successfully completed
- `failed`: Failed after max retries

**Error Responses:**
- `404 Not Found`: Job not found

#### List Jobs

Get a list of jobs.

**Endpoint:** `GET /v1/jobs`

**Authentication:** Required (admin only)

**Query Parameters:**
- `status` (string): Filter by status
- `type` (string): Filter by type
- `dataset_id` (string): Filter by dataset
- `limit` (integer): Items per page (default: 100, max: 1000)
- `offset` (integer): Items to skip (default: 0)
- `sort` (string): Sort field and direction (default: created_at:desc)

**Response:** `200 OK`

```json
{
  "data": [
    {
      "id": "750e8400-e29b-41d4-a716-446655440000",
      "dataset_id": "650e8400-e29b-41d4-a716-446655440000",
      "type": "parse_dataset",
      "status": "completed",
      "attempts": 1,
      "created_at": "2024-01-15T10:00:00Z",
      "completed_at": "2024-01-15T10:05:00Z"
    }
  ],
  "meta": {
    "total": 500,
    "limit": 100,
    "offset": 0,
    "has_more": true
  }
}
```

---

### Health

#### Health Check

Check if the API is healthy.

**Endpoint:** `GET /health`

**Authentication:** None

**Response:** `200 OK`

```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime_seconds": 86400,
  "timestamp": "2024-01-15T10:00:00Z",
  "checks": {
    "database": "connected",
    "storage": "accessible",
    "worker": "running"
  }
}
```

**Unhealthy Response:** `503 Service Unavailable`

```json
{
  "status": "unhealthy",
  "version": "1.0.0",
  "timestamp": "2024-01-15T10:00:00Z",
  "checks": {
    "database": "disconnected",
    "storage": "accessible",
    "worker": "stopped"
  }
}
```

#### Readiness Check

Check if the API is ready to serve traffic.

**Endpoint:** `GET /ready`

**Authentication:** None

**Response:** `200 OK` (ready) or `503 Service Unavailable` (not ready)

```json
{
  "ready": true
}
```

---

## OpenAPI Specification

The complete OpenAPI 3.0 specification is available at:

**Endpoint:** `GET /v1/openapi.json`

**Authentication:** None

This provides machine-readable API documentation compatible with tools like Swagger UI, Postman, and code generators.

---

## Webhooks (Future)

For future implementation, webhooks can notify external systems of events:

### Events

- `dataset.uploaded`: New dataset uploaded
- `dataset.completed`: Dataset processing completed
- `dataset.failed`: Dataset processing failed
- `job.completed`: Background job completed
- `job.failed`: Background job failed

### Webhook Payload

```json
{
  "event": "dataset.completed",
  "timestamp": "2024-01-15T10:05:00Z",
  "data": {
    "dataset_id": "650e8400-e29b-41d4-a716-446655440000",
    "country_code": "US",
    "row_count": 32033,
    "processing_time_seconds": 300
  }
}
```

---

## SDK Examples

### JavaScript/TypeScript

```typescript
import { AffirmNameClient } from '@affirm-name/sdk';

const client = new AffirmNameClient({
  apiKey: 'ak_1234567890abcdef',
  baseURL: 'https://api.affirm-name.com/v1'
});

// List names
const names = await client.names.list({
  country: 'US',
  year: 2020,
  gender: 'F',
  limit: 10
});

// Get trend
const trend = await client.trends.get('Emma', {
  country: 'US',
  startYear: 2010,
  endYear: 2020
});

// Upload dataset
const upload = await client.datasets.upload({
  file: fileBuffer,
  countryId: 'uuid',
  metadata: { year: 2020 }
});
```

### Python

```python
from affirm_name import Client

client = Client(api_key='ak_1234567890abcdef')

# List names
names = client.names.list(
    country='US',
    year=2020,
    gender='F',
    limit=10
)

# Get trend
trend = client.trends.get(
    'Emma',
    country='US',
    start_year=2010,
    end_year=2020
)

# Upload dataset
with open('yob2020.txt', 'rb') as f:
    upload = client.datasets.upload(
        file=f,
        country_id='uuid',
        metadata={'year': 2020}
    )
```

### Go

```go
package main

import (
    "github.com/affirm-name/go-sdk"
)

func main() {
    client := sdk.NewClient("ak_1234567890abcdef")
    
    // List names
    names, err := client.Names.List(context.Background(), &sdk.NameListParams{
        Country: "US",
        Year:    2020,
        Gender:  "F",
        Limit:   10,
    })
    
    // Get trend
    trend, err := client.Trends.Get(context.Background(), "Emma", &sdk.TrendParams{
        Country:   "US",
        StartYear: 2010,
        EndYear:   2020,
    })
    
    // Upload dataset
    file, _ := os.Open("yob2020.txt")
    upload, err := client.Datasets.Upload(context.Background(), &sdk.UploadParams{
        File:      file,
        CountryID: "uuid",
        Metadata:  map[string]interface{}{"year": 2020},
    })
}
```

---

## Versioning Strategy

### API Versioning

The API uses URL-based versioning:
- Current version: `/v1`
- Future versions: `/v2`, `/v3`, etc.

### Version Support Policy

- Current version (v1): Fully supported
- Previous version: Supported for 12 months after new version release
- Deprecated versions: 6-month sunset period with warnings

### Breaking Changes

Breaking changes require a new API version. Examples:
- Removing endpoints or fields
- Changing response structure
- Changing authentication method
- Changing required parameters

### Non-Breaking Changes

Non-breaking changes can be made to current version:
- Adding new endpoints
- Adding optional parameters
- Adding new fields to responses
- Adding new error codes

---

## Best Practices

### 1. Use Appropriate HTTP Methods

- `GET`: Retrieve resources (idempotent, cacheable)
- `POST`: Create resources or trigger actions
- `PUT`: Replace entire resource
- `PATCH`: Partial update
- `DELETE`: Remove resource

### 2. Handle Errors Gracefully

Always check status codes and handle errors:

```javascript
try {
  const response = await fetch('/v1/names?country=US&year=2020');
  if (!response.ok) {
    const error = await response.json();
    console.error('API Error:', error.error.message);
    return;
  }
  const data = await response.json();
  // Process data
} catch (error) {
  console.error('Network Error:', error);
}
```

### 3. Implement Retry Logic

For transient errors (500, 503), implement exponential backoff:

```javascript
async function fetchWithRetry(url, options, maxRetries = 3) {
  for (let i = 0; i < maxRetries; i++) {
    try {
      const response = await fetch(url, options);
      if (response.ok || response.status < 500) {
        return response;
      }
    } catch (error) {
      if (i === maxRetries - 1) throw error;
    }
    await new Promise(resolve => setTimeout(resolve, Math.pow(2, i) * 1000));
  }
}
```

### 4. Use Pagination for Large Datasets

Always paginate when fetching large result sets:

```javascript
async function fetchAllNames(params) {
  const allNames = [];
  let offset = 0;
  const limit = 1000;
  
  while (true) {
    const response = await client.names.list({
      ...params,
      limit,
      offset
    });
    
    allNames.push(...response.data);
    
    if (!response.meta.has_more) break;
    offset += limit;
  }
  
  return allNames;
}
```

### 5. Cache Responses

Cache stable data (historical statistics):

```javascript
const cache = new Map();

async function getCachedTrend(name, country) {
  const key = `${name}:${country}`;
  
  if (cache.has(key)) {
    return cache.get(key);
  }
  
  const trend = await client.trends.get(name, { country });
  cache.set(key, trend);
  
  // Cache for 24 hours
  setTimeout(() => cache.delete(key), 24 * 60 * 60 * 1000);
  
  return trend;
}
```

### 6. Respect Rate Limits

Monitor rate limit headers and back off when needed:

```javascript
function checkRateLimit(response) {
  const remaining = parseInt(response.headers.get('X-RateLimit-Remaining'));
  const reset = parseInt(response.headers.get('X-RateLimit-Reset'));
  
  if (remaining < 10) {
    const waitTime = reset - Math.floor(Date.now() / 1000);
    console.warn(`Rate limit low. Resets in ${waitTime}s`);
  }
}
```

---

## Performance Tips

### 1. Use Specific Filters

More specific queries are faster:

```javascript
// Slow: Fetch all names then filter
const all = await client.names.list({ country: 'US', year: 2020 });
const filtered = all.data.filter(n => n.gender === 'F');

// Fast: Filter at database level
const filtered = await client.names.list({
  country: 'US',
  year: 2020,
  gender: 'F'
});
```

### 2. Request Only Needed Fields

Use field selection when available (future feature):

```javascript
// Request only needed fields
const names = await client.names.list({
  country: 'US',
  year: 2020,
  fields: ['name', 'count']  // Future feature
});
```

### 3. Batch Requests

Use batch endpoints when available:

```javascript
// Instead of multiple requests
const emma = await client.trends.get('Emma', { country: 'US' });
const olivia = await client.trends.get('Olivia', { country: 'US' });

// Use compare endpoint
const comparison = await client.trends.compare({
  names: ['Emma', 'Olivia'],
  country: 'US'
});
```

### 4. Use Compression

Enable gzip compression for responses:

```javascript
fetch('/v1/names', {
  headers: {
    'Accept-Encoding': 'gzip'
  }
});
```

---

## Security Considerations

### 1. Protect API Keys

Never expose API keys in client-side code:

```javascript
// ❌ Bad: API key in frontend
const client = new AffirmNameClient({
  apiKey: 'ak_1234567890abcdef'  // Exposed to users!
});

// ✅ Good: API key in backend
// Frontend calls your backend, which calls Affirm Name API
```

### 2. Use HTTPS

Always use HTTPS in production:

```javascript
const client = new AffirmNameClient({
  baseURL: 'https://api.affirm-name.com/v1'  // Not http://
});
```

### 3. Validate Input

Validate user input before sending to API:

```javascript
function validateYear(year) {
  const y = parseInt(year);
  if (isNaN(y) || y < 1970 || y > 2030) {
    throw new Error('Invalid year');
  }
  return y;
}
```

### 4. Handle Sensitive Data

Don't log sensitive information:

```javascript
// ❌ Bad: Logs API key
console.log('Request:', { apiKey: key, ...params });

// ✅ Good: Redacts sensitive data
console.log('Request:', { apiKey: '[REDACTED]', ...params });
```

---

## Testing

### Example Test Cases

```javascript
describe('Names API', () => {
  it('should list names for valid country and year', async () => {
    const response = await client.names.list({
      country: 'US',
      year: 2020,
      gender: 'F',
      limit: 10
    });
    
    expect(response.data).toHaveLength(10);
    expect(response.data[0]).toHaveProperty('name');
    expect(response.data[0]).toHaveProperty('count');
    expect(response.meta.total).toBeGreaterThan(0);
  });
  
  it('should return 400 for invalid year', async () => {
    await expect(
      client.names.list({ country: 'US', year: 1800 })
    ).rejects.toThrow('Invalid year');
  });
  
  it('should handle pagination correctly', async () => {
    const page1 = await client.names.list({
      country: 'US',
      year: 2020,
      limit: 100,
      offset: 0
    });
    
    const page2 = await client.names.list({
      country: 'US',
      year: 2020,
      limit: 100,
      offset: 100
    });
    
    expect(page1.data[0].name).not.toBe(page2.data[0].name);
  });
});
```

---

## Changelog

### v1.0.0 (2024-01-15)

**Initial Release**
- Countries CRUD endpoints
- Dataset upload and management
- Name query and search
- Trend analysis
- Gender probability
- Job status tracking
- Health checks

---

## Support

### Documentation

- Architecture: `/docs/architecture.md`
- Database Schema: `/docs/database-schema.md`
- API Specification: `/docs/api-specification.md` (this document)

### Contact

- Technical Support: support@affirm-name.com
- API Issues: api-issues@affirm-name.com
- Security: security@affirm-name.com

### Rate Limit Increase

For higher rate limits, contact support with:
- Use case description
- Expected request volume
- Business justification

---

## Appendix: Complete Request/Response Examples

### Upload and Query Flow

```bash
# 1. Create country (admin)
curl -X POST https://api.affirm-name.com/v1/countries \
  -H "Authorization: Bearer <admin_key>" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "US",
    "name": "United States",
    "source_url": "https://www.ssa.gov/oact/babynames/",
    "attribution": "Social Security Administration"
  }'

# Response: 201 Created
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "code": "US",
    "name": "United States",
    ...
  }
}

# 2. Upload dataset (admin)
curl -X POST https://api.affirm-name.com/v1/datasets/upload \
  -H "Authorization: Bearer <admin_key>" \
  -F "file=@yob2020.txt" \
  -F "country_id=550e8400-e29b-41d4-a716-446655440000"

# Response: 202 Accepted
{
  "data": {
    "dataset_id": "650e8400-e29b-41d4-a716-446655440000",
    "job_id": "750e8400-e29b-41d4-a716-446655440000",
    "status": "pending"
  }
}

# 3. Check job status (admin)
curl https://api.affirm-name.com/v1/jobs/750e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <admin_key>"

# Response: 200 OK
{
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440000",
    "status": "completed",
    "result": {
      "rows_processed": 32033
    }
  }
}

# 4. Query names (public)
curl "https://api.affirm-name.com/v1/names?country=US&year=2020&gender=F&limit=10"

# Response: 200 OK
{
  "data": [
    {
      "name": "Emma",
      "gender": "F",
      "count": 15581,
      "year": 2020,
      "country_code": "US",
      "rank": 1
    },
    ...
  ],
  "meta": {
    "total": 18252,
    "limit": 10,
    "offset": 0,
    "has_more": true
  }
}

# 5. Get trend (public)
curl "https://api.affirm-name.com/v1/trends/Emma?country=US&start_year=2010&end_year=2020"

# Response: 200 OK
{
  "data": {
    "name": "Emma",
    "country_code": "US",
    "trends": [
      {
        "year": 2010,
        "gender": "F",
        "count": 17345,
        "rank": 1
      },
      ...
    ],
    "summary": {
      "total_count": 190000,
      "peak_year": 2020,
      "trend_direction": "increasing"
    }
  }
}
```

---

## Conclusion

This API specification provides a complete reference for integrating with the Affirm Name platform. The API is designed to be:

✅ **RESTful**: Follows REST principles and conventions
✅ **Consistent**: Predictable patterns across all endpoints
✅ **Well-documented**: Clear examples and error messages
✅ **Performant**: Optimized for common use cases
✅ **Secure**: Authentication and rate limiting
✅ **Extensible**: Versioned for future enhancements

For implementation details, refer to the architecture and database schema documents.