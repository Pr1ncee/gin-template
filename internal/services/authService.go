package services

import (
	"GinBox/config"
	"GinBox/internal/appErrors"
	db "GinBox/internal/postgresql"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type IAuthService interface {
	generateToken(userClaims UserClaims, issuer, subject string) (string, error)
	GenerateAccessToken(userClaims UserClaims) (string, error)
	GenerateRefreshToken(userClaims UserClaims) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	ParseToken(tokenString string) (*JWTClaims, error)
	GeneratePassword(password []byte) ([]byte, error)
	ValidatePassword(hashedPassword, password []byte) error
}

type UserClaims struct {
	UserID int32       `json:"user_id"`
	Email  string      `json:"email"`
	Role   db.UserRole `json:"role"`
}

type JWTClaims struct {
	UserClaims
	jwt.RegisteredClaims
}

// AuthService handles JWT token operations
type AuthService struct {
	cfg    config.Config
	logger *zap.Logger
}

func NewAuthService(cfg config.Config, logger *zap.Logger) *AuthService {
	return &AuthService{cfg: cfg, logger: logger}
}

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

// GenerateAccessToken generates a new access token
func (a *AuthService) GenerateAccessToken(userClaims UserClaims) (string, error) {
	return a.generateToken(userClaims, a.cfg.Auth.TokenIssuer, "access-token")
}

// GenerateRefreshToken generates a new refresh token
func (a *AuthService) GenerateRefreshToken(userClaims UserClaims) (string, error) {
	return a.generateToken(userClaims, a.cfg.Auth.TokenIssuer, "refresh-token")
}

// ValidateToken validates and parses a JWT token
func (a *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, appErrors.ErrInvalidSigningMethod
		}
		return a.cfg.Auth.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, appErrors.ErrInvalidToken
}

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
		return nil, appErrors.ErrInvalidToken
	}

	a.logger.Info("Claims parsed", zap.Any("claims", claims))
	return claims, nil
}

func (a *AuthService) GeneratePassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, a.cfg.Auth.PasswordCost)
}

func (a *AuthService) ValidatePassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
