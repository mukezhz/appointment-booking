# Testing Strategy

This guide outlines our overall testing strategy and setup for the Go Clean Architecture project.

## Testing Framework and Tools

### Core Tools
- **Ginkgo**: BDD testing framework for Go
- **Testify**: Assertion library with support for mocks
- **TestContainers**: For database integration testing
- **apitest**: HTTP testing utility

### Additional Tools
- **GoMock**: For interface mocking
- **SQLMock**: For database query testing
- **httptest**: Go's built-in HTTP testing package

## Test Organization

Our tests follow the domain-driven structure of the application:

```
domain/
└── feature/
    ├── feature_suite_test.go  # Test suite setup
    ├── service_test.go        # Business logic tests
    ├── route_test.go          # API endpoint tests
    ├── repository_test.go     # Data access tests
    └── integration_test.go    # Integration tests
```

## Test Suite Setup

### 1. Base Test Suite

```go
package feature_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "clean-architecture/pkg/utils"
)

func TestFeature(t *testing.T) {
    utils.ChDir() // Ensures tests run from project root
    RegisterFailHandler(Fail)
    RunSpecs(t, "Feature Suite")
}

var t GinkgoTInterface
var _ = BeforeSuite(func() {
    t = GinkgoT()
})
```

### 2. Integration Test Suite

```go
var testDB *gorm.DB
var container *testcontainers.Container

var _ = BeforeSuite(func() {
    // Start test database container
    container, testDB = testutil.SetupTestDB()
})

var _ = AfterSuite(func() {
    // Cleanup
    if container != nil {
        container.Terminate(context.Background())
    }
})
```

## Test Environments

### 1. Local Development
- Uses local database for integration tests
- Fast feedback loop
- Suitable for TDD

### 2. CI Environment
- Uses containerized dependencies
- Runs full test suite
- Enforces coverage requirements

## Test Categories

1. **Unit Tests**
   - Test individual components
   - Use mocks for dependencies
   - Fast execution
   
2. **Integration Tests**
   - Test component interactions
   - Use test database
   - Test actual HTTP endpoints

3. **Contract Tests**
   - Verify API contracts
   - Test request/response formats
   - Validate error scenarios

## Test Coverage

### Coverage Requirements
- Unit tests: 80% coverage minimum
- Integration tests: Key flows covered
- API tests: All endpoints tested

### Coverage Reporting
```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

## Test Data Management

### 1. Fixtures
```go
var (
    TestUsers = []models.User{
        {
            UUID:  types.NewBinaryUUID(),
            Email: "test@example.com",
            Role:  "doctor",
        },
    }
)
```

### 2. Factories
```go
func NewTestAppointment(opts ...func(*models.Appointment)) *models.Appointment {
    appt := &models.Appointment{
        UUID:     types.NewBinaryUUID(),
        Status:   models.StatusPending,
        DateTime: time.Now().Add(24 * time.Hour),
    }
    
    for _, opt := range opts {
        opt(appt)
    }
    
    return appt
}
```

## Best Practices

1. **Test Organization**
   - Group related tests
   - Use descriptive names
   - Follow consistent patterns
   - Keep tests focused

2. **Test Independence**
   - Tests should not depend on each other
   - Clean up test data
   - Use fresh fixtures
   - Avoid global state

3. **Test Maintainability**
   - Use helper functions
   - Share test utilities
   - Document complex setups
   - Keep assertions clear

4. **Continuous Integration**
   - Run tests on every PR
   - Maintain coverage thresholds
   - Use test reports
   - Automate test runs

## Example: Complete Test Setup

```go
package appointment_test

import (
    "testing"
    "context"
    
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/stretchr/testify/mock"
    
    "clean-architecture/domain/appointment"
    "clean-architecture/pkg/testutil"
)

// Test suite setup
func TestAppointment(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Appointment Suite")
}

var (
    testDB     *gorm.DB
    container  *testcontainers.Container
    ctx        context.Context
    repo       *appointment.Repository
    service    *appointment.Service
    controller *appointment.Controller
)

var _ = BeforeSuite(func() {
    ctx = context.Background()
    
    // Setup test database
    container, testDB = testutil.SetupTestDB()
    
    // Initialize components
    repo = appointment.NewRepository(testDB)
    service = appointment.NewService(repo)
    controller = appointment.NewController(service)
})

var _ = AfterSuite(func() {
    if container != nil {
        container.Terminate(ctx)
    }
})

var _ = BeforeEach(func() {
    // Clean database before each test
    testutil.CleanDB(testDB)
    
    // Seed test data
    testutil.SeedTestData(testDB)
})
```
