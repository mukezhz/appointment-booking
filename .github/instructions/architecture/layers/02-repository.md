# Repository Layer

## Overview

The repository layer handles all database operations and data persistence. It's responsible for implementing the data access patterns and managing the interaction with the database.

## Key Responsibilities

1. **Database Operations**
   - CRUD operations
   - Complex queries
   - Transaction management
   - Error mapping

2. **Data Access Patterns**
   - Pagination
   - Filtering
   - Sorting
   - Eager loading

## Implementation Guide

### 1. Basic Repository Structure

```go
// domain/<feature>/repository.go
type Repository struct {
    db     infrastructure.Database  // GORM database instance
    logger framework.Logger        // Structured logger
}

// Constructor for dependency injection
func NewRepository(
    db infrastructure.Database,
    logger framework.Logger,
) *Repository {
    return &Repository{
        db:     db,
        logger: logger,
    }
}

// Error maps for database error handling
var RepositoryErrMap = map[error]bool{
    ErrEntityNotFound: true,
    ErrDuplicateKey:  true,
}
```

### 2. CRUD Operations

```go
// Create
func (r *Repository) Create(ctx context.Context, entity *models.Entity) error {
    if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
        r.logger.Error("failed to create entity: ", err)
        return common.HandleDBError(err, RepositoryErrMap)
    }
    return nil
}

// Read
func (r *Repository) GetByID(ctx context.Context, id types.BinaryUUID) (*models.Entity, error) {
    var entity models.Entity
    if err := r.db.WithContext(ctx).First(&entity, "uuid = ?", id).Error; err != nil {
        return nil, common.HandleDBError(err, RepositoryErrMap)
    }
    return &entity, nil
}

// Update
func (r *Repository) Update(ctx context.Context, entity *models.Entity) error {
    if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
        return common.HandleDBError(err, RepositoryErrMap)
    }
    return nil
}

// Delete
func (r *Repository) Delete(ctx context.Context, id types.BinaryUUID) error {
    if err := r.db.WithContext(ctx).Delete(&models.Entity{}, "uuid = ?", id).Error; err != nil {
        return common.HandleDBError(err, RepositoryErrMap)
    }
    return nil
}
```

### 3. List Operations with Filtering

```go
type ListParams struct {
    Status    []string
    FromDate  *time.Time
    ToDate    *time.Time
    Page      int
    PerPage   int
}

func (r *Repository) List(ctx context.Context, params ListParams) ([]models.Entity, int64, error) {
    var entities []models.Entity
    var total int64

    query := r.db.WithContext(ctx).Model(&models.Entity{})

    // Apply filters
    if len(params.Status) > 0 {
        query = query.Where("status IN ?", params.Status)
    }
    if params.FromDate != nil {
        query = query.Where("created_at >= ?", params.FromDate)
    }
    if params.ToDate != nil {
        query = query.Where("created_at <= ?", params.ToDate)
    }

    // Get total count
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, common.HandleDBError(err, RepositoryErrMap)
    }

    // Apply pagination
    offset := (params.Page - 1) * params.PerPage
    if err := query.Offset(offset).Limit(params.PerPage).Find(&entities).Error; err != nil {
        return nil, 0, common.HandleDBError(err, RepositoryErrMap)
    }

    return entities, total, nil
}
```

### 4. Transaction Handling

```go
func (r *Repository) TransactionalUpdate(ctx context.Context, fn func(tx *gorm.DB) error) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := fn(tx); err != nil {
            return common.HandleDBError(err, RepositoryErrMap)
        }
        return nil
    })
}

// Usage example
func (r *Repository) UpdateWithRelations(ctx context.Context, entity *models.Entity) error {
    return r.TransactionalUpdate(ctx, func(tx *gorm.DB) error {
        // Update main entity
        if err := tx.Save(entity).Error; err != nil {
            return err
        }

        // Update related entities
        if err := tx.Save(&entity.Relations).Error; err != nil {
            return err
        }

        return nil
    })
}
```

## Best Practices

1. **Error Handling**
   ```go
   // Always use common.HandleDBError with appropriate error maps
   if err != nil {
       return common.HandleDBError(err, RepositoryErrMap)
   }
   ```

2. **Context Usage**
   ```go
   // Always use WithContext
   r.db.WithContext(ctx).Find(&result)
   ```

3. **Logging**
   ```go
   // Log operations with context
   r.logger.Info("fetching entity by ID: ", id.String())
   if err != nil {
       r.logger.Error("failed to fetch entity: ", err)
   }
   ```

4. **Query Building**
   ```go
   // Build queries incrementally
   query := r.db.WithContext(ctx)
   if filter.Status != "" {
       query = query.Where("status = ?", filter.Status)
   }
   if filter.Type != "" {
       query = query.Where("type = ?", filter.Type)
   }
   ```

5. **Eager Loading**
   ```go
   // Load relationships efficiently
   if err := r.db.WithContext(ctx).
       Preload("Relations").
       Preload("Details").
       First(&entity, "uuid = ?", id).
       Error; err != nil {
       return nil, common.HandleDBError(err, RepositoryErrMap)
   }
   ```

## Testing Repositories

```go
Describe("Repository", func() {
    var (
        repo *Repository
        ctx  context.Context
    )

    BeforeEach(func() {
        // Setup test database
        db := setupTestDB()
        logger := framework.NewLogger()
        repo = NewRepository(db, logger)
        ctx = context.Background()
    })

    Describe("Create", func() {
        It("should create entity successfully", func() {
            entity := &models.Entity{
                Name: "Test Entity",
            }
            err := repo.Create(ctx, entity)
            Expect(err).To(BeNil())
            Expect(entity.ID).NotTo(BeZero())
        })

        It("should handle duplicate key error", func() {
            // Test duplicate creation
            err := repo.Create(ctx, duplicateEntity)
            Expect(err).To(Equal(ErrDuplicateKey))
        })
    })

    // Add more test cases
})
