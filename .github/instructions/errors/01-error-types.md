# Error Types Guide

This guide describes the different error types used in our application and when to use them.

## Base Error Types

The application defines several core error types in `pkg/errorz`:

```go
// Base error type for all application errors
type Error struct {
    Code    int    // HTTP status code
    Message string // User-facing error message
    Err     error  // Original error for logging
}

// Specific error types
type ValidationError struct {
    Field   string // Field that failed validation
    Message string // User-friendly error message
}

type ResourceNotFoundError struct {
    Resource string // Resource type (e.g., "user", "appointment")
    ID      string // Identifier that wasn't found
}

type DuplicateResourceError struct {
    Resource string // Resource type
    Field    string // Field causing the duplicate
    Value    string // Value that caused the duplicate
}
```

## Error Categories

### 1. Validation Errors
Use when:
- Request payload is invalid
- Input parameters fail validation
- Business rule validation fails

Example:
```go
return &ValidationError{
    Field:   "slot_duration",
    Message: "slot duration must be between 15 and 120 minutes",
}
```

### 2. Resource Errors
Use when:
- Entity not found in database
- Duplicate entry detected
- Resource state is invalid

Example:
```go
return &ResourceNotFoundError{
    Resource: "doctor",
    ID:      doctorID.String(),
}
```

### 3. Authorization Errors
Use when:
- User is not authenticated
- User lacks required permissions
- Token is invalid or expired

Example:
```go
return &AuthorizationError{
    Message: "insufficient permissions",
    Required: "doctor",
}
```

### 4. Business Logic Errors
Use when:
- Business rules are violated
- State transitions are invalid
- Operations are not allowed

Example:
```go
return &BusinessError{
    Code:    "SLOT_UNAVAILABLE",
    Message: "appointment slot is already booked",
}
```

### 5. Infrastructure Errors
Use when:
- Database operations fail
- External services are unavailable
- Network issues occur

Example:
```go
return &InfrastructureError{
    Service: "database",
    Message: "connection timeout",
}
```

## Creating Custom Errors

1. **Define Error Type**
```go
type ConflictError struct {
    Resource    string
    Identifier  string
    Reason      string
}

func (e *ConflictError) Error() string {
    return fmt.Sprintf("%s %s: %s", e.Resource, e.Identifier, e.Reason)
}
```

2. **Create Constructor**
```go
func NewConflictError(resource, id, reason string) *ConflictError {
    return &ConflictError{
        Resource:   resource,
        Identifier: id,
        Reason:     reason,
    }
}
```

3. **Add to Error Map**
```go
var ErrMap = map[error]error{
    ErrConflict: &ConflictError{
        Resource: "appointment",
        Reason:   "time slot conflict",
    },
}
```

## Error Wrapping

Use error wrapping to preserve context:

```go
if err := validateInput(req); err != nil {
    return fmt.Errorf("invalid request: %w", err)
}
```

## Best Practices

1. **Be Specific**
   - Use specific error types
   - Include relevant context
   - Provide clear messages

2. **Error Context**
   ```go
   type AppointmentError struct {
       DoctorID  string
       DateTime  time.Time
       Reason    string
   }
   ```

3. **Error Constants**
   ```go
   var (
       ErrSlotUnavailable = &BusinessError{
           Code:    "SLOT_UNAVAILABLE",
           Message: "appointment slot is not available",
       }
       ErrDoctorNotFound = &ResourceNotFoundError{
           Resource: "doctor",
       }
   )
   ```

4. **Public vs Internal**
   - Public errors: User-friendly messages
   - Internal errors: Detailed for logging
   ```go
   // Public error
   return &ValidationError{
       Message: "Invalid appointment time",
   }

   // Internal logging
   logger.Error("appointment validation failed",
       "error", err,
       "doctorID", req.DoctorID,
       "startTime", req.StartTime,
   )
   ```

5. **Error Grouping**
   ```go
   // Group related errors
   var (
       // Validation errors
       ErrInvalidSlotDuration = &ValidationError{...}
       ErrInvalidTimeRange    = &ValidationError{...}

       // Business errors
       ErrSlotUnavailable     = &BusinessError{...}
       ErrPastBooking         = &BusinessError{...}

       // Resource errors
       ErrDoctorNotFound      = &ResourceNotFoundError{...}
       ErrPatientNotFound     = &ResourceNotFoundError{...}
   )
   ```

6. **Error Documentation**
   ```go
   // ErrSlotUnavailable is returned when attempting to book an appointment
   // in a time slot that is already taken or outside availability.
   var ErrSlotUnavailable = &BusinessError{
       Code:    "SLOT_UNAVAILABLE",
       Message: "appointment slot is not available",
   }
   ```
