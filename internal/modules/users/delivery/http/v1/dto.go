package v1

import (
	"time"

	"hexagonal-go-backend/internal/modules/users/domain"
)

type CreateUserRequest struct {
	Name     string      `json:"name" binding:"required,min=2,max=100"`
	Email    string      `json:"email" binding:"required,email"`
	Password string      `json:"password" binding:"required,min=8"`
	Role     domain.Role `json:"role" binding:"omitempty,oneof=admin user"`
}

type UpdateUserRequest struct {
	Name   string      `json:"name" binding:"omitempty,min=2,max=100"`
	Email  string      `json:"email" binding:"omitempty,email"`
	Active *bool       `json:"active"`
	Role   domain.Role `json:"role" binding:"omitempty,oneof=admin user"`
}

type UserResponse struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Role      domain.Role `json:"role"`
	Active    bool        `json:"active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func userResponse(user *domain.User) UserResponse {
	return UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, Active: user.Active, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
}
