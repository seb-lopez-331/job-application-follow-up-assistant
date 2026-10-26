
# API Contract



## Job Applications
These define the API endpoints for registered job applications.
| Method | Endpoint | Purpose | Returns |
| -------| -------- | ------- | ------- |
| `POST` | `/api/v1/applications` | Register a job application | `201`, `400`, `401`, `403`, `409` |
| `GET` | `/api/v1/applications` | List job applications with filters | `200`, `401`, `403` |
| `GET` | `/api/v1/applications/{id}` | Get job application by ID | `200`, `401`, `403`, `404` |
| `PATCH` | `/api/v1/applications/{id}` | Update selected fields | `200`, `400`, `401`, `403`, `404` |
| `DELETE` | `/api/v1/applications/{id}` | Unregister a job application | `204`, `401`, `403`, `404` |

### Register a job application
Registers a job application into the database. Returns `201 Created` or `400 Bad Request` if the client is authenticated and the database has no such job, depending on whether the JSON is passed in with the following structure:

```json
{
    "user": "johndoe",
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

#### Fields Specification
| Field | Type | Required | Description |
| :--- | :--- | :---: | :--- |
| `applied_on` | String | Yes | The day when the user applied to the job. |
| `company_name` | String | Yes | The name of the company the user is applying to. |
| `location` | String | Yes | Where the job will be performed. |
| `role_title` | String | Yes | The name of the role the user is applying to. |
| `user` | String | Yes | The user's username. |
| `work_arrangement` | String | Yes | Accepts the values of `remote`, `hybrid`, or `on-site`. |
| `description` | String | No | The description of the job the user is applying to. |
| `notes` | String | No | Any notes the user may record about the job. |
| `job_url` | String | No | The URL of the job posting. |
| `posted_at` | String | No | The date when the job was posted. |
| `status` | String | No | The status of the job application. Defaults to `applied`. |

Returns `401 Unauthorized` if the client is unauthenticated. Returns `403 Forbidden` if the client makes a request for a different user, unless the client is an admin. Returns `409 Conflict` if a job application with a matching user, company name, and role title are already in the database.

```
POST /api/v1/applications
{
    "user": "johndoe",
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
    "user": "johndoe",
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
Lists all job applications registered with the user, or all users if the user is an admin. If query parameters are specified, the list can be filtered by job application stage, job role title, number of entries, and so forth:

```
/api/v1/applications?stage=applied&search=platform&limit=20&cursor=...
```

Additionally, users can filter the list by user. Returns `200 OK` if client is authenticated, `401 Unauthorized` otherwise, and `403 Forbidden` if the client tries to filter by a different username.

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
            "user": "johndoe",
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
            "user": "johndoe",
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
            "user": "johndoe",
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
            "user": "johndoe",
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
Looks up a registered application by its UUID. Returns `200 OK` if an application with the specified UUID is found for the user, `403 Forbidden` if a non-admin user attempts to lookup an application that doesn't belong to them, and a `404 Not Found` otherwise.

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
    "user": "johndoe",
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
Updates the job application's fields to what's specified in the input JSON. Returns `200 OK` only if there is at least one valid field in the JSON and no invalid fields, and a `400 Bad Request` otherwise. Here's a reference for the valid fields:

#### Fields Specification
| Field | Type  | Description |
| :--- | :--- | :---: | :--- |
| `applied_on` | String | The day when the user applied to the job. |
| `company_name` | String | The name of the company the user is applying to. |
| `location` | String |  Where the job will be performed. |
| `role_title` | String |  The name of the role the user is applying to. |
| `work_arrangement` | String | Accepts the values of `remote`, `hybrid`, or `on-site`. |
| `description` | String | The description of the job the user is applying to. |
| `notes` | String | Any notes the user may record about the job. |
| `job_url` | String | The URL of the job posting. |
| `posted_at` | String | The date when the job was posted. |
| `status` | String | The status of the job application. Defaults to `applied`. |

Returns `401 Unauthorized` if the client isn't authenticated, `403 Forbidden` if a non-admin user requests this for another user's job application, or `404 Not Found` if the job application ID cannot be found.

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
    "user": "johndoe",
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
Deletes a job application associated with the ID passed in the URL. Returns `204 No Content` only if an entry associated with that ID is found. Returns `401 Unauthorized` if the client isn't authenticated, `403 Forbidden` if a non-admin user requests to delete a job application associated with another user, and `404 Not Found` if no such job application matches the ID.

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
            "user": "johndoe",
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
            "user": "johndoe",
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
Request:
```
POST /api/v1/applications/8e34296b-ef82-4d6a-9b1f-30c49e3f0c11/follow-ups
{
  "due_on": "2026-11-04"
}
```

Response:
```
201 Created
{
  "data": {
    "id": "c66a3e81-27c0-4b10-b47d-b1f3b2f20b31",
    "application_id": "8e34296b-ef82-4d6a-9b1f-30c49e3f0c11",
    "due_on": "2026-11-04",
    "status": "pending",
    "message_draft": null,
    "sent_at": null,
    "completed_at": null
  }
}
```

### List pending or due follow-ups
Request:
```
GET /api/v1/follow-ups
```

Response:
```
200 OK
{
    "data": [
        {
            "id": "c66a3e81-27c0-4b10-b47d-b1f3b2f20b31",
            "application_id": "8e34296b-ef82-4d6a-9b1f-30c49e3f0c11",
            "due_on": "2026-11-04",
            "status": "pending",
            "message_draft": null,
            "sent_at": null,
            "completed_at": null
        },
        {
            "id": "3c12c4f2-e79e-4f16-bca9-a50273fc7e16",
            "application_id": "5213fcaf-15ea-4892-8419-d0d7077ebe16",
            "due_on": "2026-11-10",
            "status": "pending",
            "message_draft": null,
            "sent_at": null,
            "completed_at": null
        },        
        {
            "id": "e9e4c11e-16cf-49e8-96a6-d453d030c233",
            "application_id": "d02d9a11-f57b-46ab-9cf6-52d4c4992630",
            "due_on": "2026-10-31",
            "status": "pending",
            "message_draft": null,
            "sent_at": null,
            "completed_at": null
        },
    ]
}
```

### Get a specific follow up by ID
Request:
```
GET /api/v1/follow-ups/3c12c4f2-e79e-4f16-bca9-a50273fc7e16
```

Response:
```
200 OK
{
    "data":  {
        "id": "3c12c4f2-e79e-4f16-bca9-a50273fc7e16",
        "application_id": "5213fcaf-15ea-4892-8419-d0d7077ebe16",
        "due_on": "2026-11-10",
        "status": "pending",
        "message_draft": null,
        "sent_at": null,
        "completed_at": null
    }
}
```

### Reschedule or update a follow-up
Request:
```
PATCH /api/v1/follow-ups/e9e4c11e-16cf-49e8-96a6-d453d030c233
{
    "due_on": "2026-12-01"
}
```

Response:
```
200 OK
{
    "data": {
        "id": "e9e4c11e-16cf-49e8-96a6-d453d030c233",
        "application_id": "d02d9a11-f57b-46ab-9cf6-52d4c4992630",
        "due_on": "2026-12-01",
        "status": "pending",
        "message_draft": null,
        "sent_at": null,
        "completed_at": null
    }
}
```

Request:
```
PATCH /api/v1/follow-ups/e9e4c11e-16cf-49e8-96a6-d453d030c233
{
    "message_draft": "Hi Recruiter, Hope this week has been treating you well. Please hire me, I desperately need the money and my girlfriend that I worked so hard to impress will leave me if I don't buy her a Louis Vuitton purse."
}
```

Response:
```
200 OK
{
    "data": {
        "id": "e9e4c11e-16cf-49e8-96a6-d453d030c233",
        "application_id": "d02d9a11-f57b-46ab-9cf6-52d4c4992630",
        "due_on": "2026-12-01",
        "status": "cancelled",
        "message_draft": "Hi Recruiter, Hope this week has been treating you well. Please hire me, I desperately need the money and my girlfriend that I worked so hard to impress will leave me if I don't buy her a Louis Vuitton purse.",
        "sent_at": null,
        "completed_at": null
    }
}
```

### Mark follow-up as complete
Request:
```
POST /api/v1/follow-ups/e9e4c11e-16cf-49e8-96a6-d453d030c233/complete
```

Response:
```
200 OK
{
    "data": {
        "id": "e9e4c11e-16cf-49e8-96a6-d453d030c233",
        "application_id": "d02d9a11-f57b-46ab-9cf6-52d4c4992630",
        "due_on": "2026-12-01",
        "status": "cancelled",
        "message_draft": "Hi Recruiter, Hope this week has been treating you well. Please hire me, I desperately need the money and my girlfriend that I worked so hard to impress will leave me if I don't buy her a Louis Vuitton purse.",
        "sent_at": null,
        "completed_at": "2026-12-01T13:23:00Z"
    }
}
```

## Dashboard and AI
These define the API endpoints for the dashboard and AI-integrated message generation.
| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/api/v1/dashboard` | Return application counts and today's tasks |
| `POST` | `/api/v1/applications/{id}/generate-follow-up` | Generate a message draft |

Note: Implement the dashboard and message-generation endpoints after the basic CRUD and scheduling flows work. They depend on those underlying capabilities.

### Return application counts and today's tasks
Request:
```
GET /api/v1/dashboard
```

Response:
```
{
    "data": {
        "applied": 24,
        "follow_up": 5,
        "interviews": 3,
        "applications": [
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
        ], ...
    }
}

```

### Generates a message draft
Request:
```
POST /api/v1/applications/8e34296b-ef82-4d6a-9b1f-30c49e3f0c11/generate-follow-up
{
  "tone": "professional",
  "context": "I spoke with the recruiter last week."
}
```

Response:
```
200 OK
{
  "data": {
    "application_id": "8e34296b-ef82-4d6a-9b1f-30c49e3f0c11",
    "message_draft": "Hi [Recruiter Name], ...",
    "generated_at": "2026-10-26T16:00:00Z"
  }
}
```
