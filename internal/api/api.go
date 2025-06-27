package api

import (
	v1 "GinBox/internal/api/v1"
	"GinBox/internal/handlers"
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
	logger  *zap.Logger
}

func NewApiGroup(handler *handlers.UserHandler, apiV1 *v1.ApiV1, logger *zap.Logger) *ApiGroup {
	return &ApiGroup{handler: handler, ApiV1: apiV1, logger: logger}
}

func (a *ApiGroup) InitRouterGroups(router *gin.Engine) {
	router.Use(cors.Default())
	router.Use(gin.Recovery())
	//router.Use(appErrors.HandleErr)  TODO HANDLE ERRORS HERE
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
