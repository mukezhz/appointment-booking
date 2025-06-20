# Service Layer Guide

The service layer contains all business logic. This guide explains how to implement and modify services in our clean architecture.

## Key Responsibilities

- Implements business rules and workflows
- Validates domain models
- Coordinates between repositories
- Handles domain-specific errors
- Works with models, not DTOs

## Service Structure

```go
type Service struct {
    repo    *Repository          // Primary repository
    logger  framework.Logger     // Structured logger
    // Optional: other dependencies
    userRepo   *user.Repository     // Other domain repositories if needed
    aws       *services.AWSService  // External services if needed
}

// Validation functions for domain models
func ValidateEntity(entity *models.Entity) error {
    if entity.Name == "" {
        return ErrInvalidInput
    }
    // Add more validation rules
    return nil
}
```

## Service Method Pattern

Service methods always:
1. Accept and return models (not DTOs)
2. Use context for cancellation/timeouts
3. Validate inputs
4. Log operations and errors
5. Return domain-specific errors

Example implementation:

```go
func (s *Service) CreateEntity(ctx context.Context, doctorID types.BinaryUUID, data *models.Entity) (*models.Entity, error) {
    s.logger.Info("creating entity for doctor: ", doctorID.String())

    // 1. Input Validation
    if err := ValidateEntity(data); err != nil {
        s.logger.Error("validation failed: ", err)
        return nil, err
    }

    // 2. Business Logic
    entity := &models.Entity{
        UUID:     types.NewBinaryUUID(),
        DoctorID: doctorID,
        // ... set other fields based on business rules
    }

    // 3. Domain Rules
    if err := s.checkBusinessRules(ctx, entity); err != nil {
        return nil, err
    }

    // 4. Persistence
    if err := s.repo.Create(ctx, entity); err != nil {
        s.logger.Error("failed to create entity: ", err)
        return nil, err
    }

    s.logger.Info("created entity: ", entity.UUID.String())
    return entity, nil
}
```

## Implementation Guidelines

When implementing a service:

1. Keep business logic in the service layer
2. Use domain models, not DTOs
3. Handle error cases explicitly
4. Log operations and errors
5. Implement validation functions
6. Use meaningful error messages

## Testing Services

Services should have comprehensive unit tests:

1. Test happy path scenarios
2. Test validation cases
3. Test business rule violations
4. Mock dependencies (repositories, external services)
5. Test error handling

Example test structure:
```go
func TestService_CreateEntity(t *testing.T) {
    tests := []struct {
        name    string
        input   *models.Entity
        wantErr bool
        errType error
    }{
        // Define test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mocks
            // Run test
            // Assert results
        })
    }
}
```
