package v1

import (
	"GinBox/config"
	"GinBox/internal/handlers"
	"GinBox/internal/middlewares"
	"GinBox/internal/repositories"
	"GinBox/internal/routes"
	"GinBox/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApiV1 struct {
	healthCheckHandler *handlers.HealthCheckHandler
	authService        services.IAuthService
	userHandler        *handlers.UserHandler
	redisRepo          *repositories.RedisRepository
	cfg                *config.Config
	logger             *zap.Logger
}

func NewApiV1(userHandler *handlers.UserHandler, authService services.IAuthService, healthCheckHandler *handlers.HealthCheckHandler, redisRepo *repositories.RedisRepository, cfg *config.Config, logger *zap.Logger) *ApiV1 {
	return &ApiV1{
		userHandler:        userHandler,
		authService:        authService,
		healthCheckHandler: healthCheckHandler,
		redisRepo:          redisRepo,
		cfg:                cfg,
		logger:             logger,
	}
}

func (a *ApiV1) InitApiV1Groups(group *gin.RouterGroup) {
	v1Router := group.Group("/v1")

	cache := middlewares.NewCacheMiddleware(a.redisRepo, a.cfg, a.logger)
	v1Router.Use(cache.Handle())
	{
		v1Router.GET("/health-check", a.healthCheckHandler.HealthCheck)
		routes.UserRoutes(v1Router, a.userHandler, a.authService)
	}
}
