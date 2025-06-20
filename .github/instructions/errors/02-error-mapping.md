# Error Mapping Guide

This guide explains how to map between different types of errors in our application.

## Error Maps

Each domain defines its error mappings in an `errorz.go` file to convert between different error types:

### 1. Database Error Mapping

```go
// domain/availability/errorz.go
var ErrMap = map[error]error{
    gorm.ErrRecordNotFound: &errorz.ResourceNotFoundError{
        Resource: "availability",
    },
    ErrorSlotAlreadyBooked: &errorz.DuplicateResourceError{
        Resource: "appointment",
        Field:    "slot",
    },
}

// Common error mapping function
func HandleDBError(err error, errMap map[error]error) error {
    if mappedErr, ok := errMap[err]; ok {
        return mappedErr
    }
    return err
}
```

### 2. HTTP Status Mapping

```go
// pkg/responses/error_handler.go
var statusMap = map[error]int{
    &errorz.ValidationError{}:      http.StatusBadRequest,
    &errorz.ResourceNotFoundError{}: http.StatusNotFound,
    &errorz.AuthorizationError{}:   http.StatusForbidden,
    &errorz.BusinessError{}:        http.StatusUnprocessableEntity,
}

func getStatusCode(err error) int {
    for errType, status := range statusMap {
        if errors.As(err, &errType) {
            return status
        }
    }
    return http.StatusInternalServerError
}
```

## Error Mapping Functions

### 1. Database Error Mapping

```go
// pkg/common/db_error.go
func HandleDBError(err error, errMap map[error]error) error {
    // Handle specific database errors
    switch {
    case errors.Is(err, gorm.ErrRecordNotFound):
        if mappedErr, ok := errMap[gorm.ErrRecordNotFound]; ok {
            return mappedErr
        }
        return err
    case mysqlErr, ok := err.(*mysql.MySQLError); ok:
        // Handle MySQL specific errors
        switch mysqlErr.Number {
        case 1062: // Duplicate entry
            return handleDuplicateError(mysqlErr, errMap)
        case 1452: // Foreign key constraint fails
            return handleForeignKeyError(mysqlErr, errMap)
        }
    }
    
    // Check error map for direct matches
    if mappedErr, ok := errMap[err]; ok {
        return mappedErr
    }
    
    return err
}

func handleDuplicateError(err *mysql.MySQLError, errMap map[error]error) error {
    // Extract field name from error message
    field := extractFieldFromError(err.Message)
    return &errorz.DuplicateResourceError{
        Field: field,
        Message: fmt.Sprintf("duplicate entry for %s", field),
    }
}
```

### 2. Domain Error Mapping

```go
// domain/availability/errorz.go
func MapDomainError(err error) error {
    switch err := err.(type) {
    case *models.ValidationError:
        return &errorz.ValidationError{
            Field:   err.Field,
            Message: err.Message,
        }
    case *models.BusinessError:
        return &errorz.BusinessError{
            Code:    err.Code,
            Message: err.Message,
        }
    default:
        return err
    }
}
```

### 3. API Error Mapping

```go
// pkg/responses/error_handler.go
func MapAPIError(err error) APIError {
    switch e := err.(type) {
    case *errorz.ValidationError:
        return APIError{
            Code:    "VALIDATION_ERROR",
            Message: e.Message,
            Details: map[string]interface{}{
                "field": e.Field,
            },
        }
    case *errorz.ResourceNotFoundError:
        return APIError{
            Code:    "NOT_FOUND",
            Message: fmt.Sprintf("%s not found", e.Resource),
            Details: map[string]interface{}{
                "resource": e.Resource,
                "id":      e.ID,
            },
        }
    default:
        return APIError{
            Code:    "INTERNAL_ERROR",
            Message: "An internal error occurred",
        }
    }
}
```

## Error Response Formats

### 1. Validation Error

```json
{
    "code": "VALIDATION_ERROR",
    "message": "Invalid input",
    "details": {
        "field": "slot_duration",
        "reason": "must be between 15 and 120 minutes"
    }
}
```

### 2. Resource Not Found

```json
{
    "code": "NOT_FOUND",
    "message": "Doctor not found",
    "details": {
        "resource": "doctor",
        "id": "550e8400-e29b-41d4-a716-446655440000"
    }
}
```

### 3. Business Error

```json
{
    "code": "SLOT_UNAVAILABLE",
    "message": "Appointment slot is already booked",
    "details": {
        "doctor_id": "123e4567-e89b-12d3-a456-426614174000",
        "start_time": "2024-01-01T09:00:00Z"
    }
}
```

## Best Practices

1. **Layer-Specific Error Mapping**
```go
// Repository layer: Map DB errors to domain errors
if err := r.db.Create(entity).Error; err != nil {
    return common.HandleDBError(err, ErrMap)
}

// Service layer: Map domain errors to API errors
if err := s.repo.Create(entity); err != nil {
    return MapDomainError(err)
}

// Controller layer: Map to HTTP responses
if err := s.service.Create(req); err != nil {
    responses.HandleError(ctx, err)
    return
}
```

2. **Consistent Error Keys**
```go
// Define error codes as constants
const (
    ErrCodeValidation     = "VALIDATION_ERROR"
    ErrCodeNotFound       = "NOT_FOUND"
    ErrCodeUnauthorized   = "UNAUTHORIZED"
    ErrCodeBadRequest     = "BAD_REQUEST"
    ErrCodeInternalError  = "INTERNAL_ERROR"
)
```

3. **Error Context Preservation**
```go
func (r *Repository) Create(ctx context.Context, entity *models.Entity) error {
    if err := r.db.Create(entity).Error; err != nil {
        return fmt.Errorf("failed to create %s: %w",
            entity.TableName(),
            common.HandleDBError(err, ErrMap))
    }
    return nil
}
```

4. **Safe Error Messages**
```go
// Internal error details for logging
logger.Error("database connection failed",
    "error", err,
    "host", config.DBHost,
    "database", config.DBName,
)

// Public error response
return APIError{
    Code:    "DATABASE_ERROR",
    Message: "Service temporarily unavailable",
}
