package middlewares

import (
	db "GinBox/internal/postgresql"
	"GinBox/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strings"
)

func AuthMiddleware(authService services.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid auth header"})
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := authService.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user", claims.UserClaims)
		c.Next()
	}
}

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
