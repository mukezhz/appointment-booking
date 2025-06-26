# API Implementation Guide

This guide explains how to implement APIs in our clean architecture pattern.

## 1. File Structure

For each feature (e.g., appointments, users), create these files under `domain/<feature>/`:

```
domain/<feature>/
├── controller.go   # HTTP handlers
|-- constants.go<optional> # Add constants related to feature
├── dto.go         # Request/Response objects
├── errorz.go      # Domain-specific errors
├── module.go      # Dependency injection
├── repository.go  # Database operations
├── route.go       # API routes
└── service.go     # Business logic
```

## 2. Implementation Steps

### 2.1. Define Domain Model (`domain/models/<feature>.go`)

```go
type Feature struct {
    gorm.Model
    UUID   types.BinaryUUID `json:"uuid" gorm:"type:binary(16);uniqueIndex"`
    Name   string          `json:"name" gorm:"type:varchar(255)"`
    Status string          `json:"status" gorm:"type:varchar(20)"`
}

func (Feature) TableName() string {
    return "features"
}

func (u *Feature) BeforeCreate(tx *gorm.DB) error {
	if u.UUID.String() == (types.BinaryUUID{}).String() {
		id, err := uuid.NewRandom()
		u.UUID = types.BinaryUUID(id)
		return err
	}
	return nil
}

func (f *Feature) Validate() error {
    // Add validation logic
    return nil
}
```

### 2.2. Create DTOs (`domain/<feature>/dto.go`)

```go
// Request DTOs
type CreateFeatureRequest struct {
    Name   string `json:"name" binding:"required"`
    Status string `json:"status" binding:"required,oneof=active inactive"`
}

// Response DTOs
type FeatureResponse struct {
    UUID      string    `json:"uuid"`
    Name      string    `json:"name"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"createdAt"`
}
```

### 2.3. Implement Repository (`domain/<feature>/repository.go`)

```go
type Repository struct {
	db     infrastructure.Database
	logger framework.Logger
}

func NewRepository(
	db infrastructure.Database,
	logger framework.Logger,
) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) Create(ctx context.Context, feature *models.Feature) error {
    return r.db.Model(&models.Feature{}).Create(feature).Error
}
```

### 2.4. Implement Service (`domain/<feature>/service.go`)

```go
type Service struct {
    repo   *Repository
    logger framework.Logger
}

func NewService(
    repo *Repository, 
    logger framework.Logger,
) *Service {
    return &Service{
        repo: repo, 
        logger: logger,
    }
}

func (s *Service) Create(ctx context.Context, req *CreateFeatureRequest) (*FeatureResponse, error) {
    // Business logic implementation
    feature := &models.Feature{
        UUID:   types.NewBinaryUUID(),
        Name:   req.Name,
        Status: req.Status,
    }
    
    if err := s.repo.Create(ctx, feature); err != nil {
        return nil, err
    }
    
    return &FeatureResponse{
        UUID:      feature.UUID.String(),
        Name:      feature.Name,
        Status:    feature.Status,
        CreatedAt: feature.CreatedAt,
    }, nil
}
```

### 2.5. Implement Controller (`domain/<feature>/controller.go`)

```go
type Controller struct {
    service *Service
    logger  framework.Logger
}

func NewController(
    service *Service, 
    logger framework.Logger,
) *Controller {
    return &Controller{
        service: service, 
        logger: logger,
    }
}

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
        Message: "Feature created successfully",
    })
}
```

### 2.6. Define Routes (`domain/<feature>/route.go`)

```go
type Route struct {
    logger     framework.Logger
    handler    infrastructure.Router
    controller *Controller
}

func NewRoute(
    logger framework.Logger,
    handler infrastructure.Router,
    controller *Controller,
) *Route {
    return &Route{
        handler:    handler,
        logger:     logger,
        controller: controller,
    }
}

func RegisterRoutes(r *Route) {
    api := r.handler.Group("/api/features")
    
    // Public routes
    api.GET("", r.controller.List)
    
    // Protected routes with authentication
    protected := api.Use(r.authMiddleware.HandleAuthWithRole())
    {
        protected.POST("", r.controller.Create)
        protected.PUT("/:id", r.controller.Update)
        protected.DELETE("/:id", r.controller.Delete)
    }
}
```

