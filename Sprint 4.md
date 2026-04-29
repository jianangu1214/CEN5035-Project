# Sprint 4

## This document contains the User Stories and Issues specific to Sprint 4 only. It does not represent the complete project scope.

# User Stories

## 1. Fix Issues Found During Sprint 3

As a developer, I want to fix any issues discovered during Sprint 3, so that the final application is stable and ready for submission.

## 2. Polish the Final User Experience

As a user, I want the application to feel complete and easy to use, so that I can smoothly navigate hotels, flights, map, and summary features.

## 3. Verify Final Test Coverage

As a developer, I want to run all existing tests and add tests for any new fixes or final adjustments, so that the completed project is reliable.

## 4. Prepare the Final Project Demo

As a team, we want to present the completed project clearly, so that we can demonstrate all frontend functionality, explain the backend API, and show final test results.

---

# Issues Planned for Sprint 4

## Frontend Issues

### Ao Wang (GitHub: Ao Wang, William-A-Wang)

- [Frontend] Refine Summary page behavior and final presentation.
- [Frontend] Polish final frontend layout, navigation flow, and styling across the application.
- [Frontend] Add or update unit tests for final frontend fixes and UI adjustments.

### Ray Chen (GitHub: rkc8626, kzc-qubit)

- [Frontend] Polish TravelMap behavior, marker/route interaction, and fallback handling.
- [Frontend] Review and fix frontend issues discovered during Sprint 3 integration.
- [Frontend] Add or update Cypress coverage for final integrated frontend flows.

## Backend Issues

### Wentao Chen (GitHub: wentao-chen227)

- [Backend] Re-run and verify all backend tests, including Sprint 3 `/map` and `/summary` tests.
- [Backend] Add tests for any backend fixes made during Sprint 4.
- [Backend] Review auth, routing, and response consistency across all final endpoints.

### Jianan Gu (GitHub: jianangu1214, Griffin)

- [Backend] Refine `/map` and `/summary` behavior based on final frontend integration and Sprint 3 findings.
- [Backend] Update backend API documentation in this file to reflect the final completed project.
- [Backend] Help finalize the front-page README with backend run/test/use instructions.

---

# Sprint 4 Completion Status

## Issues the team planned to address

- **Frontend (Ao Wang):** Summary page refinement, UI polish, frontend unit tests.
- **Frontend (Ray Chen):** TravelMap polish, frontend bug fixes, Cypress coverage.
- **Backend (Wentao Chen):** Backend regression testing, new backend tests for fixes, endpoint consistency review.
- **Backend (Jianan Gu):** Map/Summary refinement, final API documentation, README backend instructions.

## Successfully completed

- **Frontend (Ao Wang):** Summary page refinement, UI polish, frontend unit tests.
- **Frontend (Ray Chen):** TravelMap polish, frontend bug fixes, Cypress coverage.
- **Backend (Wentao Chen):** Backend regression testing, new backend tests for fixes, endpoint consistency review.
- **Backend (Jianan Gu):** Map/Summary refinement, final API documentation, README backend instructions.

All team member finish their own work

## Not completed (and why)

All finished

---

# Frontend Tests (Sprint 4)

## Existing frontend tests carried into Sprint 4

- Cypress E2E
  - `cypress/e2e/smoke.cy.js`: logs in, creates a hotel, and verifies it appears in the list.
  - `cypress/e2e/flight.cy.js`: logs in, opens Flights page, opens Add Flight form, and checks required fields.
- Unit tests (Vitest + Testing Library)
  - `src/components/PrivateRoute.test.jsx`
  - `src/pages/HotelList.test.jsx`
  - `src/pages/HotelForm.test.jsx`
  - `src/pages/FlightList.test.jsx`
  - `src/pages/FlightForm.test.jsx`

## Sprint 4 test additions / updates

- Add or update frontend tests for final map behavior.
- Add or update frontend tests for final summary-related fixes or UI changes.
- Re-run all frontend unit and Cypress tests for the final submission.

---

# Backend Unit Tests (Sprint 4)

## Existing backend tests carried into Sprint 4

- `routes/router_test.go`
  - Auth: register/login, JWT validation, `/me` access.
  - Flights: create, list (user isolation), delete, invalid time input, auth required.
  - Hotels: create, list (user isolation).
  - Map: success, unauthorized access, own-data isolation.
  - Summary: month aggregation, unauthorized access, invalid `type` fallback.

## Sprint 4 test additions / updates

- Add backend tests for any fixes made during Sprint 4.
- Re-run all backend tests from Sprints 2 and 3 for final verification.
- Confirm final backend behavior for auth, ownership, and response consistency.

## How to run backend tests

1. Start PostgreSQL:
   - `cd travellog-backend && docker-compose up -d`
2. Run tests:
   - `cd travellog-backend/backend`
   - `go test ./routes -v`
3. If cache/permission issues appear:
   - `GOMODCACHE=$PWD/.gocache GOCACHE=$PWD/.gocache go test ./routes -v`

---

# Backend API Documentation

## Map

- `GET /map` (JWT required)
- Response: `200 { "hotels": [ { id, hotel_name, city, country, latitude, longitude } ], "flights": [ { id, airline, flight_number, from_airport, to_airport, from_lat, from_lng, to_lat, to_lng, depart_time, arrive_time, price, notes } ] }`
- Notes: common airports have built-in lat/lng; unknown airports may return `0,0`, so the frontend should support fallback geocoding; data is scoped to the authenticated user.

## Summary

- `GET /summary?type=month|quarter|year` (default `month`, JWT required)
- Response: `200 { "type": "<type>", "summaries": [ { period, flights, hotels, nights, spend }, ... ] }`
- Notes: periods are sorted ascending; nights = `(check_out - check_in)` days with minimum `0`; spend = hotel price + flight price totals per period; the `type` query is normalized before validation, and unsupported values fall back to `month`.

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

- All `/auth`, `/hotels`, `/flights`, `/map`, and `/summary` endpoints require a proper Bearer token where noted.
- Data is scoped to the authenticated user (`user_id`); accessing another user’s data should not be allowed.

## Error Response Shape

```json
{ "error": "<message>" }
```
