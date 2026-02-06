package http

import (
	"errors"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

type AuthMiddleware struct {
	tokenService TokenService
}

type TokenService interface {
	Validate(token string) (string, error)
}

func NewAuthMiddleware(ts TokenService) *AuthMiddleware {
	return &AuthMiddleware{tokenService: ts}
}

// Recovery returns a middleware that recovers from panics and logs the error
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				log.Printf("[PANIC RECOVERED] %v\n%s", err, debug.Stack())

				// Send generic error response to client
				SendError(c, http.StatusInternalServerError, ErrCodeInternal, "Something went wrong")
				c.Abort()
			}
		}()
		c.Next()
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		userID, err := m.tokenService.Validate(token)
		if err != nil {
			SendError(c, http.StatusUnauthorized, ErrCodeUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user ID in context
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

// GetUserID extracts user ID from gin context
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", errors.New("user ID not found in context")
	}
	return userID.(string), nil
}
