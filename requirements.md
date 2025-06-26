# ✅ Core Features to Build for Demo

Build the essential backend API functionality for an appointment booking system.

---

## 📆 Availability Management (Doctor)

- Endpoint to create available time slots per day.
- Allow setting start and end time.
- Prevent overlapping availability entries (optional for demo).

---

## 🕒 Appointment Booking (Patient)

- Endpoint to fetch available time slots for a doctor.
- Endpoint to book a time slot with a doctor.
- Ensure a time slot cannot be double-booked.
- Basic input validation and conflict check.

---

## 🔁 Reschedule or Cancel Appointment

- Endpoint to reschedule an existing appointment.
- Endpoint to cancel an appointment.
- Prevent rescheduling to already-booked slots.

---

## 📃 Basic User Entities

- Patient and Doctor users (use hardcoded/mock users for demo).
- No authentication required (can use user IDs directly).

---

## 🗂 Suggested API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/doctors/:id/availability` | Set availability |
| `GET`  | `/doctors/:id/slots` | Get available time slots |
| `POST` | `/appointments` | Book an appointment |
| `PUT`  | `/appointments/:id/reschedule` | Reschedule an appointment |
| `DELETE` | `/appointments/:id` | Cancel appointment |
| `GET`  | `/appointments/:userId` | View appointments for a user |

---

## 🔄 Conflict Handling

- Prevent double booking by checking existing appointments for the doctor.
- Use simple database-level unique constraints or logic in the booking logic.