### 2.7. Setup Dependency Injection (`domain/<feature>/module.go`)

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

## 3. Standard Response Format

### Success Response
```json
{
    "item": {},     // single item
    "items": [],    // list of items
    "message": "",  // success message
    "meta": {       // optional metadata for pagination
        "page": 1,
        "limit": 10,
        "total": 100
    }
}
```

### Error Response
```json
{
    "error": "Error message",
    "details": {
        "field": "error detail"
    }
}
```

## 4. Best Practices

### 4.1. Error Handling
- Define domain-specific errors in `errorz.go`
- Use appropriate HTTP status codes
- Include helpful error messages
- Log errors with context

### 4.2. Input Validation
- Use gin binding tags for request validation
- Implement domain model validation
- Return clear validation error messages

### 4.3. Security
- Use authentication middleware for protected routes
- Validate user permissions
- Sanitize input data
- Use HTTPS in production

### 4.4. Performance
- Use pagination for list endpoints
- Cache frequently accessed data
- Use database indexes appropriately
- Handle concurrent requests safely

### 4.5. Documentation
- Use clear naming for routes and handlers
- Document request/response formats
- Include examples in comments
- Keep API versioning consistent

## 5. Common Patterns

### 5.1. Pagination
```go
type PaginationQuery struct {
    Page  int `form:"page" binding:"required,min=1"`
    Limit int `form:"limit" binding:"required,min=1,max=100"`
}

func (r *Repository) List(ctx context.Context, p *PaginationQuery) ([]models.Feature, int64, error) {
    var total int64
    var features []models.Feature
    
    offset := (p.Page - 1) * p.Limit
    
    tx := r.db.Model(&models.Feature{})
    tx.Count(&total)
    
    if err := tx.Offset(offset).Limit(p.Limit).Find(&features).Error; err != nil {
        return nil, 0, err
    }
    
    return features, total, nil
}
```

### 5.2. Search and Filtering
```go
type ListFeatureQuery struct {
    PaginationQuery
    Status string `form:"status"`
    Search string `form:"search"`
}

func (r *Repository) List(ctx context.Context, q *ListFeatureQuery) ([]models.Feature, int64, error) {
    tx := r.db.Model(&models.Feature{})
    
    if q.Status != "" {
        tx = tx.Where("status = ?", q.Status)
    }
    
    if q.Search != "" {
        search := "%" + q.Search + "%"
        tx = tx.Where("name LIKE ?", search)
    }
    
    // ... pagination logic
}
```

### 5.3. Batch Operations
```go
func (r *Repository) BatchCreate(ctx context.Context, features []*models.Feature) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        for _, feature := range features {
            if err := tx.Create(feature).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

## 6. API Documentation with Bruno

The project uses Bruno for API documentation and testing. API documentation is located in the `docs/` directory.

### 6.1. Appointments API Endpoints

#### Availability Management
- `POST /api/appointments/availability` - Create availability slot
- `GET /api/appointments/availability` - List availability slots

#### Booking Management
- `POST /api/public/appointments/book` - Create booking (public)
- `GET /api/appointments/bookings` - List bookings
- `PATCH /api/appointments/bookings/:id/status` - Update booking status

See the Bruno files in `docs/appointments/` for detailed request/response examples and testing.

### 6.2. Bruno File Structure
```
docs/
└── appointments/
    ├── folder.bru              # Module description
    ├── CreateAvailability.bru  # Create availability slot
    ├── GetAvailabilities.bru   # List availability slots
    ├── CreateBooking.bru       # Create booking
    ├── GetBookings.bru         # List bookings
    └── UpdateBookingStatus.bru # Update booking status
```

### 6.3. Example Bruno Test
```bruno
meta {
  name: CreateBooking
  type: http
}

post {
  url: {{baseURL}}/api/public/appointments/book
  body: json
}

body:json {
  {
    "userId": 1,
    "date": "2024-06-24",
    "startTime": "10:00",
    "endTime": "11:00",
    "timeZone": "UTC",
    "clientName": "John Doe",
    "clientEmail": "john@example.com"
  }
}
```
