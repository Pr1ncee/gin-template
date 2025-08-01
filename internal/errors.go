package internal

import "errors"

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
