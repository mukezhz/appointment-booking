# API Testing Guide

This guide explains how to write and maintain tests for our REST APIs.

## Test Types

1. **Unit Tests**
   - Test individual components in isolation
   - Use mocks for dependencies
   - Focus on business logic

2. **Integration Tests**
   - Test component interactions
   - Use test database
   - Test actual HTTP endpoints

3. **API Tests**
   - End-to-end testing
   - Test complete request flows
   - Validate response formats

## Testing Structure

### 1. Test Setup

```go
// test_setup.go
package tests

import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type TestSuite struct {
    suite.Suite
    app    *bootstrap.App
    router *gin.Engine
    db     *gorm.DB
}

func (s *TestSuite) SetupSuite() {
    // Initialize test environment
    s.app = bootstrap.NewTestApp()
    s.router = s.app.Router
    s.db = s.app.DB
}

func (s *TestSuite) TearDownSuite() {
    // Cleanup
    s.db.Close()
}

func (s *TestSuite) SetupTest() {
    // Run migrations
    // Seed test data
}

func (s *TestSuite) TearDownTest() {
    // Clean test data
}
```

### 2. Unit Tests

```go
// service_test.go
func TestAvailabilityService(t *testing.T) {
    tests := []struct {
        name        string
        input       *CreateAvailabilityRequest
        mockSetup   func(*mocks.Repository)
        wantErr     bool
        expectedErr error
    }{
        {
            name: "success",
            input: &CreateAvailabilityRequest{
                StartTime:    time.Now(),
                EndTime:     time.Now().Add(1 * time.Hour),
                SlotDuration: 30,
            },
            mockSetup: func(r *mocks.Repository) {
                r.On("Create", mock.Anything, mock.Anything).Return(nil)
            },
            wantErr: false,
        },
        {
            name: "invalid duration",
            input: &CreateAvailabilityRequest{
                SlotDuration: 10, // Too short
            },
            wantErr:     true,
            expectedErr: ErrInvalidSlotDuration,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            mockRepo := new(mocks.Repository)
            if tt.mockSetup != nil {
                tt.mockSetup(mockRepo)
            }

            service := NewService(mockRepo)

            // Test
            _, err := service.CreateAvailability(context.Background(), tt.input)

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

### 3. Integration Tests

```go
// route_test.go
func (s *TestSuite) TestCreateAvailability() {
    // Setup
    token := s.generateTestToken("doctor")
    
    tests := []struct {
        name         string
        request      CreateAvailabilityRequest
        setupDB      func()
        wantStatus   int
        wantResponse responses.DetailResponseType[AvailabilityDTO]
    }{
        {
            name: "success",
            request: CreateAvailabilityRequest{
                StartTime:    time.Now(),
                EndTime:     time.Now().Add(1 * time.Hour),
                SlotDuration: 30,
            },
            wantStatus: http.StatusCreated,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        s.Run(tt.name, func() {
            if tt.setupDB != nil {
                tt.setupDB()
            }

            // Make request
            w := httptest.NewRecorder()
            reqBody, _ := json.Marshal(tt.request)
            req, _ := http.NewRequest("POST", "/api/v1/availability", bytes.NewReader(reqBody))
            req.Header.Set("Authorization", "Bearer "+token)

            // Serve request
            s.router.ServeHTTP(w, req)

            // Assert response
            s.Equal(tt.wantStatus, w.Code)
            if tt.wantStatus == http.StatusCreated {
                var response responses.DetailResponseType[AvailabilityDTO]
                s.NoError(json.NewDecoder(w.Body).Decode(&response))
                s.Equal(tt.wantResponse, response)
            }
        })
    }
}
```

### 4. API Tests

```go
// api_test.go
func (s *TestSuite) TestAppointmentBookingFlow() {
    // 1. Setup test data
    doctor := s.createTestDoctor()
    patient := s.createTestPatient()
    availability := s.createTestAvailability(doctor.ID)

    // 2. Test listing available slots
    slots := s.listAvailableSlots(doctor.ID)
    s.NotEmpty(slots)

    // 3. Test booking appointment
    appointment := s.bookAppointment(patient.ID, slots[0].ID)
    s.NotNil(appointment)

    // 4. Test retrieving appointment
    retrieved := s.getAppointment(appointment.ID)
    s.Equal(appointment.ID, retrieved.ID)

    // 5. Test canceling appointment
    s.cancelAppointment(appointment.ID)

    // 6. Verify cancellation
    cancelled := s.getAppointment(appointment.ID)
    s.Equal("cancelled", cancelled.Status)
}
```

## Mocking

### 1. Repository Mocks

```go
// mocks/repository.go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, entity *models.Entity) error {
    args := m.Called(ctx, entity)
    return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id types.BinaryUUID) (*models.Entity, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*models.Entity), args.Error(1)
}
```

### 2. Service Mocks

```go
// mocks/service.go
type MockService struct {
    mock.Mock
}

func (m *MockService) CreateEntity(ctx context.Context, req *CreateEntityRequest) (*models.Entity, error) {
    args := m.Called(ctx, req)
    return args.Get(0).(*models.Entity), args.Error(1)
}
```

## Test Utilities

### 1. Test Helpers

```go
// testutil/helpers.go
func GenerateTestToken(role string) string {
    // Create test JWT token
}

func CreateTestUser(role string) *models.User {
    // Create test user in DB
}

func CleanTestData(db *gorm.DB) {
    // Clean up test data
}
```

### 2. Request Helpers

```go
// testutil/request.go
func MakeRequest(router *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
    w := httptest.NewRecorder()
    reqBody, _ := json.Marshal(body)
    req, _ := http.NewRequest(method, path, bytes.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    if token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }
    router.ServeHTTP(w, req)
    return w
}
```

## Test Data

### 1. Fixtures

```go
// testutil/fixtures.go
var (
    TestUsers = []models.User{
        {
            UUID:     types.NewBinaryUUID(),
            Email:    "test@example.com",
            Role:     "doctor",
            Status:   "active",
        },
        // ... more test users
    }

    TestAvailability = []models.Availability{
        {
            UUID:      types.NewBinaryUUID(),
            DoctorID:  TestUsers[0].UUID,
            StartTime: time.Now(),
            EndTime:   time.Now().Add(1 * time.Hour),
        },
        // ... more test data
    }
)
```

### 2. Seeds

```go
// testutil/seeds.go
func SeedTestData(db *gorm.DB) error {
    for _, user := range TestUsers {
        if err := db.Create(&user).Error; err != nil {
            return err
        }
    }
    // ... seed other test data
    return nil
}
```

## Best Practices

1. **Test Organization**
   - Group related tests together
   - Use descriptive test names
   - Follow Given-When-Then pattern
   - Keep tests focused and simple

2. **Test Coverage**
   - Test happy paths
   - Test error cases
   - Test edge cases
   - Test validation rules

3. **Clean Test Data**
   - Use transactions for isolation
   - Clean up after tests
   - Don't rely on test order
   - Use fresh data for each test

4. **Assertions**
   - Use appropriate assertions
   - Check specific fields
   - Validate error types
   - Verify side effects

5. **Maintainability**
   - Keep tests DRY
   - Use helper functions
   - Document complex tests
   - Update tests with code changes
