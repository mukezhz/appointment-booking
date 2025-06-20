# Dependency Injection Guide

This guide explains how we use dependency injection in our clean architecture using `uber-go/fx`.

## Key Principles

- Constructor injection
- Interface-based dependencies
- Clear dependency graph
- Lazy initialization

## Module Structure

Each feature module is defined using `fx.Module`:

```go
var Module = fx.Module(
    "feature",
    fx.Provide(
        NewRepository,
        NewService,
        NewController,
        NewRoute,
    ),
    fx.Invoke(RegisterRoutes),
)
```

## Component Registration

### 1. Constructors

Each layer has a constructor that declares its dependencies:

```go
// Repository constructor
func NewRepository(
    db infrastructure.Database,
    logger framework.Logger,
) *Repository {
    return &Repository{
        db:     db,
        logger: logger,
    }
}

// Service constructor
func NewService(
    repo *Repository,
    logger framework.Logger,
) *Service {
    return &Service{
        repo:   repo,
        logger: logger,
    }
}

// Controller constructor
func NewController(
    service *Service,
    logger framework.Logger,
) *Controller {
    return &Controller{
        service: service,
        logger: logger,
    }
}
```

### 2. Registering Routes

Routes are registered using `fx.Invoke`:

```go
func RegisterRoutes(r *Route) {
    // Route registration logic
}
```

## Module Organization

### 1. Application Module (`bootstrap/modules.go`)

```go
var Module = fx.Options(
    infrastructure.Module,
    framework.Module,
    services.Module,
    feature1.Module,
    feature2.Module,
)
```

### 2. Feature Module (`domain/feature/module.go`)

```go
var Module = fx.Module(
    "feature",
    fx.Provide(
        NewRepository,
        NewService,
        NewController,
        NewRoute,
    ),
    fx.Invoke(RegisterRoutes),
)
```

## Best Practices

1. **Interface Dependencies**: Prefer interfaces over concrete types for flexibility
    ```go
    type FeatureService interface {
        Create(context.Context, *models.Feature) error
        Get(context.Context, types.BinaryUUID) (*models.Feature, error)
    }
    ```

2. **Group Related Dependencies**: Use options struct for multiple dependencies
    ```go
    type ServiceOptions struct {
        Repository  *Repository
        Logger     framework.Logger
        Config     *Config
    }

    func NewService(opt ServiceOptions) *Service {
        return &Service{
            repo:   opt.Repository,
            logger: opt.Logger,
            config: opt.Config,
        }
    }
    ```

3. **Lifecycle Management**: Use `fx.Lifecycle` for start/stop hooks
    ```go
    func NewService(lc fx.Lifecycle, repo *Repository) *Service {
        svc := &Service{repo: repo}
        lc.Append(fx.Hook{
            OnStart: func(ctx context.Context) error {
                return svc.Initialize()
            },
            OnStop: func(ctx context.Context) error {
                return svc.Cleanup()
            },
        })
        return svc
    }
    ```

4. **Error Handling**: Handle initialization errors in constructors
    ```go
    func NewService(repo *Repository) (*Service, error) {
        if repo == nil {
            return nil, errors.New("repository is required")
        }
        return &Service{repo: repo}, nil
    }
    ```

## Testing with Dependencies

1. **Mock Dependencies**: Use interfaces for easy mocking
    ```go
    type mockRepository struct {
        mock.Mock
    }

    func (m *mockRepository) Create(ctx context.Context, f *models.Feature) error {
        args := m.Called(ctx, f)
        return args.Error(0)
    }
    ```

2. **Unit Testing**: Create test doubles
    ```go
    func TestService_Create(t *testing.T) {
        mockRepo := &mockRepository{}
        service := NewService(mockRepo)
        
        mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
        
        err := service.Create(context.Background(), &models.Feature{})
        assert.NoError(t, err)
    }
    ```

## Common Issues and Solutions

1. **Circular Dependencies**
   - Use interfaces to break cycles
   - Consider moving shared code to a common package

2. **Order of Initialization**
   - Use `fx.Invoke` for setup that must run after all providers
   - Add dependencies to constructors to enforce order

3. **Testing**
   - Create test-specific modules
   - Use dependency injection in tests for better isolation
