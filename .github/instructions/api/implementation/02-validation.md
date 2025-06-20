# Validation Guide

This guide explains how to implement request validation in our REST APIs.

## Input Validation

### 1. Binding Tags

Use Gin's binding tags to validate request DTOs:

```go
type CreateAvailabilityRequest struct {
    StartTime    time.Time `json:"start_time" binding:"required"`
    EndTime      time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
    SlotDuration int       `json:"slot_duration" binding:"required,min=15,max=120"`
    IsRecurring  bool      `json:"is_recurring" binding:"required"`
}
```

Common binding tags:
- `required`: Field must be present
- `omitempty`: Field can be omitted
- `min`: Minimum value for numbers
- `max`: Maximum value for numbers
- `len`: Exact length requirement
- `email`: Must be valid email format
- `url`: Must be valid URL format
- `uuid`: Must be valid UUID format

### 2. Custom Validation

Implement custom validation functions for complex rules:

```go
func ValidateAvailability(availability *models.Availability) error {
    // 1. Check business hours
    if availability.StartTime.Hour() < 8 || availability.EndTime.Hour() > 18 {
        return ErrInvalidBusinessHours
    }

    // 2. Check slot duration
    if availability.SlotDuration < 15 || availability.SlotDuration > 120 {
        return ErrInvalidSlotDuration
    }

    // 3. Check time sequence
    if availability.EndTime.Before(availability.StartTime) {
        return ErrInvalidTimeSequence
    }

    return nil
}
```

## Validation in Controllers

### 1. Request Validation

```go
func (c *Controller) CreateAvailability(ctx *gin.Context) {
    // 1. Parse and validate request
    var req CreateAvailabilityRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        responses.HandleValidationError(ctx, c.logger, err)
        return
    }

    // 2. Convert to domain model
    availability := &models.Availability{
        UUID:         types.NewBinaryUUID(),
        DoctorID:     ctx.MustGet(framework.UserIDKey).(types.BinaryUUID),
        StartTime:    req.StartTime,
        EndTime:      req.EndTime,
        SlotDuration: req.SlotDuration,
        IsRecurring:  req.IsRecurring,
    }

    // 3. Validate domain model
    if err := ValidateAvailability(availability); err != nil {
        responses.HandleValidationError(ctx, c.logger, err)
        return
    }

    // ... rest of the implementation
}
```

### 2. Path Parameter Validation

```go
func (c *Controller) GetAvailability(ctx *gin.Context) {
    // Validate UUID format
    id, err := types.ParseUUID(ctx.Param("id"))
    if err != nil {
        responses.HandleValidationError(ctx, c.logger,
            errors.New("invalid UUID format"))
        return
    }

    // ... rest of the implementation
}
```

### 3. Query Parameter Validation

```go
type ListAvailabilityParams struct {
    StartDate time.Time `form:"start_date" binding:"required" time_format:"2006-01-02"`
    EndDate   time.Time `form:"end_date" binding:"required,gtfield=StartDate" time_format:"2006-01-02"`
    DoctorID  string    `form:"doctor_id" binding:"omitempty,uuid4"`
    Page      int       `form:"page" binding:"min=1"`
    PerPage   int       `form:"per_page" binding:"required,min=1,max=100"`
}

func (c *Controller) ListAvailability(ctx *gin.Context) {
    var params ListAvailabilityParams
    if err := ctx.ShouldBindQuery(&params); err != nil {
        responses.HandleValidationError(ctx, c.logger, err)
        return
    }

    // ... rest of the implementation
}
```

## Custom Validators

### 1. Register Custom Validators

```go
func RegisterCustomValidators(v *validator.Validate) {
    // Custom time validator
    v.RegisterValidation("business_hours", validateBusinessHours)

    // Custom format validator
    v.RegisterValidation("phone_number", validatePhoneNumber)
}

// Business hours validator (8 AM to 6 PM)
func validateBusinessHours(fl validator.FieldLevel) bool {
    if t, ok := fl.Field().Interface().(time.Time); ok {
        hour := t.Hour()
        return hour >= 8 && hour <= 18
    }
    return false
}

// Phone number validator
func validatePhoneNumber(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    match, _ := regexp.MatchString(`^\+[1-9]\d{1,14}$`, phone)
    return match
}
```

### 2. Use Custom Validators

```go
type CreateAppointmentRequest struct {
    DoctorID    string    `json:"doctor_id" binding:"required,uuid4"`
    PatientID   string    `json:"patient_id" binding:"required,uuid4"`
    StartTime   time.Time `json:"start_time" binding:"required,business_hours"`
    PhoneNumber string    `json:"phone_number" binding:"required,phone_number"`
}
```

## Validation Response Format

### 1. Field Validation Error

```json
{
    "error": "Validation failed",
    "code": "VALIDATION_ERROR",
    "details": {
        "start_time": "must be within business hours (8 AM - 6 PM)",
        "phone_number": "must be a valid phone number",
        "slot_duration": "must be between 15 and 120 minutes"
    }
}
```

### 2. Custom Validation Error

```json
{
    "error": "Invalid appointment request",
    "code": "VALIDATION_ERROR",
    "details": {
        "reason": "Appointment slot overlaps with existing booking",
        "existing_appointment": {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "start_time": "2024-01-01T10:00:00Z",
            "end_time": "2024-01-01T10:30:00Z"
        }
    }
}
```

## Best Practices

1. **Layer-Specific Validation**
   - Controllers: Request format and type validation
   - Services: Business rule validation
   - Models: Domain logic validation

2. **Validation Chains**
   ```go
   func (s *Service) CreateAppointment(ctx context.Context, req *CreateAppointmentRequest) (*models.Appointment, error) {
       // 1. Basic validation
       if err := validateRequest(req); err != nil {
           return nil, err
       }

       // 2. Business rule validation
       if err := s.validateBusinessRules(ctx, req); err != nil {
           return nil, err
       }

       // 3. State validation
       if err := s.validateState(ctx, req); err != nil {
           return nil, err
       }

       // Proceed with creation
       return s.repo.Create(ctx, appointment)
   }
   ```

3. **Reusable Validation Functions**
   ```go
   package validators

   func ValidateTimeRange(start, end time.Time) error {
       if end.Before(start) {
           return ErrInvalidTimeRange
       }
       if end.Sub(start) > 24*time.Hour {
           return ErrTimeTooLong
       }
       return nil
   }
   ```

4. **Contextual Validation**
   ```go
   func (s *Service) validateAppointmentAvailability(
       ctx context.Context,
       doctorID types.BinaryUUID,
       startTime time.Time,
   ) error {
       // Check doctor's schedule
       available, err := s.availabilityRepo.IsAvailable(ctx, doctorID, startTime)
       if err != nil {
           return err
       }
       if !available {
           return ErrDoctorUnavailable
       }

       // Check existing appointments
       exists, err := s.appointmentRepo.HasOverlapping(ctx, doctorID, startTime)
       if err != nil {
           return err
       }
       if exists {
           return ErrSlotTaken
       }

       return nil
   }
   ```
