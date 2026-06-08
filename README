# Event App API

> REST API for event and attendee management written in Go.

## Stack

```text
Go
Gin
PostgreSQL
JWT
Swagger
Air
```

## Setup

### Clone

```bash
git clone https://github.com/SV5623/go_learning.git
cd go_learning/back
```

### Environment

Create:

```text
cmd/api/.env
```

```env
PORT=5623
DATABASE_URL=postgres://postgres:postgres@localhost:5432/event_app?sslmode=disable
JWT_SECRET=your-secret-key
```

### Run

```bash
cd cmd/api
air
```

Server:

```text
http://localhost:5623
```

---

## Migrations

Create migration:

```bash
migrate create -ext sql -dir ./cmd/migrate/migrations -seq create_attendees_table
```

Apply:

```bash
go run cmd/migrate/main.go up
```

Rollback:

```bash
go run cmd/migrate/main.go down
```

Force version:

```bash
migrate force <version>
```

Example:

```bash
migrate force 2
```

---

## Swagger

Install:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Generate:

```bash
swag init --dir cmd/api --parseDependency --parseInternal --parseDepth 1
```

Open:

```text
http://localhost:5623/swagger/index.html
```

---

## Routes

### Auth

```http
POST /api/v1/auth/register
POST /api/v1/auth/login
```

### Events

```http
GET    /api/v1/events
GET    /api/v1/events/:id

POST   /api/v1/events
PUT    /api/v1/events/:id
DELETE /api/v1/events/:id
```

### Attendees

```http
GET    /api/v1/events/:id/attendees
GET    /api/v1/attendees/:id/events

POST   /api/v1/events/:id/attendees/:userId
DELETE /api/v1/events/:id/attendees/:userId
```

---

## Authorization

```http
Authorization: Bearer <jwt-token>
```

---

## Status

```text
[x] Authentication
[x] JWT Authorization
[x] Event CRUD
[x] Attendee Management
[x] PostgreSQL
[x] Swagger

[ ] Frontend (React)
[ ] Docker Compose
```
