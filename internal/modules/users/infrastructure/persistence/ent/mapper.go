package ent

import (
	entuser "hexagonal-go-backend/internal/ent"
	"hexagonal-go-backend/internal/modules/users/domain"
)

func toDomain(user *entuser.User) *domain.User {
	if user == nil {
		return nil
	}
	return &domain.User{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         domain.Role(user.Role),
		Active:       user.Active,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
