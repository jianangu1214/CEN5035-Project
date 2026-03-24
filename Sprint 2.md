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

- (Fill after sprint progress)

- Implemented JWT-based authentication and verified protected routes (`/me`, `/flights`, `/hotels`).
- Configured backend to use PostgreSQL via Docker and ensured correct DB connectivity.
- Implemented backend integration tests covering:
  - user registration and login
  - JWT validation
  - protected route access
  - flight creation, retrieval, and deletion
  - hotel creation and retrieval
  - user-level data isolation
  - invalid input handlings

## Not completed (and why)

- (Fill after sprint progress)
