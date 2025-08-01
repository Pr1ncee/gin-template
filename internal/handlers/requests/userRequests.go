/*
Package requests provides structs that describe input schemas for different endpoints.

Specifically, this file describes schemas for User's endpoints.
*/
package requests

// RefreshTokenRequest represents an input schema for refreshing access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LoginRequest represents an input schema for logging a User entity.
type LoginRequest struct {
	ID       int32  `json:"id"`
	Password string `json:"password"`
}

// CreateUserRequest represents an input schema for creating a new User entity.
type CreateUserRequest struct {
	FirstName   string `json:"first_name" binding:"required,min=2,max=255"`
	LastName    string `json:"last_name" binding:"required,min=2,max=255"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Role        string `json:"role,omitempty" binding:"oneof=User Admin"`
}

// UpdateUserRequest represents an input schema for updating a User via PATCH HTTP request.
type UpdateUserRequest struct {
	FirstName       string `json:"first_name,omitempty" binding:"omitempty,min=2,max=255"`
	LastName        string `json:"last_name,omitempty" binding:"omitempty,min=2,max=255"`
	Email           string `json:"email,omitempty" binding:"omitempty,email"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	CurrentPassword string `json:"current_password,omitempty" binding:"omitempty"`
	NewPassword     string `json:"new_password,omitempty" binding:"omitempty,min=8"`
	Role            string `json:"role,omitempty" binding:"omitempty,oneof=User Admin"`
}
