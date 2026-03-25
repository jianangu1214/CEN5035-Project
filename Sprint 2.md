# Sprint 2

## This document contains the User Stories and Issues specific to Sprint 2 only. It does not represent the complete project scope. User Stories and Issues for other features (Map View, Summary) will be added in subsequent sprints.

# User Stories

## 1. Add Flight Record

As a logged-in user, I want to create a flight record so I can log my trips.

## 2. View Flight List

As a logged-in user, I want to view all my flights in a list so I can quickly browse history.

## 3. Edit Flight Record

As a logged-in user, I want to edit an existing flight record to correct or update information.

## 4. Delete Flight Record

As a logged-in user, I want to delete a flight record to remove entries I no longer need.

## 5. Integrated Auth + Data Access

As a user, I want frontend and backend to use JWT so that only my authenticated session can access my hotel and flight data.

## 6. Basic Frontend Test Coverage

As a developer, I want minimal E2E (Cypress) and unit tests to guard core flows.

---

# Issues Planned for Sprint 2

## Frontend Issues

### Ao Wang (GitHub: Ao Wang, William-A-Wang)

- [Frontend] Flight list page (read/render `/flights`, with empty/loading/error states).
- [Frontend] Flight create/edit form (call POST/PUT `/flights`).
- [Frontend] JWT integration: persist token after login; gated access to hotel/flight pages.

### Ray Chen (GitHub: rkc8626, kzc-qubit)

- [Frontend] Cypress smoke test (e.g., login + create flight or hotel).
- [Frontend] Component/function unit tests, target 1:1 coverage for key pieces (FlightList, FlightForm, PrivateRoute, etc.).

## Backend Issues

### Wentao Chen (GitHub: wentao-chen227)

- [Backend] DB/config alignment for flights: ensure docker-compose/config support flight table; migrations/auto-migrate.
- [Backend] Auth/JWT support and backend unit tests (auth middleware, handlers), target 1:1 coverage of core logic.

### Jianan Gu (GitHub: jianangu1214, Griffin)

- [Backend] Flight model refinement and validation.
- [Backend] Implement Flight CRUD APIs (`GET/POST/PUT/DELETE /flights`) and add routes with JWT protection.
- [Backend] Write/update API documentation in this file for Auth, Hotel, Flight (request/response/error examples).

---

# Sprint 2 Completion Status

## Issues the team planned to address

- **Frontend (Ao Wang):** Flight list page; Flight form page; JWT integration.
- **Frontend (Ray Chen):** Cypress smoke test; frontend unit tests (1:1 for key components).
- **Backend (Wentao Chen):** DB/config for flights; Auth/JWT support; backend integration tests.
- **Backend (Jianan Gu):** Flight CRUD; API documentation update.

## Successfully completed

- Frontend:
  - Flight list and form pages wired to backend; JWT persisted for protected routes.
  - Cypress E2E: login + create hotel flow, and flight navigation/add form flow.
  - Unit tests (Vitest/RTL): PrivateRoute, HotelList, HotelForm, FlightList, FlightForm.
- Backend:
  - Flight model finalized; Flight CRUD handlers implemented and routed with JWT protection.
  - API documentation added (Auth/Hotel/Flight) in this file.
  - JWT-based authentication verified on protected routes (`/me`, `/flights`, `/hotels`).
  - PostgreSQL via Docker configured and connectivity verified.
  - Backend integration tests (routes) covering:
    - user registration and login
    - JWT validation
    - protected route access
    - flight creation, retrieval, and deletion
    - hotel creation and retrieval
    - user-level data isolation
    - invalid input handling

## Not completed (and why)

- all issues in Sprint 2 have been finished

---

# Backend Unit Tests (Sprint 2)

- `routes/router_test.go`
  - Auth: register/login, JWT validation, `/me` access.
  - Flights: create, list (user isolation), delete, invalid time input, auth required.
  - Hotels: create, list (user isolation).
- Target ratio: ~1:1 for handlers’ main paths (auth, flights, hotels). Add more cases if time permits (e.g., update, 404 not found).
- How to run (full and targeted):
  1. Start Postgres (default creds `travellog/travellog123` on `localhost:5432`):
     - `cd travellog-backend && docker-compose up -d`
     - verify: `docker ps | grep travellog`
  2. Run all backend tests:
     - `cd travellog-backend/backend`
     - `go test ./routes -v`
     - If cache/permission issues: `GOMODCACHE=$PWD/.gocache GOCACHE=$PWD/.gocache go test ./routes -v`
  3. Run a single test case (example: create flight):
     - `go test ./routes -run TestCreateFlight -v`

---

