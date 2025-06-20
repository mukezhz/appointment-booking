# Integration Testing Guide

This guide covers integration testing patterns and best practices for our Go Clean Architecture project.

## Integration Test Setup

### 1. Test Container Setup

```go
// testutil/db.go
func SetupTestDB() (*testcontainers.Container, *gorm.DB) {
    ctx := context.Background()
    
    // Create MySQL container
    req := testcontainers.ContainerRequest{
        Image:        "mysql:8",
        ExposedPorts: []string{"3306/tcp"},
        Env: map[string]string{
            "MYSQL_ROOT_PASSWORD": "test",
            "MYSQL_DATABASE":     "test_db",
        },
        WaitingFor: wait.ForLog("port: 3306  MySQL Community Server"),
    }
    
    container, err := testcontainers.GenericContainer(ctx, req)
    if err != nil {
        panic(err)
    }
    
    // Get connection details
    port, _ := container.MappedPort(ctx, "3306")
    host, _ := container.Host(ctx)
    
    // Connect to database
    dsn := fmt.Sprintf("root:test@tcp(%s:%s)/test_db?parseTime=true", 
        host, port.Port())
    db, err := gorm.Open(mysql.Open(dsn))
    if err != nil {
        panic(err)
    }
    
    // Run migrations
    if err := RunMigrations(db); err != nil {
        panic(err)
    }
    
    return container, db
}
```

### 2. Test Suite Setup

```go
// feature/integration_test.go
var (
    container *testcontainers.Container
    testDB    *gorm.DB
    app       *bootstrap.App
    router    *gin.Engine
)

var _ = BeforeSuite(func() {
    // Setup test database
    container, testDB = testutil.SetupTestDB()
    
    // Initialize application
    app = bootstrap.NewTestApp(testDB)
    router = app.Router
})

var _ = AfterSuite(func() {
    if container != nil {
        container.Terminate(context.Background())
    }
})

var _ = BeforeEach(func() {
    // Clean database
    testutil.CleanDB(testDB)
    
    // Seed test data
    testutil.SeedTestData(testDB)
})
```

## API Integration Tests

### 1. Testing HTTP Endpoints

```go
func TestAppointmentAPI(t *testing.T) {
    tests := []struct {
        name         string
        method       string
        path         string
        body         interface{}
        setupAuth    func() string
        setupDB      func()
        wantStatus   int
        wantResponse interface{}
    }{
        {
            name:   "create appointment",
            method: "POST",
            path:   "/api/v1/appointments",
            body: CreateAppointmentRequest{
                DoctorID:  "test-doctor",
                DateTime:  time.Now().Add(24 * time.Hour),
                Duration: 30,
            },
            setupAuth: func() string {
                return testutil.GenerateTestToken("patient")
            },
            setupDB: func() {
                testutil.CreateTestDoctor(testDB, "test-doctor")
            },
            wantStatus: http.StatusCreated,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.setupDB != nil {
                tt.setupDB()
            }
            
            token := ""
            if tt.setupAuth != nil {
                token = tt.setupAuth()
            }
            
            // Make request
            w := httptest.NewRecorder()
            reqBody, _ := json.Marshal(tt.body)
            req, _ := http.NewRequest(tt.method, tt.path, bytes.NewReader(reqBody))
            
            if token != "" {
                req.Header.Set("Authorization", "Bearer "+token)
            }
            
            // Serve request
            router.ServeHTTP(w, req)
            
            // Assert response
            assert.Equal(t, tt.wantStatus, w.Code)
            if tt.wantResponse != nil {
                var response interface{}
                assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))
                assert.Equal(t, tt.wantResponse, response)
            }
        })
    }
}
```

### 2. Testing Complete Flows

```go
func TestAppointmentBookingFlow(t *testing.T) {
    // 1. Setup test data
    doctor := testutil.CreateTestDoctor(testDB)
    patient := testutil.CreateTestPatient(testDB)
    availability := testutil.CreateTestAvailability(testDB, doctor.ID)
    
    patientToken := testutil.GenerateTestToken("patient", patient.ID)
    
    // 2. List available slots
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET",
        fmt.Sprintf("/api/v1/doctors/%s/availability", doctor.ID),
        nil)
    req.Header.Set("Authorization", "Bearer "+patientToken)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    var slotsResponse responses.ListResponseType[AvailabilityDTO]
    assert.NoError(t, json.NewDecoder(w.Body).Decode(&slotsResponse))
    assert.NotEmpty(t, slotsResponse.Items)
    
    // 3. Book appointment
    slot := slotsResponse.Items[0]
    bookingReq := CreateAppointmentRequest{
        DoctorID:     doctor.ID.String(),
        DateTime:     slot.StartTime,
        Duration:     30,
        Description:  "Test appointment",
    }
    
    w = httptest.NewRecorder()
    reqBody, _ := json.Marshal(bookingReq)
    req, _ = http.NewRequest("POST", "/api/v1/appointments",
        bytes.NewReader(reqBody))
    req.Header.Set("Authorization", "Bearer "+patientToken)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusCreated, w.Code)
    var appointmentResponse responses.DetailResponseType[AppointmentDTO]
    assert.NoError(t, json.NewDecoder(w.Body).Decode(&appointmentResponse))
    
    // 4. Verify appointment details
    appointmentID := appointmentResponse.Item.UUID
    
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET",
        fmt.Sprintf("/api/v1/appointments/%s", appointmentID),
        nil)
    req.Header.Set("Authorization", "Bearer "+patientToken)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    var getResponse responses.DetailResponseType[AppointmentDTO]
    assert.NoError(t, json.NewDecoder(w.Body).Decode(&getResponse))
    assert.Equal(t, appointmentID, getResponse.Item.UUID)
}
```

