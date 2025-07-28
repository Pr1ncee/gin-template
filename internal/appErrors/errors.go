package appErrors

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Res struct {
	Message string `json:"message"`
	Code    int    `json:"status"`
}

var (
	ErrInvalidSigningMethod         = errors.New("the signing method is invalid")
	ErrInvalidToken                 = errors.New("the token is invalid")
	ErrPasswordDoesNotMatch         = errors.New("the password does not match")
	ErrFailedToGenerateAccessToken  = errors.New("failed to generate access token")
	ErrFailedToGenerateRefreshToken = errors.New("failed to generate refresh token")
	ErrFailedToCreateUser           = errors.New("failed to create user")
	ErrUserNotFound                 = errors.New("user not found")
)

func HandleErr(ctx *gin.Context) {
	ctx.Next()
	for _, err := range ctx.Errors {
		var res Res
		switch {
		case errors.Is(err.Err, ErrInvalidToken):
			res = Res{ErrInvalidToken.Error(), http.StatusBadRequest}
		case errors.Is(err.Err, ErrInvalidSigningMethod):
			res = Res{ErrInvalidSigningMethod.Error(), http.StatusBadRequest}
		case errors.Is(err.Err, ErrPasswordDoesNotMatch):
			res = Res{ErrPasswordDoesNotMatch.Error(), http.StatusForbidden}
		case errors.Is(err.Err, ErrFailedToGenerateAccessToken):
			res = Res{ErrFailedToGenerateAccessToken.Error(), http.StatusBadRequest}
		case errors.Is(err.Err, ErrFailedToGenerateRefreshToken):
			res = Res{ErrFailedToGenerateRefreshToken.Error(), http.StatusBadRequest}
		case errors.Is(err.Err, ErrUserNotFound):
			res = Res{ErrUserNotFound.Error(), http.StatusNotFound}
		default:
			res = Res{fmt.Sprintf("unhandled error occurred: %s", err.Err.Error()), http.StatusInternalServerError}
		}
		log.Printf("Registred error: %d|%s", res.Code, err.Error())
		ctx.AbortWithStatusJSON(res.Code, res.Message)
	}
}
