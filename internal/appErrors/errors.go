/*
Package appErrors provides all error types that should be handled.

Specifically, this file provides errors and the function that handles error responses in the middleware layer.
*/
package appErrors

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// ResponseBody represents the structure of the response if error occurs.
type ResponseBody struct {
	Message string `json:"message"`
	Code    int    `json:"status"`
}

var (
	ErrInvalidSigningMethod         = errors.New("the signing method is invalid")
	ErrInvalidToken                 = errors.New("the token is invalid")
	ErrInsufficientTokenClaims      = errors.New("the token is insufficient")
	ErrPasswordDoesNotMatch         = errors.New("the password does not match")
	ErrFailedToGenerateAccessToken  = errors.New("failed to generate access token")
	ErrFailedToGenerateRefreshToken = errors.New("failed to generate refresh token")
	ErrFailedToCreateUser           = errors.New("failed to create user")
	ErrUserNotFound                 = errors.New("user not found")
)

// HandleErr handles different errors. It's a function that then passed into middleware in the router.
func HandleErr(ctx *gin.Context) {
	ctx.Next()
	for _, err := range ctx.Errors {
		var res ResponseBody
		switch {
		case errors.Is(err.Err, ErrInvalidSigningMethod):
			res = ResponseBody{ErrInvalidSigningMethod.Error(), http.StatusBadRequest}
		case errors.Is(err.Err, ErrInvalidToken):
			res = ResponseBody{ErrInvalidToken.Error(), http.StatusBadRequest}
		case errors.Is(err.Err, ErrInsufficientTokenClaims):
			res = ResponseBody{ErrInsufficientTokenClaims.Error(), http.StatusBadRequest}
		case errors.Is(err.Err, ErrPasswordDoesNotMatch):
			res = ResponseBody{ErrPasswordDoesNotMatch.Error(), http.StatusForbidden}
		case errors.Is(err.Err, ErrFailedToGenerateAccessToken):
			res = ResponseBody{ErrFailedToGenerateAccessToken.Error(), http.StatusBadRequest}
		case errors.Is(err.Err, ErrFailedToGenerateRefreshToken):
			res = ResponseBody{ErrFailedToGenerateRefreshToken.Error(), http.StatusBadRequest}
		case errors.Is(err.Err, ErrUserNotFound):
			res = ResponseBody{ErrUserNotFound.Error(), http.StatusNotFound}
		default:
			res = ResponseBody{
				fmt.Sprintf("unhandled error occurred: %s", err.Err.Error()),
				http.StatusInternalServerError,
			}
		}
		log.Printf("Registred error: %d|%s", res.Code, err.Error())
		ctx.AbortWithStatusJSON(res.Code, res.Message)
	}
}
