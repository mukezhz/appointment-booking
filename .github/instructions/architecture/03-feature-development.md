# Feature Development Guide

This guide walks you through the practical steps of adding new features to the application while following our clean architecture principles.

## Development Workflow

### 1. Planning Phase
1. Define the feature requirements
2. Design the domain model and its relationships
3. Plan the API endpoints and request/response structures
4. Identify required business rules and validations

### 2. Implementation Phase

Follow our layered architecture pattern, implementing each component in order:

1. **Domain Model** - Define your core business entities
2. **Repository** - Implement data access
3. **Service** - Add business logic
4. **Controller** - Handle HTTP requests
5. **Routes** - Configure API endpoints
6. **Module** - Set up dependency injection

See our [Project Overview](01-project-overview.md) for architectural principles and [API Overview](../api/01-api-overview.md) for API design guidelines.

### Implementation Example

Here's a complete example of implementing a new feature:

#### A. Domain Model
```go
// domain/models/feature.go
type Feature struct {
    gorm.Model
    UUID      types.BinaryUUID `json:"uuid" gorm:"type:binary(16);uniqueIndex"`
    Name      string           `json:"name" gorm:"type:varchar(255)"`
    Status    FeatureStatus    `json:"status" gorm:"type:varchar(20)"`
}

// Add domain methods here
func (f *Feature) Validate() error {
    // Implement validation logic
}
```

#### B. Feature Package Structure
```
domain/
└── feature/
    ├── controller.go   # HTTP handlers
    ├── dto.go         # Request/response structures
    ├── errorz.go      # Domain-specific errors
    ├── module.go      # DI configuration
    ├── repository.go  # Database operations
    ├── route.go       # Route definitions
    └── service.go     # Business logic
```

#### C. Implementation Steps

1. **Define DTOs (dto.go)**
```go
type CreateFeatureRequest struct {
    Name   string        `json:"name" binding:"required"`
    Status FeatureStatus `json:"status" binding:"required,oneof=active inactive"`
}

type FeatureResponse struct {
    UUID      string        `json:"uuid"`
    Name      string        `json:"name"`
    Status    FeatureStatus `json:"status"`
    CreatedAt time.Time     `json:"created_at"`
}
```

2. **Implement Repository Layer**
```go
type Repository interface {
    Create(ctx context.Context, feature *models.Feature) error
    FindByUUID(ctx context.Context, uuid types.BinaryUUID) (*models.Feature, error)
    // Add other methods
}
```

3. **Implement Service Layer**
```go
type Service interface {
    CreateFeature(ctx context.Context, req *CreateFeatureRequest) (*FeatureResponse, error)
}

func (s *service) CreateFeature(ctx context.Context, req *CreateFeatureRequest) (*FeatureResponse, error) {
    // Implement business logic here
}
```

4. **Setup Module and DI**
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

### 3. Testing Phase

1. **Write Unit Tests**
   - Test each layer independently
   - Use mocks for dependencies
   - Cover error cases

2. **Write Integration Tests**
   - Test the complete feature flow
   - Use test containers for database tests
   - Test API endpoints

3. **Write Documentation**
   - Update API documentation
   - Add examples
   - Document any special considerations

## Validation Checklist

- [ ] Domain model properly structured
- [ ] All layers implemented (Repository, Service, Controller)
- [ ] Error handling follows project standards
- [ ] Input validation implemented
- [ ] Tests written and passing
- [ ] Documentation updated
- [ ] Code follows project style guide
- [ ] Proper logging implemented
- [ ] Metrics added (if required)

```go
type Repository struct {
    db     infrastructure.Database
    logger framework.Logger
}

func (r *Repository) Create(ctx context.Context, feature *models.Feature) error {
    if err := r.db.WithContext(ctx).Create(feature).Error; err != nil {
        return common.HandleDBError(err, FeatureErrMap)
    }
    return nil
}

func (r *Repository) GetByID(ctx context.Context, id types.BinaryUUID) (*models.Feature, error) {
    var feature models.Feature
    if err := r.db.WithContext(ctx).First(&feature, "uuid = ?", id).Error; err != nil {
        return nil, common.HandleDBError(err, FeatureErrMap)
    }
    return &feature, nil
}
```

