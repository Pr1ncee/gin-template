package routes

import (
	"GinBox/internal/handlers"
	"GinBox/internal/middlewares"
	db "GinBox/internal/postgresql"
	"GinBox/internal/services"
	"github.com/gin-gonic/gin"
)

// UserRoutes sets up all user-related routes
func UserRoutes(
	router *gin.RouterGroup,
	userHandler *handlers.UserHandler,
	authService services.IAuthService,
) {
	users := router.Group("/users")

	users.POST("/login", userHandler.Login)
	users.POST("/refresh", userHandler.RefreshToken)

	// Routes that require authentication
	authenticated := users.Group("")
	authenticated.Use(middlewares.AuthMiddleware(authService))
	{
		authenticated.POST("", middlewares.RequireRole([]db.UserRole{db.UserRoleAdmin}), userHandler.CreateUser)                         // POST /api/v1/users
		authenticated.GET("", middlewares.RequireRole([]db.UserRole{db.UserRoleUser, db.UserRoleAdmin}), userHandler.ListUsers)          // GET /api/v1/users?page=1&limit=10&role=Admin&q=test@gmail.com
		authenticated.GET("/:id", middlewares.RequireRole([]db.UserRole{db.UserRoleUser, db.UserRoleAdmin}), userHandler.GetUserByID)    // GET /api/v1/users/123
		authenticated.PUT("/:id", middlewares.RequireRole([]db.UserRole{db.UserRoleAdmin}), userHandler.UpdateUser)                      // PUT /api/v1/users/123
		authenticated.DELETE("/:id", middlewares.RequireRole([]db.UserRole{db.UserRoleAdmin}), userHandler.DeleteUser)                   // DELETE /api/v1/users/123
		authenticated.GET("/export", middlewares.RequireRole([]db.UserRole{db.UserRoleUser, db.UserRoleAdmin}), userHandler.ExportUsers) // GET /api/v1/users/export?page=1&limit=10&role=Admin&q=test@gmail.com
	}
}
