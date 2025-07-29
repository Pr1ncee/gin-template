package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// RequestLogger represents logger middleware that logs each request and calculates the duration of the latter.
type RequestLogger struct {
	logger *zap.Logger
}

// NewRequestLogger creates a new RequestLogger object.
func NewRequestLogger(logger *zap.Logger) *RequestLogger {
	return &RequestLogger{logger: logger}
}

// Handle implements logging functionality of each request.
func (r *RequestLogger) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		r.logger.Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
