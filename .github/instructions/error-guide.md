# Complete Error Handling Guide

> This guide consolidates all error handling documentation into a single, comprehensive resource.

## Table of Contents
- [Complete Error Handling Guide](#complete-error-handling-guide)
  - [Table of Contents](#table-of-contents)
  - [Error System Architecture](#error-system-architecture)
    - [1. Base Error Types (`pkg/errorz`)](#1-base-error-types-pkgerrorz)
    - [2. Error Creation Helpers](#2-error-creation-helpers)
    - [3. Database Error Handling (`domain/common`)](#3-database-error-handling-domaincommon)
    - [4. HTTP Response Handling (`pkg/responses`)](#4-http-response-handling-pkgresponses)
  - [Best Practices](#best-practices)
    - [1. Repository Layer](#1-repository-layer)
    - [2. Service Layer](#2-service-layer)
    - [3. Controller Layer](#3-controller-layer)
  - [Error Types Hierarchy](#error-types-hierarchy)
  - [Integration with External Systems](#integration-with-external-systems)
    - [AWS Error Handling](#aws-error-handling)
    - [Database Error Mapping](#database-error-mapping)
  - [Logging and Monitoring](#logging-and-monitoring)
  - [Error Response Format](#error-response-format)
  - [Testing Error Handling](#testing-error-handling)
  - [API Error Handling](#api-error-handling)

## Error System Architecture

### 1. Base Error Types (`pkg/errorz`)
```go
// APIError - Base error type with HTTP status
type APIError struct {
    StatusCode int
    Message    string
}

// Pre-defined errors
ErrBadRequest         // 400
ErrUnauthorized       // 401
ErrForbidden          // 403
ErrNotFound           // 404
ErrConflict           // 409
ErrUnprocessable      // 422
ErrInternal           // 500
ErrServiceUnavailable // 503
```

### 2. Error Creation Helpers
```go
// Create new errors
NewNotFoundError(message string)
NewBadRequestError(message string)
NewUnauthorizedError(message string)
NewForbiddenError(message string)
NewInternalError(message string)

// Join errors
JoinError(message string, base error)

// Wrap errors with context
Wrap(err error, message string)
```

### 3. Database Error Handling (`domain/common`)
```go
// Maps DB errors to application errors
func HandleDBError(err error, domainErrs map[error]bool) error {
    // Maps common DB errors:
    // - ErrRecordNotFound -> ErrNotFound
    // - ErrInvalidData -> ErrBadRequest
    // - ErrDuplicatedKey -> ErrAlreadyExists
    // etc...
}
```

### 4. HTTP Response Handling (`pkg/responses`)
```go
// Handle different types of errors in HTTP responses
HandleError(ctx *gin.Context, logger framework.Logger, err error)
HandleValidationError(ctx *gin.Context, logger framework.Logger, err error)
HandleErrorWithStatus(ctx *gin.Context, logger framework.Logger, statusCode int, err error)
```

## Best Practices

### 1. Repository Layer
```go
func (r *Repository) Operation(ctx context.Context) error {
    if err := r.db.WithContext(ctx).Operation().Error; err != nil {
        r.logger.Error("operation failed: ", err)
        return common.HandleDBError(err, DomainErrMap)
    }
    return nil
}
```

### 2. Service Layer
```go
func (s *Service) Operation(ctx context.Context) error {
    // Use error creation helpers
    if invalidInput {
        return NewBadRequestError("invalid input")
    }
    
    // Wrap errors with context
    if err := s.repo.Operation(ctx); err != nil {
        return Wrap(err, "operation failed")
    }
    return nil
}
```

### 3. Controller Layer
```go
func (c *Controller) Handle(ctx *gin.Context) {
    if err := c.service.Operation(ctx); err != nil {
        responses.HandleError(ctx, c.logger, err)
        return
    }
    // Success response...
}
```

## Error Types Hierarchy

1. **Base Errors** (pkg/errorz/base.go)
   - APIError structure
   - HTTP status based errors

2. **Common Errors** (pkg/errorz/common_errors.go)
   - Domain-specific common errors
   - Reusable error types

3. **Domain Errors** (domain/*/errorz.go)
   - Module-specific errors
   - Custom error types

## Integration with External Systems

### AWS Error Handling
```go
// AWS error mapping (pkg/utils/aws_error_mapper.go)
func MapAWSError(logger framework.Logger, err error) *AWSError {
    // Maps AWS SDK errors to application errors
}
```

### Database Error Mapping
```go
// Database error mapping (domain/common/db_error.go)
var ErrUserMap = map[error]bool{
    ErrUserNotFound: true,
    ErrUserExists: true,
}
```

## Logging and Monitoring

1. **Error Logging**
   - Always log errors at the source
   - Include context and stack traces
   - Use appropriate log levels

2. **Error Monitoring**
   - Integration with Sentry
   - Capture unhandled exceptions
   - Track error frequencies

## Error Response Format

Standard JSON error response:
```json
{
    "error": "Error message"
}
```

## Testing Error Handling

1. Test error creation
2. Test error wrapping
3. Test error handling middleware
4. Test database error mapping
5. Test HTTP response handling

## API Error Handling

### Request Validation
1. **Input Validation**
   ```go
   func (c *Controller) Handle(ctx *gin.Context) {
       var req RequestDTO
       if err := ctx.ShouldBindJSON(&req); err != nil {
           responses.HandleValidationError(ctx, c.logger, err)
           return
       }
   }
   ```

2. **Business Rule Validation**
   ```go
   if !isValid {
       err := NewBadRequestError("invalid business rule")
       responses.HandleError(ctx, c.logger, err)
       return
   }
   ```

### Error Response Standards
1. **Validation Errors** (400 Bad Request)
   ```json
   {
       "error": "Invalid input: field 'name' is required"
   }
   ```

2. **Authentication Errors** (401 Unauthorized)
   ```json
   {
       "error": "Invalid or expired token"
   }
   ```

3. **Permission Errors** (403 Forbidden)
   ```json
   {
       "error": "Insufficient permissions to access resource"
   }
   ```

4. **Resource Errors** (404 Not Found)
   ```json
   {
       "error": "Resource not found"
   }
   ```

5. **Conflict Errors** (409 Conflict)
   ```json
   {
       "error": "Resource already exists"
   }
   ```

6. **Server Errors** (500 Internal Server Error)
   ```json
   {
       "error": "An unexpected error occurred"
   }
   ```

### API Error Best Practices
1. Use consistent error response format
2. Include only safe error details in production
3. Log detailed errors server-side
4. Use appropriate HTTP status codes
5. Handle validation errors separately from other errors
6. Include request ID in responses for tracking
