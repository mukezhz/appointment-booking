package middlewares

import (
	"clean-architecture/pkg/framework"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// MockAuthMiddleware is a mock implementation of AuthMiddleware
type MockAuthMiddleware struct {
	logger framework.Logger
}

func NewMockAuthMiddleware(logger framework.Logger) *MockAuthMiddleware {
	return &MockAuthMiddleware{
		logger: logger,
	}
}

// HandleAuthWithRole returns a gin middleware that mocks authentication with role checking
func (m *MockAuthMiddleware) HandleAuthWithRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// In a mock implementation, we just allow all requests to pass through
		// You can customize this to simulate different auth scenarios for testing
		authHeader := ctx.GetHeader("Authorization") // Simulate reading an Authorization header
		m.logger.Debug("MockAuthMiddleware: Authorization header received: ", authHeader)
		ctx.Set(framework.UID, "mock-user-id")
		ctx.Set(framework.Role, "mock-role") // Set the first role for testing
		if len(roles) > 0 {
			ctx.Set(framework.Role, roles[0]) // Set the first role from the provided roles
		}

		userID := ctx.GetHeader(framework.UIDHeader)
		m.logger.Debug("MockAuthMiddleware: User ID set: ", userID)
		if userID != "" {
			ctx.Set(framework.UID, userID)
		}
		if strings.ToLower(strings.TrimSpace(authHeader)) != "bearer mock" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		m.logger.Debug("MockAuthMiddleware: User authenticated with mock credentials")
		ctx.Next()
	}
}

// HandleAuth returns a gin middleware that mocks basic authentication
func (m *MockAuthMiddleware) HandleAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple mock implementation that sets a user ID
		c.Set(framework.UID, "mock-user-id")
		c.Next()
	}
}

// GetUserID is a helper method to extract user ID from the context
func (m *MockAuthMiddleware) GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get(framework.UID)
	if !exists {
		return "", nil
	}
	return userID.(string), nil
}

// GetUserRole is a helper method to extract user role from the context
func (m *MockAuthMiddleware) GetUserRole(c *gin.Context) (string, error) {
	userRole, exists := c.Get(framework.Role)
	if !exists {
		return "", nil
	}
	return userRole.(string), nil
}
