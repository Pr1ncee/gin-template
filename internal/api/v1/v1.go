package v1

import (
	"GinBox/internal/handlers"
	"GinBox/internal/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApiV1 struct {
	healthCheckHandler *handlers.HealthCheckHandler
	userHandler        *handlers.UserHandler
	logger             *zap.Logger
}

func NewApiV1(userHandler *handlers.UserHandler, healthCheckHandler *handlers.HealthCheckHandler, logger *zap.Logger) *ApiV1 {
	return &ApiV1{userHandler: userHandler, healthCheckHandler: healthCheckHandler, logger: logger}
}

func (a *ApiV1) InitApiV1Groups(group *gin.RouterGroup) {
	v1Router := group.Group("/v1")
	{
		v1Router.GET("/health-check", a.healthCheckHandler.HealthCheck)
		routes.UserRoutes(v1Router, a.userHandler, a.logger)
	}
}
