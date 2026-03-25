# TravelLog Backend

Teammates: for running the **frontend + backend together**, tests, seed user, and port troubleshooting, see the repo root **[README.md](../../README.md)** (section **11. Local development and testing**).

## Prerequisites

- Go 1.21 or higher
- PostgreSQL (optional, see below)
- Docker (optional, for PostgreSQL)

## Quick Start

### Option 1: Using SQLite (Recommended for Development)

For local development without Docker, use SQLite:

```bash
cd backend
go mod tidy
USE_SQLITE=true go run main.go
```

The server listens on `SERVER_PORT` (default **8080**): `http://localhost:8080`

If you see `listen tcp :8080: bind: address already in use`, another process is using that port. Find and stop it, for example:

```bash
lsof -i :8080 -sTCP:LISTEN
kill <PID>
```

Or run on another port (remember to point the frontend proxy at the same port; see repo root **README** section 11):

```bash
SERVER_PORT=8081 USE_SQLITE=true go run main.go
```

### Test user for Cypress / manual QA

There is no separate “admin” role in the API; this command adds a **normal user** with fixed credentials for automated tests:

```bash
cd backend
USE_SQLITE=true go run ./cmd/seedtestuser
```

Defaults: email `admin@test.local`, password `AdminTest123!`. Override with `SEED_EMAIL` and `SEED_PASSWORD` (password must be at least 8 characters). The command is idempotent: if the user already exists, it does nothing.

### Option 2: Using PostgreSQL with Docker

1. Start PostgreSQL using docker-compose:

```bash
docker-compose up -d
```

2. Run the backend:

```bash
cd backend
go mod tidy
go run main.go
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_HOST | localhost | Database host |
| DB_PORT | 5432 | Database port |
| DB_USER | travellog | Database user |
| DB_PASSWORD | travellog123 | Database password |
| DB_NAME | travellog | Database name |
| JWT_SECRET | travellog-secret-key-2024 | JWT signing secret |
| SERVER_PORT | 8080 | Server port |
| USE_SQLITE | false | Set to "true" to use SQLite |

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /auth/register | Register a new user |
| POST | /auth/login | Login and get JWT token |

### Hotels (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /hotels | Get all hotels for user |
| GET | /hotels/:id | Get specific hotel |
| POST | /hotels | Create new hotel |
| PUT | /hotels/:id | Update hotel |
| DELETE | /hotels/:id | Delete hotel |

### Flights (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /flights | Get all flights for user |
| GET | /flights/:id | Get specific flight |
| POST | /flights | Create new flight |
| PUT | /flights/:id | Update flight |
| DELETE | /flights/:id | Delete flight |

## Example Requests

### Register

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}'
```

### Create Hotel (with token)

```bash
curl -X POST http://localhost:8080/hotels \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "hotel_name": "Grand Hotel",
    "city": "Paris",
    "country": "France",
    "check_in": "2024-06-01",
    "check_out": "2024-06-05",
    "price": 500.00,
    "notes": "Great view of the Eiffel Tower",
    "latitude": 48.8566,
    "longitude": 2.3522
  }'
```

### Get All Hotels

```bash
curl -X GET http://localhost:8080/hotels \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```