# Frontend Tests (Sprint 2)

- Cypress E2E
  - `cypress/e2e/smoke.cy.js`: logs in (uses `TEST_EMAIL`/`TEST_PASSWORD` env), creates a hotel, verifies it appears in the list.
  - `cypress/e2e/flight.cy.js`: logs in (demo creds in test), opens Flights page, opens Add Flight form, checks required fields render.
- Unit tests (Vitest + Testing Library)
  - `src/components/PrivateRoute.test.jsx`: protects routes without token.
  - `src/pages/HotelList.test.jsx`: loads and renders hotels; delete flow calls API.
  - `src/pages/HotelForm.test.jsx`: submits hotel form and navigates back to list.
  - `src/pages/FlightList.test.jsx`: renders flights; delete flow calls API.
  - `src/pages/FlightForm.test.jsx`: submits flight creation with RFC3339 times and uppercase IATA codes.
- How to run:
  - Unit: `cd frontend && npm test`
  - Cypress (headless): `cd frontend && npm run cypress:run` (set `TEST_EMAIL`/`TEST_PASSWORD`)
  - Cypress (headed): `cd frontend && npm run cypress:open`

---

# Backend API Documentation

## Auth

- `POST /auth/register`
  - Body: `{ "email": "user@example.com", "password": "min 8 chars" }`
  - Responses: `201 { "id": 1, "email": "user@example.com" }`; `409` email exists; `400` invalid input.
- `POST /auth/login`
  - Body: `{ "email": "...", "password": "..." }`
  - Responses: `200 { "token": "<jwt>", "user": { "id": 1, "email": "..." } }`; `401` invalid credentials; `400` invalid input.
- `GET /me` (JWT required)
  - Header: `Authorization: Bearer <token>`
  - Responses: `200 { "user_id": 1 }`; `401` missing/invalid token.

## Hotels (JWT required)

- `GET /hotels`
  - Response: `200 { "hotels": [ { id, hotel_name, city, country, check_in, check_out, price, notes, latitude, longitude } ] }`
- `GET /hotels/:id`
  - Response: `200 { "hotel": { ... } }`; `404` not found (or not owned).
- `POST /hotels`
  - Body: `{ "hotel_name": "...", "city": "...", "country": "...", "check_in": "YYYY-MM-DD", "check_out": "YYYY-MM-DD", "price": 123.45, "notes": "...", "latitude": 0, "longitude": 0 }`
  - Responses: `201 { "hotel": { ... } }`; `400` invalid body/date; `401` unauthorized.
- `PUT /hotels/:id`
  - Body: partial fields same as create.
  - Responses: `200 { "hotel": { ...updated } }`; `404` not found; `400` invalid body.
- `DELETE /hotels/:id`
  - Responses: `200 { "message": "Hotel deleted successfully" }`; `404` not found.

## Flights (JWT required)

- `GET /flights`
  - Response: `200 { "flights": [ { id, airline, flight_number, from_airport, to_airport, depart_time (RFC3339), arrive_time (RFC3339), price, notes } ] }`
- `GET /flights/:id`
  - Response: `200 { "flight": { ... } }`; `404` not found (or not owned).
- `POST /flights`
  - Body: `{ "airline": "...", "flight_number": "...", "from_airport": "JFK", "to_airport": "LAX", "depart_time": "2026-03-01T10:00:00Z", "arrive_time": "2026-03-01T14:00:00Z", "price": 320.5, "notes": "..." }`
  - Responses: `201 { "flight": { ... } }`; `400` invalid body/time; `401` unauthorized.
- `PUT /flights/:id`
  - Body: partial fields same as create.
  - Responses: `200 { "flight": { ...updated } }`; `404` not found; `400` invalid body.
- `DELETE /flights/:id`
  - Responses: `200 { "message": "Flight deleted successfully" }`; `404` not found.

## Auth & Ownership Rules

- All `/hotels` and `/flights` endpoints require `Authorization: Bearer <token>`.
- Data is scoped to the authenticated user (`user_id`); accessing others’ data returns `404`.

## Error Response Shape

```json
{ "error": "<message>" }
```

---

# Frontend / Backend Testing Summary

- Frontend: Cypress (`cypress/e2e/smoke.cy.js`, `cypress/e2e/flight.cy.js`); Unit tests (`src/components/PrivateRoute.test.jsx`, `src/pages/HotelList.test.jsx`, `src/pages/HotelForm.test.jsx`, `src/pages/FlightList.test.jsx`, `src/pages/FlightForm.test.jsx`).
- Backend: `routes/router_test.go` (auth, hotels, flights, JWT, isolation, bad input).
