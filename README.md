# CEN5035-Project
# TravelLog — Hotel & Flight Log Web App

A simple full-stack web application for logging hotels and flights, showing travel routes and hotel nights on a map, and generating month / quarter / year summaries.  
This project is designed for a **4-week software engineering class project** and a **team of four**.

---
## Members


- Ao Wang : Front-end
- Ray Chen : Front-end
- Wentao Chen : Back-end
- Jianan Gu : Back-end
## 1. Project Goals

- Practice frontend + backend collaboration
- Build a realistic CRUD web app with visualization
- Keep scope **small, clear, and finishable in 4 weeks**

---

## 2. Core Features

### 1) Hotel Log
- Add / edit / delete hotel stays
- Fields:
  - Hotel name
  - City, country
  - Check-in date, check-out date
  - Total price
  - Notes (optional)
- List hotels in a table

### 2) Flight Log
- Add / edit / delete flights
- Fields:
  - Airline, flight number
  - From airport (IATA), to airport (IATA)
  - Departure time, arrival time
  - Total price
  - Notes (optional)
- List flights in a table

### 3) Map View
- Show flight routes as lines on a map
- Show hotel stays as markers
- Clicking a marker or route shows basic details

### 4) Time Summary
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
travellog/
README.md
docker-compose.yml
frontend/
src/
package.json
backend/
src/
prisma/
package.json


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
- Hotel log UI (table + form)

Frontend Developer B
- Flight log UI
- Map page and summary page

Backend Developer A
- Auth (login / register)
- Hotel CRUD APIs

Backend Developer B
- Flight CRUD APIs
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

## 11. License

MIT

