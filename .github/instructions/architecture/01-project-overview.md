# Project Overview

## Introduction

This project implements a modern web application using Go and Clean Architecture principles. The codebase is organized to be maintainable, testable, and scalable while following best practices in software design.

## Key Principles

1. **Clean Architecture**
   - Clear separation of concerns
   - Independent of frameworks and external agencies
   - Highly testable design
   - Independent of UI, database, or external services

2. **Domain-Driven Design**
   - Business logic centered around domain models
   - Rich domain models with behavior
   - Clear boundaries between domains
   - Ubiquitous language in code

3. **Dependency Injection**
   - Uses uber-go/fx for DI
   - Loose coupling between components
   - Clear dependency graphs
   - Easy testing and mocking

4. **Standard Patterns**
   - Consistent code organization
   - Standard error handling
   - Common response formats
   - Reusable components

## Key Technologies

1. **Core Framework**
   - Go (1.23+)
   - Gin Web Framework
   - GORM for database access
   - uber-go/fx for dependency injection

2. **Testing**
   - Ginkgo for BDD-style testing
   - TestContainers for integration tests
   - Mocking utilities

3. **Infrastructure**
   - MySQL database
   - Redis for caching (optional)
   - AWS services integration

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
