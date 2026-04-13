# Sprint 3

## This document contains the User Stories and Issues specific to Sprint 3 only. It does not represent the complete project scope.

# User Stories

## 1. View Map with Flight Routes and Hotel Markers

As a logged-in user, I want to see my flight routes as lines and my hotel stays as markers on a map, so that I can visualize my travel history geographically.

## 2. View Travel Summary

As a logged-in user, I want to see a summary of my travel data by month, quarter, or year, so that I can understand my travel patterns and spending.

## 3. Interactive Map Details

As a logged-in user, I want to click on a marker or route on the map to see basic details, so that I can quickly review my travel information.

## 4. Comprehensive Testing

As a developer, I want to maintain test coverage for new functionality, so that the application remains reliable.

---

# Issues Planned for Sprint 3

## Frontend Issues

### Ao Wang (GitHub: Ao Wang, William-A-Wang)

- [Frontend] Map page layout with Leaflet/Mapbox integration.
- [Frontend] Summary page layout with monthly/quarterly/yearly toggle.
- [Frontend] Unit tests for map and summary UI components.

### Ray Chen (GitHub: rkc8626, kzc-qubit)

- [Frontend] Map markers for hotels (city coordinates) and polylines for flight routes (airport coordinates).
- [Frontend] Map interaction: click markers/routes to show details modal/tooltips.
- [Frontend] Cypress E2E covering map load + summary view (smoke).
- [Frontend] Install leaflet and react-leaflet dependencies
- [Frontend] Create airport IATA code with coordinates lookup table
- [Frontend] Create geocoding utility for hotel city/country
- [Frontend] Create TravelMap page component with markers and polylines
- [Frontend] Add /map route to App.jsx and NavBar link
- [Frontend] Add map-related CSS styles

## Backend Issues

### Wentao Chen (GitHub: wentao-chen227)

- [Backend] Backend tests for map/summary handlers (happy path + auth/ownership).
- [Backend] Support/auth review for new endpoints (ensure JWT + routing consistent).

### Jianan Gu (GitHub: jianangu1214, Griffin)

- [Backend] Implement `GET /map` (hotels with lat/lng, flights with from/to lat/lng).
- [Backend] Implement `GET /summary?type=month|quarter|year` returning flights count, hotels count, total nights, total spending.
- [Backend] Update API documentation in this file for Map/Summary.

---

# Sprint 3 Completion Status

## Issues the team planned to address

- **Frontend (Ao Wang):** Map page layout, Summary page layout, unit tests.
- **Frontend (Ray Chen):** Hotel markers, flight route lines, interactive map details.
- **Backend (Wentao Chen):** Map API implementation, unit tests.
- **Backend (Jianan Gu):** Summary API implementation, API documentation.

## Successfully completed

_Update this section after Sprint 3 implementation_

## Not completed (and why)

_Update this section if any issues remain incomplete_

---

# Backend API Documentation

Map/Summary are Sprint 3 additions; Auth/Hotels/Flights are retained from Sprint 2 for completeness.\_

## Map (Sprint 3)

- `GET /map` (JWT required)
- Response: `200 { "hotels": [ { id, hotel_name, city, country, latitude, longitude } ], "flights": [ { id, airline, flight_number, from_airport, to_airport, from_lat, from_lng, to_lat, to_lng, depart_time, arrive_time, price, notes } ] }`
- Notes: common airports have built-in lat/lng; unknown airports return `0,0` (frontend should geocode fallback); data scoped to user.

## Summary (Sprint 3)

- `GET /summary?type=month|quarter|year` (default `month`, JWT required)
- Response: `200 { "type": "<type>", "summaries": [ { period, flights, hotels, nights, spend }, ... ] }`
- Notes: periods sorted ascending; nights = (check_out - check_in) days, min 0; spend = flights + hotels price per period.

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

- All `/auth`, `/hotels`, `/flights`, `/map`, `/summary` endpoints require proper Bearer token where noted.
- Data is scoped to the authenticated user (`user_id`); accessing others’ data yields 404.

## Error Response Shape

```json
{ "error": "<message>" }
```
