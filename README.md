
# API Contract

All endpoints use the `/api/v1` prefix. Unless otherwise noted, request and
response bodies use JSON. Protected endpoints require an authenticated client.

## Job Applications
These endpoints manage registered job applications.

| Method | Endpoint | Purpose | Responses |
| -------| -------- | ------- | ------- |
| `POST` | `/api/v1/applications` | Register a job application | `201`, `400`, `401`, `403`, `409` |
| `GET` | `/api/v1/applications` | List job applications with filters | `200`, `401`, `403` |
| `GET` | `/api/v1/applications/{id}` | Get job application by ID | `200`, `400`, `401`, `403`, `404` |
| `PATCH` | `/api/v1/applications/{id}` | Update selected fields | `200`, `400`, `401`, `403`, `404` |
| `DELETE` | `/api/v1/applications/{id}` | Unregister a job application | `204`, `400`, `401`, `403`, `404` |

### Register a job application
Registers a job application. The request body must be valid JSON and include
all required fields listed below. Returns `201 Created` on success,
`400 Bad Request` if the body is invalid, `401 Unauthorized` if the client is
unauthenticated, `403 Forbidden` if a non-admin client attempts to create an
application for another user, or `409 Conflict` if an application already
exists for the same user, company, and role.

```json
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

#### Field specification
| Field | Type | Required | Description |
| :--- | :--- | :---: | :--- |
| `company_name` | String | Yes | Company name. |
| `role_title` | String | Yes | Job title. |
| `location` | String | Yes | Job location. |
| `work_arrangement` | String | Yes | One of `remote`, `hybrid`, or `on-site`. |
| `applied_on` | String (date) | Yes | Date the user applied, in `YYYY-MM-DD` format. |
| `description` | String | No | Job description. |
| `notes` | String | No | Notes about the application. |
| `job_url` | String (URI) | No | URL of the job posting. |
| `posted_at` | String (date) | No | Date the job was posted, in `YYYY-MM-DD` format. |
| `status` | String | No | Application status; defaults to `applied`. |
| `user` | String | Admin-only | Username of the user the application belongs to. |

The `user` field may be supplied only by an admin. It is included in the
response only when the requester is an admin.

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
Status: 201 Created
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

Admin request:
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

Response to a non-admin request that includes `user`:
```
403 Forbidden
```

Response to an admin request:
```
Status: 201 Created
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
Lists the authenticated user's applications. Admins may list applications
across users. Supported query parameters are:

```
GET /api/v1/applications?status=applied&search=platform&limit=20&cursor=eyJvZmZzZXQiOjIwfQ
```

| Parameter | Description |
| :--- | :--- |
| `status` | Filter by application status. |
| `search` | Search application company names and role titles. |
| `limit` | Maximum number of results to return. |
| `cursor` | Continue listing from a cursor returned by the API. |
| `user` | Admin-only username to filter by. |

List responses include a top-level `next_cursor` field. It is `null` when there
are no more results; otherwise, pass its value as the `cursor` query parameter
to retrieve the next page.

Returns `200 OK` if the client is authenticated, `401 Unauthorized` otherwise,
or `403 Forbidden` if a non-admin client tries to filter by another username.

Request:
```
GET /api/v1/applications
```

Response:
```
Status: 200 OK
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
    ],
    "next_cursor": null
}
```

Request with filter:
```
GET /api/v1/applications?status=interviewing
```

Response with filter:
```
Status: 200 OK
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
    ],
    "next_cursor": null
}
```

### Get a job application by ID
Looks up an application by its UUID. Returns `200 OK` if the application is
found and belongs to the user (or the requester is an admin), `400 Bad Request`
if the UUID is invalid, `401 Unauthorized` if the client is unauthenticated, 
`403 Forbidden` if a non-admin user attempts to access another user's application,
or `404 Not Found` if no such application exists.

Request:
```
GET /api/v1/applications/5213fcaf-15ea-4892-8419-d0d7077ebe16
```

