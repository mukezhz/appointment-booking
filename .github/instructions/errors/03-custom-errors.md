# Custom Error Guide

This guide explains how to create and use custom errors in our application.

## Creating Custom Errors

### 1. Define Error Structure

```go
// pkg/errorz/custom_errors.go
type BusinessError struct {
    Code    string                 `json:"code"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}

func (e *BusinessError) Error() string {
    return e.Message
}

type ValidationError struct {
    Field   string                 `json:"field"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

### 2. Create Error Constructors

```go
// pkg/errorz/constructors.go
func NewBusinessError(code, message string) *BusinessError {
    return &BusinessError{
        Code:    code,
        Message: message,
        Details: make(map[string]interface{}),
    }
}

func (e *BusinessError) WithDetails(details map[string]interface{}) *BusinessError {
    e.Details = details
    return e
}

func NewValidationError(field, message string) *ValidationError {
    return &ValidationError{
        Field:   field,
        Message: message,
        Details: make(map[string]interface{}),
    }
}

func (e *ValidationError) WithDetails(details map[string]interface{}) *ValidationError {
    e.Details = details
    return e
}
```

## Domain-Specific Errors

### 1. Availability Domain

```go
// domain/availability/errorz.go
var (
    ErrSlotUnavailable = NewBusinessError(
        "SLOT_UNAVAILABLE",
        "appointment slot is not available",
    )

    ErrInvalidSlotDuration = NewValidationError(
        "slot_duration",
        "slot duration must be between 15 and 120 minutes",
    )

    ErrInvalidTimeRange = NewValidationError(
        "time_range",
        "end time must be after start time",
    )
)

// With context
func NewSlotUnavailableError(doctorID string, startTime time.Time) error {
    return ErrSlotUnavailable.WithDetails(map[string]interface{}{
        "doctor_id":   doctorID,
        "start_time": startTime,
    })
}
```

### 2. Appointment Domain

```go
// domain/appointment/errorz.go
var (
    ErrPastBooking = NewBusinessError(
        "PAST_BOOKING",
        "cannot book appointments in the past",
    )

    ErrBookingLimit = NewBusinessError(
        "BOOKING_LIMIT",
        "maximum appointments limit reached",
    )
)

// With context
func NewBookingLimitError(limit int, current int) error {
    return ErrBookingLimit.WithDetails(map[string]interface{}{
        "limit":   limit,
        "current": current,
    })
}
```

## Using Custom Errors

### 1. Service Layer

```go
func (s *Service) CreateAvailability(ctx context.Context, req *CreateAvailabilityRequest) error {
    // Validate slot duration
    if req.SlotDuration < 15 || req.SlotDuration > 120 {
        return NewValidationError("slot_duration", "must be between 15 and 120 minutes").
            WithDetails(map[string]interface{}{
                "min":      15,
                "max":      120,
                "provided": req.SlotDuration,
            })
    }

    // Check for overlapping slots
    if overlaps, err := s.checkOverlap(ctx, req); err != nil {
        return err
    } else if overlaps {
        return NewSlotUnavailableError(req.DoctorID, req.StartTime)
    }

    return nil
}
```

### 2. Error Handling

```go
func (c *Controller) CreateAvailability(ctx *gin.Context) {
    var req CreateAvailabilityRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        responses.HandleValidationError(ctx, c.logger, err)
        return
    }

    if err := c.service.CreateAvailability(ctx.Request.Context(), &req); err != nil {
        switch e := err.(type) {
        case *errorz.ValidationError:
            responses.HandleValidationError(ctx, c.logger, e)
        case *errorz.BusinessError:
            responses.HandleBusinessError(ctx, c.logger, e)
        default:
            responses.HandleError(ctx, c.logger, err)
        }
        return
    }

    // Success response
}
```

## Error Response Examples

### 1. Validation Error

```json
{
    "error": {
        "code": "VALIDATION_ERROR",
        "field": "slot_duration",
        "message": "must be between 15 and 120 minutes",
        "details": {
            "min": 15,
            "max": 120,
            "provided": 10
        }
    }
}
```

### 2. Business Error

```json
{
    "error": {
        "code": "SLOT_UNAVAILABLE",
        "message": "appointment slot is not available",
        "details": {
            "doctor_id": "550e8400-e29b-41d4-a716-446655440000",
            "start_time": "2024-01-01T09:00:00Z"
        }
    }
}
```

## Best Practices

1. **Error Grouping**
```go
// Group related errors together
var (
    // Validation errors
    ErrInvalidSlotDuration = NewValidationError(...)
    ErrInvalidTimeRange    = NewValidationError(...)

    // Business logic errors
    ErrSlotUnavailable    = NewBusinessError(...)
    ErrPastBooking        = NewBusinessError(...)

    // Resource errors
    ErrDoctorNotFound     = NewNotFoundError(...)
    ErrPatientNotFound    = NewNotFoundError(...)
)
```

2. **Contextual Information**
```go
// Add relevant context to errors
return ErrSlotUnavailable.WithDetails(map[string]interface{}{
    "doctor_id":  doctorID,
    "date":      date.Format("2006-01-02"),
    "start_time": startTime.Format("15:04"),
    "reason":     "doctor unavailable",
})
```

3. **User-Friendly Messages**
```go
// Use clear, actionable messages
var (
    ErrInvalidTimeRange = NewValidationError(
        "time_range",
        "end time must be after start time and within business hours (8 AM - 6 PM)",
    )

    ErrSlotUnavailable = NewBusinessError(
        "SLOT_UNAVAILABLE",
        "this time slot is already booked, please choose another time",
    )
)
```

4. **Error Documentation**
```go
// Document error types and their use cases
// ErrBookingLimit is returned when a patient attempts to book more
// appointments than allowed within a given time period.
var ErrBookingLimit = NewBusinessError(
    "BOOKING_LIMIT",
    "maximum appointments limit reached",
).WithDetails(map[string]interface{}{
    "max_per_day":   3,
    "max_per_week":  10,
    "max_per_month": 20,
})
```
