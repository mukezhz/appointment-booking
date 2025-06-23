# Unit Testing Guide

This guide covers unit testing patterns and best practices for our Go Clean Architecture project.

## Unit Test Structure

### 1. Table-Driven Tests

```go
func TestService_CreateAppointment(t *testing.T) {
    tests := []struct {
        name        string
        input       *CreateAppointmentRequest
        setupMocks  func(*mocks.Repository)
        wantErr     bool
        expectedErr error
    }{
        {
            name: "successful creation",
            input: &CreateAppointmentRequest{
                DoctorID:  "test-doctor",
                DateTime:  time.Now().Add(24 * time.Hour),
                Duration: 30,
            },
            setupMocks: func(r *mocks.Repository) {
                r.On("Create", mock.Anything, mock.Anything).Return(nil)
            },
            wantErr: false,
        },
        {
            name: "invalid duration",
            input: &CreateAppointmentRequest{
                Duration: 10, // Too short
            },
            wantErr:     true,
            expectedErr: ErrInvalidDuration,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            mockRepo := new(mocks.Repository)
            if tt.setupMocks != nil {
                tt.setupMocks(mockRepo)
            }
            
            service := NewService(mockRepo)
            
            // Execute
            _, err := service.CreateAppointment(context.Background(), tt.input)
            
            // Assert
            if tt.wantErr {
                assert.Error(t, err)
                assert.Equal(t, tt.expectedErr, err)
            } else {
                assert.NoError(t, err)
            }
            mockRepo.AssertExpectations(t)
        })
    }
}
```

## Testing Different Layers

### 1. Service Layer Tests

```go
func TestAvailabilityService_CreateAvailability(t *testing.T) {
    // Setup
    mockRepo := new(mocks.Repository)
    mockLogger := new(mocks.Logger)
    service := NewService(mockRepo, mockLogger)
    
    // Test data
    availability := &models.Availability{
        DoctorID:     types.NewBinaryUUID(),
        StartTime:    time.Now(),
        EndTime:      time.Now().Add(time.Hour),
        SlotDuration: 30,
    }
    
    // Setup expectations
    mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    // Execute
    result, err := service.CreateAvailability(context.Background(), availability)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, availability.DoctorID, result.DoctorID)
    mockRepo.AssertExpectations(t)
}
```

### 2. Repository Layer Tests

```go
func TestRepository_GetByID(t *testing.T) {
    // Setup mock DB
    db, mock := testutil.NewMockDB()
    repo := NewRepository(db)
    
    // Test data
    id := types.NewBinaryUUID()
    expected := &models.Entity{
        UUID: id,
        Name: "Test Entity",
    }
    
    // Setup expectations
    mock.ExpectQuery("SELECT (.+) FROM entities").
        WithArgs(id).
        WillReturnRows(sqlmock.NewRows([]string{"uuid", "name"}).
            AddRow(id, "Test Entity"))
    
    // Execute
    result, err := repo.GetByID(context.Background(), id)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected.UUID, result.UUID)
    assert.Equal(t, expected.Name, result.Name)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### 3. Domain Model Tests

```go
func TestAvailability_GenerateSlots(t *testing.T) {
    tests := []struct {
        name           string
        availability   models.Availability
        date          time.Time
        expectedSlots  int
        expectedError error
    }{
        {
            name: "valid slots",
            availability: models.Availability{
                StartTime:    time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
                EndTime:      time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
                SlotDuration: 30,
            },
            date:          time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
            expectedSlots: 4,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            slots, err := tt.availability.GenerateSlots(tt.date)
            
            if tt.expectedError != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.expectedError, err)
            } else {
                assert.NoError(t, err)
                assert.Len(t, slots, tt.expectedSlots)
            }
        })
    }
}
```

## Mocking

### 1. Interface Mocking with GoMock

```go
//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks github.com/mukezhz/appointment-booking/domain/appointment Repository
type Repository interface {
    Create(context.Context, *models.Appointment) error
    GetByID(context.Context, types.BinaryUUID) (*models.Appointment, error)
}

// Test using generated mock
func TestService_WithMockGen(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockRepository(ctrl)
    service := NewService(mockRepo)
    
    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(nil)
    
    _, err := service.CreateAppointment(context.Background(), req)
    assert.NoError(t, err)
}
```

### 2. Manual Mocks with Testify

```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, entity *models.Entity) error {
    args := m.Called(ctx, entity)
    return args.Error(0)
}

// Test using manual mock
func TestService_WithTestify(t *testing.T) {
    mockRepo := new(MockRepository)
    service := NewService(mockRepo)
    
    mockRepo.On("Create", mock.Anything, mock.Anything).
        Return(nil)
    
    _, err := service.CreateAppointment(context.Background(), req)
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

## Testing Error Cases

### 1. Business Logic Errors

```go
func TestService_ValidationErrors(t *testing.T) {
    service := NewService(new(MockRepository))
    
    tests := []struct {
        name        string
        input       *CreateAppointmentRequest
        expectedErr error
    }{
        {
            name: "past date",
            input: &CreateAppointmentRequest{
                DateTime: time.Now().Add(-24 * time.Hour),
            },
            expectedErr: ErrPastDate,
        },
        {
            name: "invalid duration",
            input: &CreateAppointmentRequest{
                DateTime: time.Now().Add(24 * time.Hour),
                Duration: 10,
            },
            expectedErr: ErrInvalidDuration,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := service.CreateAppointment(context.Background(), tt.input)
            assert.Error(t, err)
            assert.Equal(t, tt.expectedErr, err)
        })
    }
}
```

### 2. Repository Errors

```go
func TestRepository_Errors(t *testing.T) {
    tests := []struct {
        name        string
        setupDB     func(*sqlmock.Sqlmock)
        expectedErr error
    }{
        {
            name: "duplicate key",
            setupDB: func(mock *sqlmock.Sqlmock) {
                (*mock).ExpectQuery("SELECT").
                    WillReturnError(&mysql.MySQLError{
                        Number: 1062,
                    })
            },
            expectedErr: ErrDuplicateKey,
        },
        {
            name: "foreign key violation",
            setupDB: func(mock *sqlmock.Sqlmock) {
                (*mock).ExpectQuery("SELECT").
                    WillReturnError(&mysql.MySQLError{
                        Number: 1452,
                    })
            },
            expectedErr: ErrForeignKeyViolation,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            db, mock := testutil.NewMockDB()
            repo := NewRepository(db)
            
            if tt.setupDB != nil {
                tt.setupDB(&mock)
            }
            
            // Execute
            _, err := repo.GetByID(context.Background(), types.NewBinaryUUID())
            
            // Assert
            assert.Error(t, err)
            assert.Equal(t, tt.expectedErr, err)
            assert.NoError(t, mock.ExpectationsWereMet())
        })
    }
}
```

## Best Practices

1. **Test Organization**
   - Group related test cases
   - Use clear, descriptive names
   - Follow consistent patterns
   - Keep tests focused

2. **Test Coverage**
   - Test happy paths
   - Test error cases
   - Test edge cases
   - Test validation rules

3. **Mock Usage**
   - Mock external dependencies
   - Verify mock expectations
   - Use appropriate mock types
   - Don't over-mock

4. **Assertions**
   - Be specific in assertions
   - Check error types
   - Verify all relevant fields
   - Use appropriate assertion functions