Response:
```
Status: 200 OK
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
Updates the job application's fields using the supplied JSON object. Include
at least one supported field; unknown fields, an invalid body, or an invalid UUID 
result in `400 Bad Request`. A successful update returns `200 OK`.

#### Field specification
| Field | Type  | Description |
| :--- | :--- | :--- |
| `applied_on` | String (date) | Date the user applied, in `YYYY-MM-DD` format. |
| `company_name` | String | Company name. |
| `location` | String | Job location. |
| `role_title` | String | Job title. |
| `work_arrangement` | String | One of `remote`, `hybrid`, or `on-site`. |
| `description` | String | Job description. |
| `notes` | String | Notes about the application. |
| `job_url` | String (URI) | URL of the job posting. |
| `posted_at` | String (date) | Date the job was posted, in `YYYY-MM-DD` format. |
| `status` | String | Application status. |

Returns `401 Unauthorized` if the client is unauthenticated, `403 Forbidden`
if a non-admin user requests another user's application, or `404 Not Found` if
the application ID does not exist.

Request:
```
PATCH /api/v1/applications/8e34296b-ef82-4d6a-9b1f-30c49e3f0c11
{
  "status": "rejected"
}
```

Response:
```
Status: 200 OK
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

### Delete a job application
Deletes the job application identified by the UUID in the URL. Returns
`204 No Content` on success, `400 Bad Request` if the UUID is invalid,
`401 Unauthorized` if the client is unauthenticated, `403 Forbidden` if 
a non-admin user attempts to delete another user's application, or `404 Not Found` 
if the application does not exist.

Request:
```
DELETE /api/v1/applications/8e34296b-ef82-4d6a-9b1f-30c49e3f0c11
```

Response:
```
Status: 204 No Content
```

Follow-up request:
```
GET /api/v1/applications
```

Follow-up response:
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

## Job Application Follow-Ups
These endpoints manage job application follow-ups.

| Method | Endpoint | Purpose | Responses |
| -------| -------- | ------- | ------- |
| `POST` | `/api/v1/applications/{id}/follow-ups` | Schedule a follow-up | `201`, `400`, `401`, `403` |
| `GET` | `/api/v1/follow-ups` | List pending or due follow-ups | `200`, `401`, `403` |
| `GET` | `/api/v1/follow-ups/{id}` | Get a specific follow-up | `200`, `400`, `401`, `403`, `404` |
| `PATCH` | `/api/v1/follow-ups/{id}` | Reschedule or update a follow-up | `200`, `400`, `401`, `403`, `404` |
| `POST` | `/api/v1/follow-ups/{id}/complete` | Mark it complete | `200`, `400`, `401`, `403`, `404`, `409` |

### Schedule a follow-up
Creates a follow-up for the application identified by `{id}`. Returns
`201 Created` on success, `400 Bad Request` if the request body is invalid,
`401 Unauthorized` if the client is unauthenticated, or `403 Forbidden` if the
application belongs to another user and the requester is not an admin.

```json
{
    "due_on": "2026-11-04"
}
```

#### Field specification
| Field | Type | Required | Description |
| :--- | :--- | :---: | :--- |
| `due_on` | String (date) | Yes | Date the follow-up is due, in `YYYY-MM-DD` format. |
| `message_draft` | String or `null` | No | Draft follow-up message. |

Request:
```
POST /api/v1/applications/8e34296b-ef82-4d6a-9b1f-30c49e3f0c11/follow-ups
{
    "due_on": "2026-11-04"
}
```

