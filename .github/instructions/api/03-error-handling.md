# API Error Handling Guide

This guide explains our API error handling patterns and standards.

## Error Types

1. **Validation Errors**
   - Request payload validation
   - Input parameter validation
   - Business rule validation

2. **Domain Errors**
   - Business logic errors
   - Entity not found
   - Permission denied
   - Invalid state transitions

3. **Infrastructure Errors**
   - Database errors
   - External service errors
   - Network errors
   - System errors

## Error Response Format

Our standard error response structure:

```json
{
    "error": "Error message here",
    "code": "ERROR_CODE",        // Optional error code
    "details": { ... }          // Optional additional details
}
```

## Error Handling Functions

### 1. Validation Errors

Use for input validation failures:

```go
// For request binding/validation errors
if err := ctx.ShouldBindJSON(&req); err != nil {
    responses.HandleValidationError(ctx, c.logger, err)
    return
}
```

Example response:
```json
{
    "error": "Validation failed",
    "code": "VALIDATION_ERROR",
    "details": {
        "slot_duration": "must be between 15 and 120 minutes",
        "end_time": "must be after start_time"
    }
}
```

### 2. General Errors

Use for domain or business logic errors:

```go
// For domain errors or other errors
if err != nil {
    responses.HandleError(ctx, c.logger, err)
    return
}
```

Example response:
```json
{
    "error": "Doctor not found",
    "code": "NOT_FOUND",
    "details": {
        "id": "550e8400-e29b-41d4-a716-446655440000"
    }
}
```

### 3. Specific Status Errors

Use when you need to specify an HTTP status code:

```go
// For errors with specific HTTP status codes
responses.HandleErrorWithStatus(ctx, c.logger, http.StatusForbidden, err)
```

Example response:
```json
{
    "error": "Insufficient permissions",
    "code": "FORBIDDEN",
    "details": {
        "required_role": "doctor",
        "user_role": "patient"
    }
}
```

## Error Mapping

Map domain errors to HTTP status codes:

```go
var errorStatusMap = map[error]int{
    ErrNotFound:     http.StatusNotFound,
    ErrInvalidInput: http.StatusBadRequest,
    ErrForbidden:    http.StatusForbidden,
}

func getStatusCode(err error) int {
    if status, ok := errorStatusMap[err]; ok {
        return status
    }
    return http.StatusInternalServerError
}
```

## Error Logging

Always log errors with context:

```go
func (c *Controller) CreateAvailability(ctx *gin.Context) {
    var req CreateAvailabilityRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        c.logger.Error("invalid request payload",
            "error", err,
            "method", "CreateAvailability",
            "user_id", ctx.GetString("user_id"),
        )
        responses.HandleValidationError(ctx, c.logger, err)
        return
    }
}
```

## Best Practices

1. **Use Domain-Specific Errors**
   ```go
   var (
       ErrAvailabilityNotFound = errorz.NewNotFoundError("availability slot not found")
       ErrInvalidSlotDuration  = errorz.NewBadRequestError("invalid slot duration")
       ErrSlotOverlap         = errorz.NewBadRequestError("slot overlaps with existing availability")
   )
   ```

2. **Include Helpful Error Details**
   ```go
   responses.DetailError(ctx, http.StatusBadRequest, "invalid slot duration",
       map[string]interface{}{
           "min_duration": 15,
           "max_duration": 120,
           "provided": req.SlotDuration,
       },
   )
   ```

3. **Log at Appropriate Levels**
   - Validation errors: WARN
   - Domain errors: ERROR
   - Infrastructure errors: ERROR
   - Security violations: ERROR
   
4. **Handle Panics**
   ```go
   router.Use(middlewares.RecoveryMiddleware(logger))
   ```

5. **Map Database Errors**
   ```go
   var RepositoryErrMap = map[error]bool{
       ErrEntityNotFound: true,
       ErrDuplicateKey:  true,
   }

   if err := r.db.Create(entity).Error; err != nil {
       return common.HandleDBError(err, RepositoryErrMap)
   }
   ```

6. **Return Safe Error Messages**
   - Never expose internal error details
   - Always map to user-friendly messages
   - Include error codes for client handling
   - Add details only when safe

7. **Use Middleware for Common Errors**
   ```go
   router.Use(
       middlewares.AuthMiddleware(),
       middlewares.RateLimitMiddleware(),
       middlewares.ValidationMiddleware(),
   )
   ```
