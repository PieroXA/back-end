package user

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=owner vet worker"`
}

type UpdateUserRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
	FullName string  `json:"full_name" validate:"required"`
	Role     string  `json:"role" validate:"required,oneof=owner vet worker"`
	IsActive *bool   `json:"is_active" validate:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