Response:
```
Status: 201 Created
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
Lists follow-ups associated with the authenticated user. Admins may list
follow-ups across users. The list can be filtered by due date, status, or (for
admins only) username:

```
GET /api/v1/follow-ups?due=overdue
GET /api/v1/follow-ups?due=today&status=pending
GET /api/v1/follow-ups?user=johndoe
```

`due` accepts `overdue` or `today`. Returns `200 OK` if the client is
authenticated, `401 Unauthorized` if unauthenticated, or `403 Forbidden` if a
non-admin client filters by another user.

Request:
```
GET /api/v1/follow-ups
```

Response:
```
Status: 200 OK
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
        }
    ]
}
```

### Get a specific follow-up by ID
Retrieves a follow-up by its UUID. Returns `200 OK` if it is found and belongs
to the user (or the requester is an admin), `400 Bad Request` if the UUID is invalid,
`401 Unauthorized` if the client is unauthenticated, `403 Forbidden` if a non-admin 
user tries to access another user's follow-up, or `404 Not Found` if the follow-up 
does not exist.

Request:
```
GET /api/v1/follow-ups/3c12c4f2-e79e-4f16-bca9-a50273fc7e16
```

Response:
```
Status: 200 OK
{
    "data": {
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
Updates the follow-up identified by the UUID in the URL. Include at least one
supported field; an empty or invalid body, or invalid UUID returns `400 Bad Request`. A
successful update returns `200 OK`.

#### Field specification
| Field | Type | Required | Description |
| :--- | :--- | :---: | :--- |
| `due_on` | String (date) | No | New due date, in `YYYY-MM-DD` format. |
| `message_draft` | String or `null` | No | Follow-up message draft. |

Returns `401 Unauthorized` if the client is unauthenticated, `403 Forbidden`
if a non-admin user requests another user's follow-up, or `404 Not Found` if
the follow-up ID does not exist.

Request:
```
PATCH /api/v1/follow-ups/e9e4c11e-16cf-49e8-96a6-d453d030c233
{
    "due_on": "2026-12-01"
}
```

Response:
```
Status: 200 OK
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
    "message_draft": "Hello [Recruiter Name], I hope you are doing well. I wanted to follow up on my application and ask whether there are any updates. Thank you for your time."
}
```

Response:
```
Status: 200 OK
{
    "data": {
        "id": "e9e4c11e-16cf-49e8-96a6-d453d030c233",
        "application_id": "d02d9a11-f57b-46ab-9cf6-52d4c4992630",
        "due_on": "2026-12-01",
        "status": "pending",
        "message_draft": "Hello [Recruiter Name], I hope you are doing well. I wanted to follow up on my application and ask whether there are any updates. Thank you for your time.",
        "sent_at": null,
        "completed_at": null
    }
}
```

### Mark follow-up as complete
Marks the follow-up identified by the UUID in the URL as complete. Returns
`200 OK` on success, `400 Bad Request` if the UUID is invalid, `401 Unauthorized` 
if the client is unauthenticated, `403 Forbidden` if a non-admin user attempts 
to complete another user's follow-up, `404 Not Found` if the follow-up does not 
exist, or `409 Conflict` if it is already complete.

Request:
```
POST /api/v1/follow-ups/e9e4c11e-16cf-49e8-96a6-d453d030c233/complete
```

Response:
```
Status: 200 OK
{
    "data": {
        "id": "e9e4c11e-16cf-49e8-96a6-d453d030c233",
        "application_id": "d02d9a11-f57b-46ab-9cf6-52d4c4992630",
        "due_on": "2026-12-01",
        "status": "completed",
        "message_draft": "Hello [Recruiter Name], I hope you are doing well. I wanted to follow up on my application and ask whether there are any updates. Thank you for your time.",
        "sent_at": null,
        "completed_at": "2026-12-01T13:23:00Z"
    }
}
```

## Dashboard and AI
These endpoints provide dashboard data and AI-assisted message generation.

| Method | Endpoint | Purpose | Responses |
| ------ | -------- | ------- | ------- |
| `GET` | `/api/v1/dashboard` | Return application counts and application summaries. | `200`, `401` |
| `POST` | `/api/v1/applications/{id}/generate-follow-up` | Generate a message draft. | `200`, `400`, `401`, `403`, `404` |

### Return the dashboard summary
Returns application counts and application summaries for the dashboard. Returns
`200 OK` if the client is authenticated or `401 Unauthorized` otherwise.

Request:
```
GET /api/v1/dashboard
```

Response:
```
Status: 200 OK
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
        ]
    }
}
```

### Generate a message draft
Generates a message draft using an LLM. The request body must be valid JSON and
include both required fields. Returns `200 OK` on success, `400 Bad Request`
if the body or UUID are invalid, `401 Unauthorized` if the client is unauthenticated,
`403 Forbidden` if a non-admin user attempts to generate a draft for another
user's application, or `404 Not Found` if the application does not exist.

```json
{
    "tone": "professional",
    "context": "I spoke with the recruiter last week."
}
```

#### Field specification
| Field | Type | Required | Description |
| :--- | :--- | :---: | :--- |
| `tone` | String | Yes | Tone to use for the generated message. |
| `context` | String | Yes | Context for the generated message. |

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
Status: 200 OK
{
  "data": {
    "application_id": "8e34296b-ef82-4d6a-9b1f-30c49e3f0c11",
    "message_draft": "Hello [Recruiter Name], I hope you are doing well. I wanted to follow up on my application.",
    "generated_at": "2026-10-26T16:00:00Z"
  }
}
```

## User Account Management
These API methods provide means to create new, and access and update existing users.
| Method | Endpoint | Purpose | Responses |
| ------ | -------- | ------- | ------- |
| `POST` | `/api/v1/users` | Register a new account. | `201`, `400`, `409` |
| `GET` | `/api/v1/users/me` | Get current profile. | `200`, `401` |
| `PATCH` | `/api/v1/users/me` | Update profile. | `200`, `400`, `401`|

### Register a New Account
Registers a new user into the database based on the JSON response body. Accepts email,
password, and profile details:

```json
{
    "username": "jdoughnut25",
    "email": "john.doe@example.com",
    "password": "notarealpassword",
    "preferences": {
        "follow_up_schedule": {
            "follow_ups": [3, 7, 14], // days
            "skip_weekends": true,
            "skip_holidays": true,
        },
        "mode": "dark",
    }
}
```

#### Field Specification
| Field | Type | Required | Description |
| :--- | :--- | :---: | :--- |
| `username` | String | Yes | The user's username. |
| `email` | String | Yes | The user's email address. |
| `password` | String | Yes | The user's password. |
| `preferences` | Object | No | The user's preferences. |

The password must be salted and hashed before being stored in the database. It
also must not be returned in the response. Returns `201 Created` on success,
`400 Bad Request` if the request body doesn't match the schema above, and
`409 Conflict` if a new user with the same username 

Request:
```
POST /api/v1/users
{
    "username": "jdoughnut25"
    "email": "john.doe@example.com",
    "password": "notarealpassword",
    "preferences": {
        "follow_up_schedule": {
            "follow_ups": [3, 7, 14], // days
            "skip_weekends": true,
            "skip_holidays": true,
        },
        "mode": "dark",
    }
}
```

Response:
```
Status: 201 Created
{
    "data": {
        "id": "723a5720-3bf4-4ecd-9b2e-a17eb509ebd8",
        "username": "jdoughnut25",
        "email": "john.doe@example.com",
        "preferences": {
            "follow_up_schedule": {
                "follow_ups": [3, 7, 14], // days
                "skip_weekends": true,
                "skip_holidays": true,
            },
            "mode": "dark",
        },
        "created_at": "2026-10-26T23:48:00Z",
        "updated_at": "2026-10-26T23:48:00Z"
    }
}
```

### Get Current Profile
Allows the logged-in user to see their own account details. Returns `200 OK` if
authenticated and `401 Unauthorized` otherwise.

Request:
```
GET /api/v1/users/me
```

Response:
```
Status: 200 OK
{
    "data": {
        "id": "723a5720-3bf4-4ecd-9b2e-a17eb509ebd8",
        "username": "jdoughnut25",
        "email": "john.doe@example.com",
        "preferences": {
            "follow_up_schedule": {
                "follow_ups": [3, 7, 14], // days
                "skip_weekends": true,
                "skip_holidays": true,
            },
            "mode": "dark",
        },
        "created_at": "2026-10-26T23:48:00Z",
        "updated_at": "2026-10-26T23:48:00Z"
    }
}
```

### Update Profile
Allows user to update email, change password, or edit preferences.
Here are the following fields that can be updated:

#### Field Specification
| Field | Type | Required | Description |
| :--- | :--- | :---: | :--- |
| `email` | String | No | The user's email address. |
| `password` | String | No | The user's password. |
| `preferences` | Object | No | The user's preferences. |

At least one of these fields is required. Returns `200 OK` on success, 
`400 Bad Request` if no fields are present or if an invalid field is
present in the JSON, or `401 Unauthorized` if the client is not
authenticated.

Request:
```
PATCH /api/v1/users/me
{
    "preferences": {
        "mode": "light"
    }
}
```

Response:
```
Status: 200 OK
{
    "data": {
        "id": "723a5720-3bf4-4ecd-9b2e-a17eb509ebd8",
        "username": "jdoughnut25",
        "email": "john.doe@example.com",
        "preferences": {
            "follow_up_schedule": {
                "follow_ups": [3, 7, 14], // days
                "skip_weekends": true,
                "skip_holidays": true,
            },
            "mode": "light",
        },
        "created_at": "2026-10-26T23:48:00Z",
        "updated_at": "2026-10-26T23:48:00Z"
    }
}
```

## Authentication & Session Management
These endpoints facilitate authentication and session management for users.
| Method | Endpoint | Purpose | Responses |
| ------ | -------- | ------- | ------- |
| `POST` | `/api/v1/auth/login` | Authenticate a user. | `200`, `400`, `401` |
| `POST` | `/api/v1/auth/logout` | End session. | `204` |

### Authenticating a User
This endpoint authenticates a user into the application. Returns `200 OK` on
success, `400 Bad Request` if either `login` or `password` are not provided, 
and `401 Unauthorized` if invalid credentials are provided. 

#### Field Specification
| Field | Type | Required | Description |
| :--- | :--- | :---: | :--- |
| `login` | String | No | The user's email address or username. |
| `password` | String | No | The user's password. |

Request:
```
POST /api/v1/auth/login
{
    "login": "jdoughnut25",
    "password": "notarealpassword"
}
```

Response:
```
Status: 200 OK
Set-Cookie: session_token=s3cr3t_s3ssi0n_str1ng_h3r3; Path=/; HttpOnly; Secure; SameSite=Lax; Expires=Tue, 27 Oct 2026 01:09:00 GMT
Content-Type: application/json
```

### Ending a Session
This endpoint ends a user session. Always returns `204 No Content`.

Request:
```
POST /api/v1/auth/logout
```

Response:
```
Status: 204 No Content
```

## Error Formats
A consistent JSON error format should be used for every endpoint:

```json
{
  "error": {
    "code": "application_not_found",
    "message": "The requested application was not found."
  }
}
```

These should generally be the HTTP status codes and when to use them:
| Status | When to use it |
| ------ | -------------- |
| `400 Bad Request` | Malformed JSON, invalid UUID, or invalid query parameters |
| `401 Unauthorized` | 	Missing or invalid authentication |
| `403 Forbidden` | Resources don't belong to the authenticated user |
| `404 Not Found` | Resource doesn't exist |
| `409 Conflict` | A genuine state confict, such as an invalid duplicate operation |
| `422 Unprocessable Entity` | Valid JSON that violates a business rule |
| `429 Too Many Requests` | Rate limit exceeded |
| `500 Internal Server Error` | Unexpected server side failure |
