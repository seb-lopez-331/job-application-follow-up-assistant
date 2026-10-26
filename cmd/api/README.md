
# API Contract

## Job Applications
These define the API endpoints for registered job applications.
| Method | Endpoint | Purpose |
| -------| -------- | ------- |
| `POST` | `/api/v1/applications` | Register a job application |
| `GET` | `/api/v1/applications` | List job applications with filters |
| `GET` | `/api/v1/applications/{id}` | Get one job application |
| `PATCH` | `/api/v1/applications/{id}` | Update selected fields |
| `DELETE` | `/api/v1/applications/{id}` | Unregister a job application |

List filters:
```
?stage=applied&search=platform&limit=20&cursor=...
```

### Register a job application
Request:
```
POST /api/v1/applications
{
  "company_name": "Acme",
  "role_title": "Platform Engineer",
  "job_url": "https://example.com/jobs/123",
  "location": "Remote, US",
  "work_arrangement": "remote",
  "status": "applied",
  "applied_on": "2026-10-26",
  "notes": "Interested in the infrastructure team."
}
```

Response:
```
201 Created
{
  "data": {
    "id": "8e34296b-ef82-4d6a-9b1f-30c49e3f0c11",
    "company_name": "Acme",
    "role_title": "Platform Engineer",
    "status": "applied",
    "applied_on": "2026-10-26",
    "next_follow_up_on": "2026-11-04",
    "created_at": "2026-10-26T16:00:00Z",
    "updated_at": "2026-10-26T16:00:00Z"
  }
}
```

### List job applications with filters
Request:
```
GET /api/v1/applications
```

Response:
```
200 OK
{
    "data": [
        {
            "id": "8e34296b-ef82-4d6a-9b1f-30c49e3f0c11",
            "company_name": "Acme",
            "role_title": "Platform Engineer",
            "status": "applied",
            "applied_on": "2026-10-26",
            "next_follow_up_on": "2026-11-04",
            "created_at": "2026-10-26T16:00:00Z",
            "updated_at": "2026-10-26T16:00:00Z"
        },
        {
            "id": "5213fcaf-15ea-4892-8419-d0d7077ebe16",
            "company_name": "CloudCo",
            "role_title": "Software Engineer",
            "status": "applied",
            "applied_on": "2026-10-22",
            "next_follow_up_on": "2026-10-29",
            "created_at": "2026-10-22T12:59:00Z",
            "updated_at": "2026-10-22T12:59:00Z"
        },
        {
            "id": "d02d9a11-f57b-46ab-9cf6-52d4c4992630",
            "company_name": "DataWorks",
            "role_title": "Infrastructure SWE",
            "status": "interviewing",
            "applied_on": "2026-10-24",
            "next_follow_up_on": "2026-10-31",
            "created_at": "2026-10-24T13:23:00Z",
            "updated_at": "2026-10-24T13:23:00Z"
        }
    ]
}
```

Request with filter:
```
GET /api/v1/applications?status=interviewing
```

Response with filter:
```
200 OK
{
    "data": [
        {
            "id": "d02d9a11-f57b-46ab-9cf6-52d4c4992630",
            "company_name": "DataWorks",
            "role_title": "Infrastructure SWE",
            "status": "interviewing",
            "applied_on": "2026-10-24",
            "next_follow_up_on": "2026-10-31",
            "created_at": "2026-10-24T13:23:00Z",
            "updated_at": "2026-10-24T13:23:00Z"
        }
    ]
}
```

### Get job application by ID
Request:
```
GET /api/v1/applications/5213fcaf-15ea-4892-8419-d0d7077ebe16
```

Response:
```
200 OK
{
  "data": {
    "id": "5213fcaf-15ea-4892-8419-d0d7077ebe16",
    "company_name": "CloudCo",
    "role_title": "Software Engineer",
    "status": "applied",
    "applied_on": "2026-10-22",
    "next_follow_up_on": "2026-10-29",
    "created_at": "2026-10-22T12:59:00Z",
    "updated_at": "2026-10-22T12:59:00Z"
    }
}

```

### Update selected fields
Request:
```
PATCH /api/v1/applications/8e34296b-ef82-4d6a-9b1f-30c49e3f0c11
{
  "status": "rejected"
}
```

Response:
```
200 OK
{
  "data": {
    "id": "8e34296b-ef82-4d6a-9b1f-30c49e3f0c11",
    "company_name": "Acme",
    "role_title": "Platform Engineer",
    "status": "rejected",
    "applied_on": "2026-10-26",
    "next_follow_up_on": "2026-11-04",
    "created_at": "2026-10-26T16:00:00Z",
    "updated_at": "2026-10-31T22:35:00Z"
  }
}
```

### Unregister a job application
Request (Delete Acme):
```
DELETE /api/v1/applications/8e34296b-ef82-4d6a-9b1f-30c49e3f0c11
```

Response:
```
204 No Content
```

Request After:
```
GET /api/v1/applications
```

Response After (Acme gone):
```
{
    "data": [
        {
            "id": "5213fcaf-15ea-4892-8419-d0d7077ebe16",
            "company_name": "CloudCo",
            "role_title": "Software Engineer",
            "status": "applied",
            "applied_on": "2026-10-22",
            "next_follow_up_on": "2026-10-29",
            "created_at": "2026-10-22T12:59:00Z",
            "updated_at": "2026-10-22T12:59:00Z"
        },
        {
            "id": "d02d9a11-f57b-46ab-9cf6-52d4c4992630",
            "company_name": "DataWorks",
            "role_title": "Infrastructure SWE",
            "status": "interviewing",
            "applied_on": "2026-10-24",
            "next_follow_up_on": "2026-10-31",
            "created_at": "2026-10-24T13:23:00Z",
            "updated_at": "2026-10-24T13:23:00Z"
        }
    ]
}
```

## Job Application Follow-ups
These define the API endpoints for job application follow-ups.
| Method | Endpoint | Purpose |
| -------| -------- | ------- |
| `POST` | `/api/v1/applications/{id}/follow-ups` | Schedule a follow-up |
| `GET` | `/api/v1/follow-ups` | List pending or due follow-ups |
| `GET` | `/api/v1/follow-ups/{id}` | Get a specific follow-up |
| `PATCH` | `/api/v1/follow-ups/{id}` | Reschedule or update a follow-up |
| `POST` | `/api/v1/follow-ups/{id}/complete` | Mark it completed |

Supports the following list filters
```
support ?due=overdue, ?due=today, and ?status=pending.
```

### Schedule a follow-up

### List pending or due follow-ups

### Get a specific follow up

### Mark follow-up as complete

## Dashboard and AI
These define the API endpoints for the dashboard and AI-integrated message generation.
| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/api/v1/dashboard` | Return application counts and today's tasks |
| `POST` | `/api/v1/applications/{id}/generate-follow-up` | Generate a message draft |

Note: Implement the dashboard and message-generation endpoints after the basic CRUD and scheduling flows work. They depend on those underlying capabilities.

### Return application counts and today's tasks

### Generates a message draft