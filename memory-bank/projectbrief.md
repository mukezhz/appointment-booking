# Project Brief: Minimal Appointment Booking System

## Overview
A lightweight backend API built using Gin (Go) for appointment scheduling functionality, focusing on simplicity, LLM integration potential, and safe concurrency handling.

## Core Requirements

### Authentication System
- User registration and login functionality
- JWT-based authentication
- Simple role-based access (user/guest)

### Availability Management
- Weekly recurring availability slots
- Time slot management (30-minute intervals)
- Timezone handling

### Booking System
- Public booking interface for guests
- Concurrent booking handling
- Email notifications (future enhancement)

### User Dashboard
- View all bookings
- Manage availability
- Simple analytics (future enhancement)

## Technical Constraints
- Framework: Gin (Go)
- Database: PostgreSQL
- Authentication: JWT
- API Documentation: OpenAPI/Swagger
- Testing: Unit tests and integration tests

## Project Goals
1. Create a robust, concurrent-safe booking system
2. Implement clean architecture principles
3. Ensure code maintainability and scalability
4. Provide comprehensive API documentation
5. Enable future LLM integration capabilities