### 4. Implement Service

```go
type Service struct {
    repo   *Repository
    logger framework.Logger
}

func (s *Service) CreateFeature(ctx context.Context, feature *models.Feature) (*models.Feature, error) {
    if err := ValidateFeature(feature); err != nil {
        return nil, err
    }
    if err := s.repo.Create(ctx, feature); err != nil {
        return nil, err
    }
    return feature, nil
}

func (s *Service) GetFeature(ctx context.Context, id types.BinaryUUID) (*models.Feature, error) {
    return s.repo.GetByID(ctx, id)
}
```

### 5. Implement Controller

```go
type Controller struct {
    service *Service
    logger  framework.Logger
}

func (c *Controller) CreateFeature(ctx *gin.Context) {
    var req CreateFeatureRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        responses.HandleValidationError(ctx, c.logger, err)
        return
    }

    feature, err := c.service.CreateFeature(ctx.Request.Context(), req.ToModel())
    if err != nil {
        responses.HandleError(ctx, c.logger, err)
        return
    }

    responses.DetailResponse(ctx, http.StatusCreated, responses.DetailResponseType[FeatureDTO]{
        Item:    c.toDTO(feature),
        Message: "Feature created successfully",
    })
}
```

### 6. Define Routes

In `route.go`:
```go
type Route struct {
    controller *Controller
}

func (r *Route) Register(router *gin.RouterGroup) {
    group := router.Group("/features")
    {
        group.POST("", r.controller.CreateFeature)
        group.GET("/:id", r.controller.GetFeature)
        // Add other routes
    }
}
```

### 7. Setup Dependency Injection

In `module.go`:
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

### 8. Add Tests

Create corresponding test files:
```bash
domain/feature/
├── service_test.go
├── repository_test.go
└── route_test.go
```

## Implementation Checklist

### Models
- [ ] Use GORM tags for database mapping
- [ ] Add validation methods
- [ ] Include domain logic methods
- [ ] Add proper JSON tags
- [ ] Define types and constants

### Repository
- [ ] Error mapping
- [ ] Transaction support
- [ ] Logging
- [ ] Context usage
- [ ] Test queries

### Service
- [ ] Validation logic
- [ ] Business rules
- [ ] Error handling
- [ ] Logging
- [ ] Transaction coordination

### Controller
- [ ] Input validation
- [ ] DTO conversions
- [ ] Error handling
- [ ] Response formatting
- [ ] Swagger docs

### Tests
- [ ] Unit tests for service
- [ ] Integration tests for repository
- [ ] API tests for routes
- [ ] Mock dependencies
- [ ] Test edge cases

## Migration Management

1. **Create Migration**
   ```bash
   make migrate-create name=create_features
   ```

2. **Edit Migration File** (`migrations/YYYYMMDDHHMMSS_create_features.sql`):
   ```sql
   -- migrate:up
   CREATE TABLE features (
       id BIGINT AUTO_INCREMENT PRIMARY KEY,
       uuid BINARY(16) NOT NULL,
       name VARCHAR(255) NOT NULL,
       status VARCHAR(20) NOT NULL,
       created_at DATETIME NOT NULL,
       updated_at DATETIME NOT NULL,
       UNIQUE KEY uk_features_uuid (uuid)
   );

   -- migrate:down
   DROP TABLE IF EXISTS features;
   ```

3. **Run Migration**
   ```bash
   make migrate-up
   ```

## Feature Documentation

Always include:

1. API Documentation
   - Swagger annotations
   - Request/response examples
   - Error scenarios

2. Domain Documentation
   - Business rules
   - Validation rules
   - Entity relationships

3. Integration Documentation
   - Dependencies
   - External services
   - Configuration
