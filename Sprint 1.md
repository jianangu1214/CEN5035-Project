# Sprint 1

This document contains the User Stories and Issues specific to Sprint 1 only. It does not represent the complete project scope. User Stories and Issues for other features (Flight Log, Map View, Summary) will be added in subsequent sprints.
---

# User Stories

## 1. Registration

As a new user.
I want to create an account with my email and password.
So that I can save my travel logs and access them securely.

## 2. Login

As a registered user.
I want to log in to my account.
So that I can access my personal travel.

## 3. Logout

As a logged-in user.
I want to log out of my account.
So that I can secure my data when I'm done.

## 4. Add Hotel Stay

As a logged-in user.
I want to add a new hotel stay record.
So that I can keep track of my accommodations.

## 5. View Hotel List

As a logged-in user.
I want to view all my hotel stays in a list.
So that I can see my accommodation history.

## 6. Edit Hotel Stay

As a logged-in user.
I want to edit an existing hotel record.
So that I can correct or update information.

## 7. Delete Hotel Stay

As a logged-in user.
I want to delete a hotel record.
So that I can remove unwanted entries.

---

# Issues Planned for Sprint 1

## Frontend Issues

### Ao Wang:

[Setup] Initialize React frontend project.
[Frontend] Create registration page.
[Frontend] Create login page.
[Frontend] Implement logout and route protection.

### Ray Chen

[Frontend] Create hotel list page.
[Frontend] Create add/edit hotel form.
[Frontend] Implement hotel delete.

## Backend Issues

### Wentao Chen

[Setup] Initialize Go backend with Gin.
[Backend] Create User model.
[Backend] Implement POST /auth/register.
[Backend] Implement POST /auth/login.

### Jianan Gu

[Setup] Set up PostgreSQL with Docker.
[Backend] Create Hotel model.
[Backend] Implement Hotel CRUD APIs.

---

# Sprint 1 Completion Status

## Issues the team planned to address

- **Frontend (Ao Wang):** Initialize React frontend, registration page, login page, logout and route protection.
- **Frontend (Ray Chen):** Hotel list page, add/edit hotel form, hotel delete.
- **Backend (Wentao Chen):** Initialize Go backend with Gin, User model, POST /auth/register, POST /auth/login.
- **Backend (Jianan Gu):** PostgreSQL with Docker, Hotel model, Hotel CRUD APIs.

## Successfully completed

- [Frontend] Create hotel list page — table of hotel stays (hotel name, city, country, check-in, check-out, price, notes); empty state; links to add/edit.
- [Frontend] Create add/edit hotel form — single form for create and edit; fields: hotel name, city, country, check-in, check-out, price, notes.
- [Frontend] Implement hotel delete — delete button per row with confirmation; removes hotel and refreshes list.

*(Other items — fill in as the team completes them.)*

## Not completed (and why)

- *(List any issues that were not done and brief reason.)*
