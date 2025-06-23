---
mode: 'agent'
description: 'Generate a project of your **Appointment Booking API** requirements for a **doctor-patient platform**'
---

# 📆 Minimal Appointment Booking System – Project Specification

A lightweight backend API built using **Gin (Go)** to demonstrate appointment scheduling functionality, with minimal features, LLM integration potential, and safe concurrency handling.

---

## ✅ Functional Requirements

### 1. User Authentication

Allow a single user to register and log in to manage their availability and view bookings.

#### Endpoints
- `POST /auth/register` – Register with email and password
- `POST /auth/login` – Login to receive JWT token

---

### 2. Set Weekly Availability

User defines their recurring weekly availability.

#### Fields
- `weekday`: string (`Monday`, `Tuesday`, etc.)
- `start_time`: string (`09:00`)
- `end_time`: string (`17:00`)

#### Endpoints
- `GET /availability` – View availability slots
- `POST /availability` – Create availability slot

---

### 3. Public Booking (Guest)

A guest can view the user’s available slots and book one by submitting their name, email, and preferred time.

#### Booking Fields
- `guest_name`: string
- `guest_email`: string
- `date`: string (`YYYY-MM-DD`)
- `time`: string (`HH:MM`)

#### Endpoints
- `GET /public/slots` – Get available slots for a date range
- `POST /public/book` – Book a slot

---

### 4. View Bookings (User)

The user can view all their appointments.

#### Endpoints
- `GET /bookings` – List of guest bookings

