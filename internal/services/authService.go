/*
Package services provides methods with core business logic for different entities and domains of the application.

Specifically, this file implements all necessary methods
for authentication (token generation/validation and password generation/validation).
*/
package services

import (
	"GinBox/config"
	"GinBox/internal"
	db "GinBox/internal/postgresql"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// IAuthService defines the interface for actions with token and password.
type IAuthService interface {
	generateToken(userClaims UserClaims, issuer, subject string) (string, error)
	GenerateAccessToken(userClaims UserClaims) (string, error)
	GenerateRefreshToken(userClaims UserClaims) (string, error)
	ParseToken(tokenString string) (*JWTClaims, error)
	GeneratePassword(password []byte) ([]byte, error)
	ValidatePassword(hashedPassword, password []byte) error
}

// UserClaims defines the mandatory User information inside each generated token.
type UserClaims struct {
	UserID int32       `json:"user_id"`
	Email  string      `json:"email"`
	Role   db.UserRole `json:"role"`
}

// JWTClaims defines the mandatory information (beyond only User one) inside each generated token.
type JWTClaims struct {
	UserClaims
	jwt.RegisteredClaims
}

// AuthService implements IAuthService and handles JWT token operations.
type AuthService struct {
	cfg    config.Config
	logger *zap.Logger
}

// NewAuthService returns AuthService with input config and logger.
func NewAuthService(cfg config.Config, logger *zap.Logger) *AuthService {
	return &AuthService{cfg: cfg, logger: logger}
}

// generateToken generates a JWT token based on input parameters.
func (a *AuthService) generateToken(userClaims UserClaims, issuer, subject string) (string, error) {
	claims := JWTClaims{
		UserClaims: userClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.cfg.Auth.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			Subject:   subject,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.cfg.Auth.JWTSecret))
}

// GenerateAccessToken generates a new access token.
func (a *AuthService) GenerateAccessToken(userClaims UserClaims) (string, error) {
	return a.generateToken(userClaims, a.cfg.App.Name, "access-token")
}

// GenerateRefreshToken generates a new refresh token.
func (a *AuthService) GenerateRefreshToken(userClaims UserClaims) (string, error) {
	return a.generateToken(userClaims, a.cfg.App.Name, "refresh-token")
}

// ParseToken validates and parses a JWT token.
func (a *AuthService) ParseToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(a.cfg.Auth.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		a.logger.Error("Error parsing token", zap.Error(err))
		return nil, err
	}

	if !token.Valid {
		a.logger.Error("Error parsing token", zap.Error(err))
		return nil, internal.ErrInvalidToken
	}

	a.logger.Info("Claims parsed", zap.Any("claims", claims))
	return claims, nil
}

// GeneratePassword generates a new password with pre-configured salt.
func (a *AuthService) GeneratePassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, a.cfg.Auth.PasswordCost)
}

// ValidatePassword validates password.
func (a *AuthService) ValidatePassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
