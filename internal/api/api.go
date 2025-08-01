package api

import (
	"GinBox/config"
	v1 "GinBox/internal/api/v1"
	"GinBox/internal/handlers"
	"GinBox/internal/middlewares"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ApiGroup struct {
	handler *handlers.UserHandler
	ApiV1   *v1.ApiV1
	cfg     *config.Config
	logger  *zap.Logger
}

func NewApiGroup(handler *handlers.UserHandler, apiV1 *v1.ApiV1, cfg *config.Config, logger *zap.Logger) *ApiGroup {
	return &ApiGroup{handler: handler, ApiV1: apiV1, cfg: cfg, logger: logger}
}

func (a *ApiGroup) InitRouterGroups(router *gin.Engine) {
	prometheusLogger := middlewares.NewPrometheusLogger(a.logger)

	corsDefaultConfig := cors.DefaultConfig()
	corsDefaultConfig.AllowAllOrigins = false
	corsDefaultConfig.AllowOrigins = a.cfg.App.Origins
	corsDefaultConfig.AllowCredentials = true
	corsDefaultConfig.AddAllowHeaders("Authorization")

	rateLimiter := middlewares.NewRateLimiter(a.cfg)

	errorMiddleware := middlewares.NewErrorHandlerMiddleware(a.logger)

	prometheusLogger.Setup(router)
	router.Use(rateLimiter.Handle())
	router.Use(cors.Default())
	router.Use(gin.Recovery())
	router.Use(errorMiddleware.Handle())
	api := router.Group("/api")
	{
		a.ApiV1.InitApiV1Groups(api)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.NoRoute(func(c *gin.Context) {
		a.logger.Warn(fmt.Sprintf("Route not found: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP()))
		c.JSON(http.StatusNotFound, gin.H{
			"error":     "Route not found",
			"code":      404,
			"message":   fmt.Sprintf("The requested route %s %s does not exist", c.Request.Method, c.Request.URL.Path),
			"timestamp": time.Now().UTC(),
		})
	})

}
