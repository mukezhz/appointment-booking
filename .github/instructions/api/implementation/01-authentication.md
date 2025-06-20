# Authentication Guide

This guide explains how to implement and use authentication in our REST APIs.

## Authentication Middleware

### 1. Configuration

The authentication middleware is configured in `pkg/middlewares/auth_middleware.go`:

```go
type AuthMiddleware struct {
    logger  framework.Logger
    cognito *services.CognitoService
}

func NewAuthMiddleware(
    logger framework.Logger,
    cognito *services.CognitoService,
) *AuthMiddleware {
    return &AuthMiddleware{
        logger:  logger,
        cognito: cognito,
    }
}
```

### 2. JWT Token Validation

```go
func (m *AuthMiddleware) ValidateToken() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        token := m.extractToken(ctx)
        if token == "" {
            responses.HandleErrorWithStatus(ctx, m.logger,
                http.StatusUnauthorized, errors.New("no token provided"))
            ctx.Abort()
            return
        }

        claims, err := m.cognito.ValidateToken(ctx.Request.Context(), token)
        if err != nil {
            responses.HandleErrorWithStatus(ctx, m.logger,
                http.StatusUnauthorized, err)
            ctx.Abort()
            return
        }

        // Set user info in context
        ctx.Set(framework.UserIDKey, claims.Subject)
        ctx.Set(framework.UserRolesKey, claims.Roles)
        ctx.Next()
    }
}
```

## Using Authentication

### 1. Protect Routes

Apply authentication to routes in your `route.go`:

```go
func (r *Route) Register(router *gin.RouterGroup) {
    // Public routes
    router.POST("/auth/login", r.controller.Login)
    router.POST("/auth/register", r.controller.Register)

    // Protected routes
    authorized := router.Group("")
    authorized.Use(r.auth.ValidateToken())
    {
        authorized.POST("/availability", r.controller.CreateAvailability)
        authorized.GET("/availability", r.controller.ListAvailability)
    }
}
```

### 2. Access User Context

In your controllers, get the authenticated user's information:

```go
func (c *Controller) CreateAvailability(ctx *gin.Context) {
    // Get authenticated user ID
    userID := ctx.MustGet(framework.UserIDKey).(types.BinaryUUID)

    // Get user roles
    roles := ctx.MustGet(framework.UserRolesKey).([]string)

    // ... rest of the implementation
}
```

## Role-Based Access Control

### 1. Role Middleware

```go
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        userRoles := ctx.MustGet(framework.UserRolesKey).([]string)
        
        if !hasAnyRole(userRoles, roles) {
            responses.HandleErrorWithStatus(ctx, m.logger,
                http.StatusForbidden,
                errors.New("insufficient permissions"))
            ctx.Abort()
            return
        }
        
        ctx.Next()
    }
}

func hasAnyRole(userRoles []string, requiredRoles []string) bool {
    for _, required := range requiredRoles {
        for _, role := range userRoles {
            if role == required {
                return true
            }
        }
    }
    return false
}
```

### 2. Apply Role Requirements

```go
func (r *Route) Register(router *gin.RouterGroup) {
    authorized := router.Group("")
    authorized.Use(r.auth.ValidateToken())

    // Doctor routes
    doctorRoutes := authorized.Group("")
    doctorRoutes.Use(r.auth.RequireRole("doctor"))
    {
        doctorRoutes.POST("/availability", r.controller.CreateAvailability)
        doctorRoutes.GET("/appointments", r.controller.ListDoctorAppointments)
    }

    // Patient routes
    patientRoutes := authorized.Group("")
    patientRoutes.Use(r.auth.RequireRole("patient"))
    {
        patientRoutes.POST("/appointments", r.controller.BookAppointment)
        patientRoutes.GET("/appointments", r.controller.ListPatientAppointments)
    }
}
```

## Error Handling

### 1. Authentication Errors

```go
// Token missing
{
    "error": "No token provided",
    "code": "AUTH_REQUIRED"
}

// Invalid token
{
    "error": "Invalid or expired token",
    "code": "INVALID_TOKEN"
}
```

### 2. Authorization Errors

```go
// Insufficient permissions
{
    "error": "Insufficient permissions",
    "code": "FORBIDDEN",
    "details": {
        "required_role": "doctor",
        "user_role": "patient"
    }
}
```

## Security Best Practices

1. **Token Security**
   - Use HTTPS only
   - Store tokens securely
   - Implement token refresh
   - Set appropriate expiration

2. **Rate Limiting**
   ```go
   router.Use(middlewares.RateLimitMiddleware(
       redis,
       &RateLimitConfig{
           RequestsPerSecond: 10,
           BurstSize:        20,
       },
   ))
   ```

3. **CORS Configuration**
   ```go
   router.Use(middlewares.CORSMiddleware(
       &CORSConfig{
           AllowOrigins:     []string{"https://api.example.com"},
           AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
           AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
           ExposeHeaders:    []string{"Content-Length"},
           AllowCredentials: true,
           MaxAge:          12 * time.Hour,
       },
   ))
   ```

4. **Logging Security Events**
   ```go
   m.logger.Warn("authentication failed",
       "ip", ctx.ClientIP(),
       "path", ctx.Request.URL.Path,
       "error", err,
   )
   ```

5. **Request Sanitization**
   - Validate all input
   - Sanitize user data
   - Prevent injection attacks
   - Use prepared statements
