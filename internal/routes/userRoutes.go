package routes

import (
	"GinBox/internal/handlers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// UserRoutes sets up all user-related routes
func UserRoutes(
	router *gin.RouterGroup,
	userHandler *handlers.UserHandler,
	logger *log.Logger,
) {
	users := router.Group("/users")

	// Apply middleware
	//users.Use(middleware.RequestLogger(logger))
	//users.Use(middleware.CORS())
	//users.Use(middleware.RateLimiter(100, time.Minute)) // 100 requests per minute

	// Routes that require authentication
	authenticated := users.Group("")
	//authenticated.Use(middleware.AuthRequired()) // You'll need to implement this
	{
		authenticated.POST("", userHandler.CreateUser)        // POST /api/v1/users
		authenticated.GET("", userHandler.ListUsers)          // GET /api/v1/users?page=1&limit=10&role=Admin&q=test@gmail.com
		authenticated.GET("/:id", userHandler.GetUserByID)    // GET /api/v1/users/123
		authenticated.PUT("/:id", userHandler.UpdateUser)     // PUT /api/v1/users/123
		authenticated.DELETE("/:id", userHandler.DeleteUser)  // DELETE /api/v1/users/123
		authenticated.GET("/export", userHandler.ExportUsers) // GET /api/v1/users/export?page=1&limit=10&role=Admin&q=test@gmail.com
	}
}
