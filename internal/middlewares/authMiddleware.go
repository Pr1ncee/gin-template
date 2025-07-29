/*
Package middlewares implements different kinds of middleware that are used in the application.

Specifically, this file implements authentication Middleware via JWT tokens
and also describes an authorization mechanism that validates a User's role before allowing the request to proceed.
*/
package middlewares

import (
	db "GinBox/internal/postgresql"
	"GinBox/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strings"
)

// AuthMiddleware represents authentication middleware via JWT tokens.
type AuthMiddleware struct {
	authService services.IAuthService
}

// NewAuthMiddleware initializes a new AuthMiddleware object.
func NewAuthMiddleware(authService services.IAuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// Handle takes the incoming request and validates token in the header.
func (a *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth header"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := a.authService.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user", claims.UserClaims)
		c.Next()
	}
}

// RequireRole represents a method that ensures a User has sufficient permission
// based on his role before processing the request.
func RequireRole(roles []db.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User not found in context"})
			return
		}
		user := claims.(services.UserClaims)
		if !slices.Contains(roles, user.Role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		c.Next()
	}
}