## Repository Integration Tests

```go
func TestAppointmentRepository_Integration(t *testing.T) {
    // Setup
    repo := appointment.NewRepository(testDB)
    ctx := context.Background()
    
    // Test data
    doctor := testutil.CreateTestDoctor(testDB)
    patient := testutil.CreateTestPatient(testDB)
    
    t.Run("create and retrieve", func(t *testing.T) {
        // Create appointment
        appt := &models.Appointment{
            UUID:        types.NewBinaryUUID(),
            DoctorID:    doctor.ID,
            PatientID:   patient.ID,
            DateTime:    time.Now().Add(24 * time.Hour),
            Duration:    30,
            Status:      models.StatusPending,
        }
        
        err := repo.Create(ctx, appt)
        assert.NoError(t, err)
        
        // Retrieve appointment
        retrieved, err := repo.GetByID(ctx, appt.UUID)
        assert.NoError(t, err)
        assert.Equal(t, appt.UUID, retrieved.UUID)
        assert.Equal(t, appt.DoctorID, retrieved.DoctorID)
        assert.Equal(t, appt.Status, retrieved.Status)
    })
    
    t.Run("list with filters", func(t *testing.T) {
        // Create test appointments
        testutil.CreateTestAppointments(testDB, doctor.ID, patient.ID)
        
        // Test filtering
        filters := &appointment.ListFilters{
            DoctorID: doctor.ID,
            Status:   models.StatusPending,
            FromDate: time.Now(),
        }
        
        appointments, total, err := repo.List(ctx, filters, &types.Pagination{
            Page:    1,
            PerPage: 10,
        })
        
        assert.NoError(t, err)
        assert.NotEmpty(t, appointments)
        assert.Greater(t, total, int64(0))
        
        // Verify filter application
        for _, appt := range appointments {
            assert.Equal(t, doctor.ID, appt.DoctorID)
            assert.Equal(t, models.StatusPending, appt.Status)
        }
    })
}
```

## Testing External Services

### 1. Email Service Integration

```go
func TestEmailService_Integration(t *testing.T) {
    // Skip in CI
    if os.Getenv("CI") != "" {
        t.Skip("Skipping in CI environment")
    }
    
    // Setup
    config := &services.EmailConfig{
        FromAddress: "test@example.com",
        SMTPHost:    "localhost",
        SMTPPort:    "1025", // MailHog
    }
    emailService := services.NewEmailService(config)
    
    t.Run("send appointment confirmation", func(t *testing.T) {
        email := &services.Email{
            To:      "patient@example.com",
            Subject: "Appointment Confirmation",
            Body:    "Your appointment has been confirmed",
        }
        
        err := emailService.Send(context.Background(), email)
        assert.NoError(t, err)
        
        // Verify in MailHog API
        messages, err := testutil.GetMailHogMessages()
        assert.NoError(t, err)
        assert.NotEmpty(t, messages)
        
        latest := messages[len(messages)-1]
        assert.Equal(t, email.To, latest.To)
        assert.Equal(t, email.Subject, latest.Subject)
    })
}
```

### 2. AWS Service Integration

```go
func TestAWSService_Integration(t *testing.T) {
    // Use LocalStack for AWS testing
    config := &services.AWSConfig{
        Endpoint: "http://localhost:4566",
        Region:   "us-east-1",
    }
    awsService := services.NewAWSService(config)
    
    t.Run("s3 operations", func(t *testing.T) {
        // Upload file
        content := []byte("test content")
        key := "test.txt"
        
        err := awsService.UploadFile(context.Background(),
            "test-bucket", key, content)
        assert.NoError(t, err)
        
        // Download and verify
        downloaded, err := awsService.GetFile(context.Background(),
            "test-bucket", key)
        assert.NoError(t, err)
        assert.Equal(t, content, downloaded)
    })
}
```

## Best Practices

1. **Test Independence**
   - Clean database between tests
   - Use fresh test data
   - Avoid test interdependence
   - Handle cleanup properly

2. **Test Data Management**
   - Use factories for test data
   - Share common setup
   - Make data relationships clear
   - Clean up test artifacts

3. **Error Scenarios**
   - Test API error responses
   - Test database constraints
   - Test concurrent operations
   - Test timeout scenarios

4. **Performance**
   - Use parallel tests when possible
   - Minimize container restarts
   - Reuse connections
   - Clean up resources

5. **Environment Handling**
   - Support different environments
   - Use environment variables
   - Skip tests appropriately
   - Document prerequisites
