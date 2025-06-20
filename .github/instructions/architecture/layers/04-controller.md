# Controller Layer Guide

The controller layer handles HTTP interactions. This guide explains how to implement and modify controllers in our clean architecture.

## Key Responsibilities

- Parse HTTP requests
- Convert between DTOs and domain models
- Handle HTTP responses
- Manage input validation
- Handle HTTP-specific errors

## Controller Structure

### Data Transfer Objects (DTOs)

```go
// Request DTOs
type CreateEntityRequest struct {
    Name      string    `json:"name" binding:"required"`
    StartDate time.Time `json:"start_date" binding:"required"`
    Status    string    `json:"status" binding:"required,oneof=active inactive"`
}

// Response DTOs
type EntityDTO struct {
    UUID      string    `json:"uuid"`
    Name      string    `json:"name"`
    StartDate time.Time `json:"start_date"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Controller Definition

```go
type Controller struct {
    service *Service           // Domain service
    logger  framework.Logger   // Structured logger
}
```

## DTO Conversion Methods

Always implement methods to convert between DTOs and domain models:

```go
func (c *Controller) toEntityDTO(model *models.Entity) *EntityDTO {
    if model == nil {
        return nil
    }
    return &EntityDTO{
        UUID:      model.UUID.String(),
        Name:      model.Name,
        StartDate: model.StartDate,
        Status:    string(model.Status),
        CreatedAt: model.CreatedAt,
    }
}
```

## Controller Method Pattern

Controller methods follow this standard pattern:

1. Parse and validate request
2. Convert DTO to model
3. Call service method
4. Handle errors
5. Convert model to DTO
6. Send response

Example implementation:

```go
// @Summary Create a new entity
// @Description Create a new entity with the provided data
// @Tags entities
// @Accept json
// @Produce json
// @Param request body CreateEntityRequest true "Entity details"
// @Success 201 {object} responses.DetailResponseType[EntityDTO]
// @Failure 400 {object} responses.DetailResponseType
// @Router /api/v1/entities [post]
func (c *Controller) CreateEntity(ctx *gin.Context) {
    // 1. Parse and validate request
    var req CreateEntityRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        responses.HandleValidationError(ctx, c.logger, err)
        return
    }

    // 2. Get authenticated user ID from context
    userID := ctx.MustGet(framework.UserIDKey).(types.BinaryUUID)

    // 3. Call service with domain model
    entity, err := c.service.CreateEntity(
        ctx.Request.Context(),
        userID,
        &models.Entity{
            Name:      req.Name,
            StartDate: req.StartDate,
            Status:    models.EntityStatus(req.Status),
        },
    )

    // 4. Handle errors
    if err != nil {
        responses.HandleError(ctx, c.logger, err)
        return
    }

    // 5. Convert to DTO and send response
    responses.DetailResponse(ctx, http.StatusCreated, responses.DetailResponseType[EntityDTO]{
        Item:    utils.SafeDeref(c.toEntityDTO(entity)),
        Message: "Entity created successfully",
    })
}
```

## Implementation Guidelines

When implementing a controller:

1. Define DTOs with validation tags
2. Add Swagger documentation
3. Implement DTO conversion methods
4. Use standard response types
5. Handle all possible errors
6. Log at appropriate levels

## Response Handling

Always use standard response types:

1. `responses.DetailResponse` for single item responses
2. `responses.ListResponse` for paginated list responses
3. `responses.HandleError` for error responses
4. `responses.HandleValidationError` for validation errors

Example:

```go
// Success response with single item
responses.DetailResponse(ctx, http.StatusOK, responses.DetailResponseType[EntityDTO]{
    Item:    dto,
    Message: "Entity retrieved successfully",
})

// List response with pagination
responses.ListResponse(ctx, http.StatusOK, responses.ListResponseType[EntityDTO]{
    Items:      items,
    TotalPages: totalPages,
    Count:      count,
    Message:    "Entities retrieved successfully",
})

// Error response
responses.HandleError(ctx, c.logger, err)

// Validation error
responses.HandleValidationError(ctx, c.logger, err)
```

## Testing Controllers

Controllers should have integration tests that:

1. Test HTTP request parsing
2. Test validation error handling
3. Test service layer integration
4. Test response formatting
5. Test authentication/authorization

Example test structure:
```go
func TestController_CreateEntity(t *testing.T) {
    tests := []struct {
        name       string
        input      CreateEntityRequest
        wantStatus int
        wantErr    bool
    }{
        // Define test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup test server
            // Make request
            // Assert response
        })
    }
}
```
