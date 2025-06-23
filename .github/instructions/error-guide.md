# Error Handling

## Core Error Types
1. Base Errors (`pkg/errorz`)
   - APIError with StatusCode & Message
   - Standard HTTP errors (400,401,403,404,409,422,500,503)

## Usage
```go
// Create errors
NewNotFoundError("item not found")
NewBadRequestError("invalid input")

// Handle errors
HandleError(ctx, logger, err)
HandleDBError(err, domainErrs)
```

## Best Practices
1. Use typed errors (`ErrNotFound`, etc.)
2. Map DB errors to domain errors
3. Consistent error responses
4. Log all errors
