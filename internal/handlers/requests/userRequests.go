package requests

type CreateUserRequest struct {
	FirstName   string `json:"first_name" binding:"required,min=2,max=255"`
	LastName    string `json:"last_name" binding:"required,min=2,max=255"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Role        string `json:"role,omitempty" binding:"oneof=User Admin"`
}

type UpdateUserRequest struct {
	FirstName       string `json:"first_name,omitempty" binding:"omitempty,min=2,max=255"`
	LastName        string `json:"last_name,omitempty" binding:"omitempty,min=2,max=255"`
	Email           string `json:"email,omitempty" binding:"omitempty,email"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	CurrentPassword string `json:"current_password,omitempty" binding:"omitempty"`
	NewPassword     string `json:"new_password,omitempty" binding:"omitempty,min=8"`
	Role            string `json:"role,omitempty" binding:"omitempty,oneof=User Admin"`
}
