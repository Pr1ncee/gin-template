package middlewares

import (
	"GinBox/internal"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// ErrorHandlerMiddleware represents an object which can be embedded into middleware.
type ErrorHandlerMiddleware struct {
	logger *zap.Logger
}

// NewErrorHandlerMiddleware creates a new ErrorHandlerMiddleware object.
func NewErrorHandlerMiddleware(logger *zap.Logger) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{logger: logger}
}

// ResponseBody represents the structure of the response if error occurs.
type ResponseBody struct {
	Message string `json:"message"`
	Code    int    `json:"status"`
}

// Handle handles different errors. It's a function that then passed into middleware in the router.
func (e *ErrorHandlerMiddleware) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) == 0 {
			return
		}

		for _, err := range ctx.Errors {
			var res ResponseBody
			switch {
			case errors.Is(err.Err, internal.ErrInvalidSigningMethod):
				res = ResponseBody{internal.ErrInvalidSigningMethod.Error(), http.StatusBadRequest}
			case errors.Is(err.Err, internal.ErrInvalidToken):
				res = ResponseBody{internal.ErrInvalidToken.Error(), http.StatusBadRequest}
			case errors.Is(err.Err, internal.ErrInsufficientTokenClaims):
				res = ResponseBody{internal.ErrInsufficientTokenClaims.Error(), http.StatusBadRequest}
			case errors.Is(err.Err, internal.ErrPasswordDoesNotMatch):
				res = ResponseBody{internal.ErrPasswordDoesNotMatch.Error(), http.StatusForbidden}
			case errors.Is(err.Err, internal.ErrFailedToGenerateAccessToken):
				res = ResponseBody{internal.ErrFailedToGenerateAccessToken.Error(), http.StatusBadRequest}
			case errors.Is(err.Err, internal.ErrFailedToGenerateRefreshToken):
				res = ResponseBody{internal.ErrFailedToGenerateRefreshToken.Error(), http.StatusBadRequest}
			case errors.Is(err.Err, internal.ErrUserNotFound):
				res = ResponseBody{internal.ErrUserNotFound.Error(), http.StatusNotFound}
			default:
				res = ResponseBody{
					fmt.Sprintf("unhandled error occurred: %s", err.Err.Error()),
					http.StatusInternalServerError,
				}
			}
			e.logger.Info("Registered error", zap.Int("Status Code", res.Code), zap.Error(err))
			ctx.AbortWithStatusJSON(res.Code, res.Message)
		}
	}
}
