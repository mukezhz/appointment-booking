# API Overview Guide

## Overview

This guide demonstrates how to implement REST APIs in our Go Clean Architecture project, using the appointment booking system as an example. This guide ensures consistent API design and implementation across the project.

## Key Concepts

1. **Clean Architecture Separation**
   - Models for domain entities
   - DTOs for API requests/responses
   - Clear layer boundaries
   - Standard response formats

2. **RESTful Best Practices**
   - Resource-based routing
   - Proper HTTP methods
   - Consistent response structures
   - Proper error handling

3. **Documentation**
   - Swagger annotations
   - Clear API descriptions
   - Request/response examples
   - Error scenarios

## Implementation Flow

Our API implementation follows this flow:
1. Design the API endpoints
2. Define request/response DTOs
3. Implement the API handlers
4. Add validation and error handling
5. Document with Swagger

## API Design Best Practices

1. **URL Design**
   - Use plural nouns for resources: `/api/v1/users`
   - Nest related resources: `/api/v1/doctors/{id}/availability`
   - Use query parameters for filtering: `?status=active&sort=created_at`

2. **HTTP Methods**
   - GET: Retrieve resources
   - POST: Create new resources
   - PUT/PATCH: Update resources
   - DELETE: Remove resources

3. **Status Codes**
   - 200: Success (GET, PUT, PATCH)
   - 201: Created (POST)
   - 204: No Content (DELETE)
   - 400: Bad Request
   - 401: Unauthorized
   - 403: Forbidden
   - 404: Not Found
   - 422: Unprocessable Entity
   - 500: Internal Server Error

4. **Security**
   - Always validate input
   - Use proper authentication middleware
   - Implement rate limiting
   - Sanitize responses
   - Log security events

5. **Performance**
   - Use pagination for list endpoints
   - Implement caching where appropriate
   - Optimize database queries
   - Monitor response times

## Implementation Components

Each API feature consists of these key files:

1. **Models** (`domain/models/`)
   - Database entities with GORM annotations
   - Pure domain models without transport/persistence concerns
   - Example: `models.Availability`, `models.Appointment`

2. **DTOs** (`domain/<feature>/dto.go`)
   - Request/Response data structures
   - Input validation using binding tags
   - Clear separation from domain models

3. **Repository** (`domain/<feature>/repository.go`)
   - Database access layer
   - CRUD operations using GORM
   - Error mapping using `common.HandleDBError`
   - Contextual logging

4. **Service** (`domain/<feature>/service.go`)
   - Business logic
   - Model validation
   - Transaction handling
   - Works with domain models only

5. **Controller** (`domain/<feature>/controller.go`)
   - HTTP request/response handling
   - DTO to model conversion
   - Standard response formats
   - Error handling using `responses.HandleError`

6. **Routes** (`domain/<feature>/route.go`)
   - Endpoint registration
   - Middleware configuration
   - API documentation with Swagger annotations
