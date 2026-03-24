# CEN5035-Project

# TravelLog — Hotel & Flight Log Web App

A simple full-stack web application for logging hotels and flights, showing travel routes and hotel nights on a map, and generating month / quarter / year summaries.
This project is designed for a **4-week software engineering class project** and a **team of four**.

---

## Members

- Ao Wang (@Ao Wang, @William-A-Wang): Front-end
- Ray Chen (@rkc8626, @kzc-qubit) : Front-end
- Wentao Chen (@wentao-chen227) : Back-end
- Jianan Gu (@jianangu1214, @Griffin) : Back-end

> Note: Due to the initial Git configuration, some contributors may appear under multiple commit identities.

## 1. Project Goals

- Practice frontend + backend collaboration
- Build a realistic CRUD web app with visualization
- Keep scope **small, clear, and finishable in 4 weeks**

---

## 2. Core Features

### 1 Hotel Log

- Add / edit / delete hotel stays
- Fields:
  - Hotel name
  - City, country
  - Check-in date, check-out date
  - Total price
  - Notes (optional)
- List hotels in a table

### 2 Flight Log

- Add / edit / delete flights
- Fields:
  - Airline, flight number
  - From airport (IATA), to airport (IATA)
  - Departure time, arrival time
  - Total price
  - Notes (optional)
- List flights in a table

### 3 Map View

- Show flight routes as lines on a map
- Show hotel stays as markers
- Clicking a marker or route shows basic details

### 4 Time Summary

- Summary by **month / quarter / year**
- Show:
  - Number of flights
  - Number of hotels
  - Total nights
  - Total spending
- Simple charts or tables

---

## 3. Tech Stack (Simple)

Frontend

- React + JavaScript
- Map: Leaflet or Mapbox
- Charts: Chart.js or Recharts

Backend

- Gin + Golang
- Database: PostgreSQL
- ORM: GORM
- Auth: simple email + password + JWT

Dev

- Docker + docker-compose
- GitHub for collaboration

---

## 4. Repository Structure

```
CEN5035-Project/
├── README.md                 # This file — start with §11 for how to run
├── frontend/                 # React + Vite ([frontend/README.md](frontend/README.md))
└── travellog-backend/
    ├── docker-compose.yml
    └── backend/              # Go API ([travellog-backend/backend/README.md](travellog-backend/backend/README.md))
```

---

## 5. Basic Data Models

User

- id
- email
- password

Hotel

- id
- user_id
- hotel_name
- city
- country
- check_in
- check_out
- price
- notes

Flight

- id
- user_id
- airline
- flight_number
- from_airport
- to_airport
- depart_time
- arrive_time
- price
- notes

---

## 6. API Endpoints (Minimal)

Auth

- POST `/auth/register`
- POST `/auth/login`

Hotels

- GET `/hotels`
- POST `/hotels`
- PUT `/hotels/:id`
- DELETE `/hotels/:id`

Flights

- GET `/flights`
- POST `/flights`
- PUT `/flights/:id`
- DELETE `/flights/:id`

Map

- GET `/map`
  Returns hotel locations and flight routes

Summary

- GET `/summary?type=month|quarter|year`

---

## 7. Team Work Split (4 People)

Frontend Developer A

- Login / register pages
- Flight log UI

Frontend Developer B

- Hotel log UI (table + form)
- Map page and summary page

Backend Developer A

- Auth (login / register)
- Flight CRUD APIs

Backend Developer B

- Hotel CRUD APIs
- Map and summary aggregation APIs

---

## 8. 4-Week Plan (Easy Pace)

### Week 1

- Project setup
- Database schema
- Login / register
- Hotel CRUD (backend + frontend)

### Week 2

- Flight CRUD (backend + frontend)
- Basic tables for hotels and flights

### Week 3

- Map page (show flights + hotels)
- Summary API and simple UI

### Week 4

- Bug fixes
- UI polish
- Final demo preparation

---

## 9. Minimum Demo Requirements

- User can log in
- User can add a hotel
- User can add a flight
- Map shows at least one route and one hotel
- Summary page shows numbers for one month or year

---

## 10. Notes

- Accuracy is not the goal — completion is
- Map can use approximate city locations
- UI can be simple tables and forms

---

## 11. Local development and testing

### Ports (defaults)

| Service | Port | Notes |
|--------|------|--------|
| Frontend (Vite) | **5173** | `npm run dev` in `frontend/` |
| Backend (Gin) | **8080** | `SERVER_PORT` (default `8080`) |
| API from the browser | `/api/*` | Vite proxies to the backend; see below |

The React app calls `fetch('/api/...')`. Vite rewrites `/api` to the backend URL. **Keep the proxy target and `SERVER_PORT` in sync** (both default to port 8080).

If the backend runs on another port (e.g. `SERVER_PORT=8081`), create `frontend/.env.local`:

```bash
VITE_API_PROXY_TARGET=http://localhost:8081
```

Restart `npm run dev` after changing env files.

### 11.1 Run the backend (SQLite, recommended for local dev)

From the repo root:

```bash
cd travellog-backend/backend
go mod tidy
USE_SQLITE=true go run main.go
```

Health check: [http://localhost:8080/health](http://localhost:8080/health)

More detail, PostgreSQL/Docker, environment variables, and `curl` examples: [travellog-backend/backend/README.md](travellog-backend/backend/README.md).

### 11.2 Seed a test user (Cypress / QA)

There is no separate admin role; this creates a **normal user** with known credentials:

```bash
cd travellog-backend/backend
USE_SQLITE=true go run ./cmd/seedtestuser
```

Defaults: **email** `admin@test.local`, **password** `AdminTest123!`. Override with `SEED_EMAIL` and `SEED_PASSWORD` (password at least 8 characters). Safe to run multiple times.

### 11.3 Run the frontend

```bash
cd frontend
npm install
npm run dev
```

Open [http://localhost:5173](http://localhost:5173). Register a new account, or log in with the seeded test user after running the seed command.

### 11.4 Unit tests (Vitest)

```bash
cd frontend
npm test
npm run test:watch
```

### 11.5 End-to-end tests (Cypress)

Needs **both** backend and frontend running (same terminals as above).

```bash
# Terminal A
cd travellog-backend/backend && USE_SQLITE=true go run main.go

# Terminal B
cd frontend && npm run dev

# Terminal C
cd frontend && npm run cypress:open
# or headless: npm run cypress:run
```

Default Cypress credentials match the seed user (`admin@test.local` / `AdminTest123!`). Override when needed:

```bash
npx cypress run --env TEST_EMAIL=you@example.com,TEST_PASSWORD='YourPass123'
```

### 11.6 `listen tcp :8080: address already in use`

Another process (often a previous `go run main.go`) is using port 8080.

On macOS/Linux:

```bash
lsof -i :8080 -sTCP:LISTEN
kill <PID>
```

Or stop the other terminal where the backend is still running. Alternatively run the backend on another port and set `VITE_API_PROXY_TARGET` as in the table above.

### 11.7 Repository layout (actual)

- `frontend/` — React + Vite app
- `travellog-backend/backend/` — Go + Gin API (`main.go`, `cmd/seedtestuser/`)

---

## 12. License

MIT
