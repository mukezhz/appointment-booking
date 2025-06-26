# Technical Context

## Technology Stack
1. **Backend**
   - Language: Go
   - Framework: Gin
   - ORM: GORM
   - DI: fx

2. **Database**
   - MySQL
   - Migration tool: Atlas

3. **Testing**
   - Go testing package
   - API testing: Bruno

## Development Setup
1. **Prerequisites**
   - Go 1.x
   - MySQL
   - Docker & Docker Compose
   - Bruno (for API testing)

2. **Project Structure**
```
/
├── bootstrap/      # App initialization
├── domain/        # Business logic & features
├── pkg/           # Shared packages
├── migrations/    # Database migrations
└── docs/         # API documentation
```

3. **Key Libraries**
- gin: Web framework
- gorm: ORM
- fx: Dependency injection
- atlas: Database migrations

## Technical Constraints
1. **Database**
   - MySQL for persistence
   - Proper indexing
   - Transaction support

2. **API**
   - RESTful principles
   - JSON payloads
   - Proper status codes

3. **Performance**
   - Efficient queries
   - Proper error handling
   - Resource cleanup

## Development Practices
1. **Code Quality**
   - Clean architecture
   - SOLID principles
   - Proper documentation

2. **Testing**
   - Unit tests
   - Integration tests
   - API documentation

3. **Version Control**
   - Feature branches
   - Pull requests
   - Code review
