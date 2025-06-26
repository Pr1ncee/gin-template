package handlers

import (
	"GinBox/internal/handlers/responses"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type BaseHandler struct {
	logger  *logrus.Logger
	timeout time.Duration
}

func NewBaseHandler(logger *logrus.Logger, timeout time.Duration) *BaseHandler {
	return &BaseHandler{
		logger:  logger,
		timeout: timeout,
	}
}

// SendError sends an error response
func (h *BaseHandler) SendError(c *gin.Context, statusCode int, message string, err error) {
	h.logger.WithError(err).Error(message)
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
