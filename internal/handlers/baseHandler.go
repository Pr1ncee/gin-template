/*
Package handlers provides implementation of endpoints.

Specifically, this file implements utils methods that help set up new handlers, handling success and error responses.
*/
package handlers

import (
	"GinBox/internal/handlers/responses"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// BaseHandler provides basic tools which can be used inside handlers.
type BaseHandler struct {
	logger  *zap.Logger
	timeout time.Duration
}

// NewBaseHandler creates a new BaseHandler with logger and timeout value.
func NewBaseHandler(logger *zap.Logger, timeout time.Duration) *BaseHandler {
	return &BaseHandler{
		logger:  logger,
		timeout: timeout,
	}
}

// SendError sends an error response
func (h *BaseHandler) SendError(c *gin.Context, statusCode int, message string, err error) {
	h.logger.Error(message, zap.Error(err))
	c.JSON(statusCode, responses.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Code:    statusCode,
		Message: message,
	})
}

// SendSuccess sends a success response
func (h *BaseHandler) SendSuccess(c *gin.Context, statusCode int, data interface{}, message string) {
	response := responses.SuccessResponse{
		Data:    data,
		Message: message,
	}
	c.JSON(statusCode, response)
}

// GetContextWithTimeout returns a context with timeout
func (h *BaseHandler) GetContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), h.timeout)
}
