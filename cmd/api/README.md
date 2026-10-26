
# API Contract

## Job Applications
| Method | Endpoint | Purpose |
| - | - |
| `POST` | `/api/v1/applications` | Register a job application |
| `GET` | `/api/v1/applications` | List job applications with filters |
| `GET` | `/api/v1/applications/{id}` | Get one job application |
| `PATCH` | `/api/v1/applications/{id}` | Update selected fields |
| `DELETE` | `/api/v1/applications/{id}` | Unregister a job application |

List filters:
```
?status=applied&search=platform&limit=20&cursor=...
```

## Follow-ups
| Method | Endpoint | Purpose |
| - | - |
| `POST` | `/api/v1/applications/{id}/follow-ups` | Schedule a follow-up |
| `GET` | `/api/v1/follow-ups` | List pending or due follow-ups |
| `GET` | `/api/v1/follow-ups/{id}` | Get a specific follow-up |
| `PATCH` | `/api/v1/follow-ups/{id}` | Reschedule or update a follow-up |
| `POST` | `/api/v1/follow-ups/{id}/complete` | Mark it completed |

Supports the following list filters
```
support ?due=overdue, ?due=today, and ?status=pending.
```

## Dashboard and AI
| Method | Endpoint | Purpose |
| - | - |
| `GET` | `/api/v1/dashboard` | Return application counts and today's tasks |
| `POST` | `/api/v1/applications/{id}/generate-follow-up` | Generate a message draft |

Note: Implement the dashboard and message-generation endpoints after the basic CRUD and scheduling flows work. They depend on those underlying capabilities.