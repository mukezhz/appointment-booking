# Request and Response Standards

This guide outlines our standard patterns for handling API requests and responses.

## Standard Response Formats

We use standard response types from the `responses` package:

### 1. Detail Response (Single Item)
```go
type DetailResponseType[T any] struct {
    Item    T      `json:"item"`
    Message string `json:"message"`
}

// Usage:
responses.DetailResponse(ctx, http.StatusOK, responses.DetailResponseType[UserDTO]{
    Item:    userDTO,
    Message: "User retrieved successfully",
})
```

### 2. List Response (Multiple Items)
```go
type ListResponseType[T any] struct {
    Items   []T    `json:"items"`
    Message string `json:"message"`
}

// Usage:
responses.ListResponse(ctx, http.StatusOK, responses.ListResponseType[UserDTO]{
    Items:   userDTOs,
    Message: "Users retrieved successfully",
})
```

### 3. Paginated Response
```go
type PaginatedResponseType[T any] struct {
    Items      []T    `json:"items"`
    Message    string `json:"message"`
    TotalItems int64  `json:"total_items"`
    TotalPages int    `json:"total_pages"`
    Page       int    `json:"page"`
    PerPage    int    `json:"per_page"`
}

// Usage:
responses.PaginatedResponse(ctx, http.StatusOK, responses.PaginatedResponseType[UserDTO]{
    Items:      userDTOs,
    Message:    "Users retrieved successfully",
    TotalItems: total,
    TotalPages: totalPages,
    Page:       page,
    PerPage:    perPage,
})
```

## DTOs and Model Conversion

### Request DTOs
```go
type CreateAvailabilityRequest struct {
    StartTime    time.Time `json:"start_time" binding:"required"`
    EndTime      time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
    SlotDuration int       `json:"slot_duration" binding:"required,min=15,max=120"`
    IsRecurring  bool      `json:"is_recurring" binding:"required"`
}
```

### Response DTOs
```go
type AvailabilityDTO struct {
    UUID         string    `json:"uuid"`
    DoctorID     string    `json:"doctor_id"`
    StartTime    time.Time `json:"start_time"`
    EndTime      time.Time `json:"end_time"`
    SlotDuration int       `json:"slot_duration"`
    IsRecurring  bool      `json:"is_recurring"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### Model to DTO Conversion
```go
func (c *Controller) toAvailabilityDTO(m *models.Availability) *AvailabilityDTO {
    if m == nil {
        return nil
    }
    return &AvailabilityDTO{
        UUID:         m.UUID.String(),
        DoctorID:     m.DoctorID.String(),
        StartTime:    m.StartTime,
        EndTime:      m.EndTime,
        SlotDuration: m.SlotDuration,
        IsRecurring:  m.IsRecurring,
        CreatedAt:    m.CreatedAt,
        UpdatedAt:    m.UpdatedAt,
    }
}
```

## Error Response Handling

### 1. Validation Errors
```go
// For request binding/validation errors
if err := ctx.ShouldBindJSON(&req); err != nil {
    responses.HandleValidationError(ctx, c.logger, err)
    return
}
```

### 2. General Errors
```go
// For domain errors or other errors
if err != nil {
    responses.HandleError(ctx, c.logger, err)
    return
}
```

### 3. Specific Status Errors
```go
// For errors with specific HTTP status codes
responses.HandleErrorWithStatus(ctx, c.logger, http.StatusForbidden, err)
```

### Error Response Format
```json
{
    "error": "Error message here",
    "code": "ERROR_CODE",        // Optional error code
    "details": { ... }          // Optional additional details
}
```

## Standard Examples

### 1. Success Response (Single Item)
```json
{
    "item": {
        "uuid": "550e8400-e29b-41d4-a716-446655440000",
        "doctor_id": "123e4567-e89b-12d3-a456-426614174000",
        "start_time": "2024-01-01T09:00:00Z",
        "end_time": "2024-01-01T17:00:00Z",
        "slot_duration": 30,
        "is_recurring": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
    },
    "message": "Availability slot created successfully"
}
```

### 2. Success Response (List)
```json
{
    "items": [
        {
            "uuid": "550e8400-e29b-41d4-a716-446655440000",
            "doctor_id": "123e4567-e89b-12d3-a456-426614174000",
            "start_time": "2024-01-01T09:00:00Z",
            "end_time": "2024-01-01T17:00:00Z",
            "slot_duration": 30,
            "is_recurring": true,
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
        }
    ],
    "message": "Availability slots retrieved successfully",
    "total_items": 1,
    "total_pages": 1,
    "page": 1,
    "per_page": 10
}
```

### 3. Error Response Example
```json
{
    "error": "Invalid availability slot duration",
    "code": "INVALID_DURATION",
    "details": {
        "min_duration": 15,
        "max_duration": 120,
        "provided_duration": 10
    }
}
