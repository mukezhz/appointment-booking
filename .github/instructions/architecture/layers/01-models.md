# Domain Models Layer

## Overview

Domain models form the core of our application. They represent business entities and implement domain-specific logic. Located in `domain/models/`, these models are independent of frameworks and external concerns.

## Key Concepts

1. **Rich Domain Models**
   - Business logic in models
   - Framework-independent
   - Self-validating
   - Clear relationships

2. **Database Mapping**
   - GORM annotations
   - Clean persistence
   - Proper indexing
   - Type safety

## Example Model

```go
// domain/models/availability.go
type Availability struct {
    gorm.Model
    UUID         types.BinaryUUID `json:"uuid" gorm:"type:binary(16);uniqueIndex;not null"`
    DoctorID     types.BinaryUUID `json:"doctor_id" gorm:"type:binary(16);index;not null"`
    StartDate    *time.Time       `json:"start_date,omitempty" gorm:"type:date;index"`
    EndDate      *time.Time       `json:"end_date,omitempty" gorm:"type:date"`
    StartTime    time.Time        `json:"start_time" gorm:"type:time;not null"`
    EndTime      time.Time        `json:"end_time" gorm:"type:time;not null"`
    SlotDuration int              `json:"slot_duration" gorm:"not null"`
    IsRecurring  bool             `json:"is_recurring" gorm:"default:false"`

    // Domain logic methods
    func (a *Availability) OverlapsWith(other *Availability) bool
    func (a *Availability) GenerateSlots(date time.Time) []Slot
    func (a *Availability) IsAvailableAt(date time.Time, startTime time.Time) bool
}
```

## Implementation Guide

### 1. Define the Model

```go
type YourModel struct {
    gorm.Model                     // Base fields (ID, timestamps)
    UUID      types.BinaryUUID     // Unique identifier
    Name      string               // Business fields
    Status    YourModelStatus      // Custom types for enums
    // Add other fields
}
```

### 2. Add Database Tags

```go
type YourModel struct {
    UUID      types.BinaryUUID `gorm:"type:binary(16);uniqueIndex;not null"`
    Name      string           `gorm:"type:varchar(255);not null"`
    Status    YourModelStatus  `gorm:"type:varchar(20);default:active"`
}
```

### 3. Add JSON Tags

```go
type YourModel struct {
    UUID      types.BinaryUUID `json:"uuid"`
    Name      string           `json:"name"`
    Status    YourModelStatus  `json:"status"`
}
```

### 4. Implement Domain Logic

```go
// Validation
func (m *YourModel) Validate() error {
    if m.Name == "" {
        return ErrInvalidName
    }
    return nil
}

// Business Rules
func (m *YourModel) CanBeModified() bool {
    return m.Status != YourModelStatusLocked
}

// Calculations
func (m *YourModel) CalculateMetric() float64 {
    // Implement business calculations
    return result
}
```

### 5. Add Table Name (Optional)

```go
func (m *YourModel) TableName() string {
    return "your_models"
}
```

## Best Practices

1. **Model Design**
   - Keep models focused and cohesive
   - Use custom types for enums/special fields
   - Include validation logic
   - Implement domain methods

2. **Database Mapping**
   - Use appropriate GORM tags
   - Set proper column types
   - Add necessary indexes
   - Handle relationships properly

3. **Validation**
   - Validate all required fields
   - Check business rule constraints
   - Return domain-specific errors
   - Log validation failures

4. **Domain Logic**
   - Keep business rules in models
   - Use meaningful method names
   - Document complex logic
   - Write unit tests

5. **Performance**
   - Use appropriate field types
   - Index frequently queried fields
   - Optimize relationships
   - Consider query patterns

## Common Patterns

### 1. Status Enums
```go
type ModelStatus string

const (
    StatusActive   ModelStatus = "active"
    StatusInactive ModelStatus = "inactive"
    StatusDeleted  ModelStatus = "deleted"
)
```

### 2. Relationships
```go
type Parent struct {
    gorm.Model
    Children []Child `gorm:"foreignKey:ParentID"`
}

type Child struct {
    gorm.Model
    ParentID uint
    Parent   Parent `gorm:"foreignKey:ParentID"`
}
```

### 3. Validation Methods
```go
func (m *Model) Validate() error {
    switch {
    case m.Name == "":
        return ErrNameRequired
    case len(m.Description) > 500:
        return ErrDescriptionTooLong
    case !m.IsValidStatus():
        return ErrInvalidStatus
    }
    return nil
}
```

### 4. Custom Types
```go
type BinaryUUID [16]byte

func (u BinaryUUID) String() string {
    return fmt.Sprintf("%x-%x-%x-%x-%x",
        u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
```

## Testing Models

```go
Describe("YourModel", func() {
    var model *models.YourModel

    BeforeEach(func() {
        model = &models.YourModel{
            Name:   "Test Model",
            Status: models.StatusActive,
        }
    })

    Describe("Validate", func() {
        It("should pass for valid model", func() {
            Expect(model.Validate()).To(BeNil())
        })

        It("should fail for invalid name", func() {
            model.Name = ""
            Expect(model.Validate()).To(Equal(models.ErrNameRequired))
        })
    })

    Describe("Business Logic", func() {
        It("should calculate correctly", func() {
            result := model.CalculateMetric()
            Expect(result).To(Equal(expectedValue))
        })
    })
})
