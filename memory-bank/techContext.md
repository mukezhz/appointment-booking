# Technical Context: Appointment Booking System

## Technology Stack

### Backend
- Language: Go 1.23+
- Framework: Gin
- Database: PostgreSQL
- ORM: GORM
- Authentication: JWT

### Development Tools
- Docker for containerization
- Make for build automation
- Swagger for API documentation
- Go modules for dependency management

## Project Structure
```
├── bootstrap/      # Application bootstrapping
├── domain/        # Business logic and models
│   ├── models/   
│   ├── services/
│   └── repositories/
├── api/           # HTTP handlers and routes
├── pkg/           # Shared packages
└── migrations/    # Database migrations
```

## Key Dependencies
1. gin-gonic/gin: Web framework
2. golang-jwt/jwt: JWT authentication
3. gorm.io/gorm: ORM
4. swaggo/swag: API documentation

## Development Setup
1. Go 1.23+ installation
2. PostgreSQL database
3. Docker and Docker Compose
4. Make utility

## Build Process
1. Generate API documentation
2. Run database migrations
3. Compile application
4. Run tests
5. Build Docker image

## Deployment Requirements
1. PostgreSQL database
2. Environment configuration
3. Docker runtime
4. SSL certificate (production)

## Security Measures
1. JWT authentication
2. Input validation
3. Rate limiting
4. SQL injection prevention
5. XSS protection

## Monitoring
1. Error logging
2. Request tracking
3. Performance metrics
4. Database monitoring
