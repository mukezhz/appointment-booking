# Project Overview

## Introduction

This project implements a modern web application using Go and Clean Architecture principles. The codebase is organized to be maintainable, testable, and scalable while following best practices in software design.

## Clean Architecture Layers

Our application follows a layered architecture pattern that enforces separation of concerns:

### 1. Domain Models Layer (`domain/models/`)
- Core business entities and logic
- Framework-independent
- Rich domain models with validation
- GORM tags for persistence
- Example: User, Booking, Availability models

### 2. Repository Layer (`domain/<feature>/repository.go`)
- Data access abstractions
- Database operations
- Framework-independent interfaces
- Implementation can use any data store
- Example: GORM implementations

### 3. Service Layer (`domain/<feature>/service.go`)
- Business logic implementation
- Uses repository interfaces
- Coordinates across repositories
- Implements use cases
- Transaction management

### 4. Controller Layer (`domain/<feature>/controller.go`)
- HTTP request handling
- Request validation
- Response formatting
- Error handling
- Uses service interfaces

## Key Principles

1. **Clean Architecture**
   - Clear separation of concerns
   - Dependencies point inward
   - Highly testable design
   - Framework independence

2. **Domain-Driven Design**
   - Business logic in domain models
   - Clear bounded contexts
   - Ubiquitous language
   - Rich domain behavior

3. **Dependency Injection**
   - Uses uber-go/fx
   - Loose coupling
   - Clear dependency graphs
   - Testable components

## Key Technologies

1. **Core Framework**
   - Go 1.21+
   - Gin Web Framework
   - GORM ORM
   - uber-go/fx

2. **Testing Tools**
   - Ginkgo (BDD-style)
   - TestContainers
   - Mocking utilities

3. **Infrastructure**
   - MySQL database
   - Redis (optional)
   - AWS services

## Project Goals

1. **Maintainability**
   - Clear, well-documented code
   - Consistent patterns
   - Easy to understand structure
   - Self-documenting architecture

2. **Testability**
   - Independent layers
   - Easy to mock dependencies
   - Comprehensive test coverage
   - Integration test support

3. **Scalability**
   - Modular design
   - Easy to add new features
   - Performance considerations
   - Cloud-native ready

4. **Security**
   - Built-in security patterns
   - Authentication/Authorization
   - Input validation
   - Error handling

## Next Steps

- [Clean Architecture Implementation](./02-clean-architecture.md)
- [Dependency Injection](./03-dependency-injection.md)
- [Project Structure](./04-project-structure.md)
