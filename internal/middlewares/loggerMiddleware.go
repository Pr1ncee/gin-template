package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
	"time"
)

// PrometheusLogger represents Prometheus middleware that logs each request with different metrics
// and custom written logger middleware that calculates the duration of the requests.
type PrometheusLogger struct {
	logger *zap.Logger
}

// NewPrometheusLogger creates a new PrometheusLogger with configured logger.
func NewPrometheusLogger(logger *zap.Logger) *PrometheusLogger {
	return &PrometheusLogger{logger: logger}
}

// Setup configures Gin compatible Prometheus application with custom logger middleware
func (pl *PrometheusLogger) Setup(router *gin.Engine) {
	p := ginprometheus.NewPrometheus("gin")
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		return c.FullPath()
	}

	// Custom middleware for zap logging
	router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		pl.logger.Info("Incoming request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
		)
	})

	p.Use(router)
}
