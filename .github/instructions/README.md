# Go Clean Architecture Documentation

## Overview

This is a Go-based appointment booking system built with clean architecture principles. The project uses Gin Framework, GORM, and follows domain-driven design patterns.

### Key Features
- Clean Architecture with Gin Web Framework
- Domain-Driven Design
- GORM for database operations
- Atlas for migrations
- Ginkgo for testing
- JWT authentication
- AWS service integration

## Project Structure

```
├── bootstrap/      # Application initialization
├── console/       # CLI commands
├── domain/        # Business logic & models
│   ├── models/    # Database models
│   └── <feature>/ # Feature modules
├── pkg/           # Shared packages
└── migrations/    # Database migrations
```

### Feature Module Structure
Each feature (e.g., appointment, user) follows:
```
domain/<feature>/
├── controller.go   # HTTP handlers
├── dto.go         # Data transfer objects
├── errorz.go      # Domain errors
├── module.go      # DI setup
├── repository.go  # Data access
├── route.go       # API routes
└── service.go     # Business logic
```

## Development Guide

### 1. Adding New Features

1. Define Domain Model in `domain/models/`
2. Create feature package in `domain/<feature>/`
3. Implement Repository, Service, Controller
4. Add routes and wire up DI
5. Write tests
6. Add API documentation

### 2. API Development

For detailed API implementation instructions, see the [API Implementation Guide](./api-guide.md). The guide covers:

- Complete implementation workflow
- Request/Response formats
- Error handling patterns
- Authentication & Security
- Testing strategies
- Performance considerations

Here's a quick reference for common tasks:

1. **Domain Model** (`models/feature.go`)
```go
type Feature struct {
    gorm.Model
    UUID      types.BinaryUUID `json:"uuid" gorm:"type:binary(16);uniqueIndex"`
    // Add model fields
}

func (f *Feature) Validate() error {
    // Add validation logic
}
```

2. **DTOs** (`domain/feature/dto.go`)
```go
type CreateFeatureRequest struct {
    // Request fields with validation tags
    Name string `json:"name" binding:"required"`
}

type FeatureResponse struct {
    // Response fields
    UUID string `json:"uuid"`
    Name string `json:"name"`
}
```

3. **Repository** (`domain/feature/repository.go`)
```go
type Repository struct {
    db     infrastructure.Database
    logger framework.Logger
}

func (r *Repository) Create(ctx context.Context, feature *models.Feature) error {
    // Implementation
}
```

4. **Service** (`domain/feature/service.go`)
```go
type Service struct {
    repo   *Repository
    logger framework.Logger
}

func (s *Service) Create(ctx context.Context, req *CreateFeatureRequest) (*FeatureResponse, error) {
    // Business logic implementation
}
```

5. **Controller** (`domain/feature/controller.go`)
```go
func (c *Controller) Create(ctx *gin.Context) {
    var req CreateFeatureRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        responses.HandleError(ctx, c.logger, err)
        return
    }
    
    result, err := c.service.Create(ctx.Request.Context(), &req)
    if err != nil {
        responses.HandleError(ctx, c.logger, err)
        return
    }
    
    responses.DetailResponse(ctx, http.StatusCreated, responses.DetailResponseType[*FeatureResponse]{
        Item:    result,
        Message: "Created successfully",
    })
}
```

6. **Routes** (`domain/feature/route.go`)
```go
func RegisterRoutes(r *Route) {
    api := r.handler.Group("/api/features")
    
    // Public routes
    api.GET("", r.controller.List)
    
    // Protected routes
    protected := api.Use(r.authMiddleware.HandleAuthWithRole())
    protected.POST("", r.controller.Create)
}
```

7. **Dependency Injection** (`domain/feature/module.go`)
```go
var Module = fx.Module("feature",
    fx.Provide(
        NewRepository,
        NewService,
        NewController,
        NewRoute,
    ),
    fx.Invoke(RegisterRoutes),
)
```

#### C. Best Practices

1. **Error Handling**
   - Define domain-specific errors in `errorz.go`
   - Use HTTP status codes appropriately
   - Return detailed error messages in development

2. **Validation**
   - Use gin binding tags for request validation
   - Implement domain model validation
   - Handle validation errors consistently

3. **Security**
   - Use middleware for authentication
   - Validate user permissions
   - Sanitize input data

4. **Documentation**
   - Use clear function and variable names
   - Add comments for complex logic
   - Document API endpoints

### 3. Testing Strategy

- Unit tests for business logic
- Integration tests for APIs
- Ginkgo for BDD-style tests
- TestContainers for database tests

### 4. Database Migrations

Using Atlas for schema management:
```bash
# Create migration
atlas migrate diff migration_name

# Apply migrations
atlas migrate apply
```

## Best Practices

1. **Clean Architecture**
   - Separate concerns
   - Dependencies point inward
   - Domain-driven design

2. **Code Organization**
   - Feature-based modules
   - Clear package boundaries
   - Consistent naming

3. **Error Handling**
   - Domain-specific errors
   - Consistent error responses
   - Proper status codes

4. **Testing**
   - Test business logic
   - Mock external dependencies
   - Integration test coverage

## Getting Started

1. Clone the repository
2. Copy `.env.example` to `.env`
3. Configure database settings
4. Run migrations
5. Start the server:
   ```bash
   go run main.go app:serve
   ```

For detailed setup instructions, see [README.md](../README.md).
