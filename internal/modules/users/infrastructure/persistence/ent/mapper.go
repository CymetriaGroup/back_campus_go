package ent

import (
	entuser "hexagonal-go-backend/internal/ent"
	"hexagonal-go-backend/internal/modules/users/domain"
)

func toDomain(user *entuser.User) *domain.User {
	if user == nil {
		return nil
	}
	status := domain.StatusActive
	if user.Status != "" {
		status = domain.UserStatus(user.Status)
	}
	return &domain.User{
		ID:        user.ID,
		AuthID:    user.AuthID,
		TenantID:  user.TenantID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      domain.Role(user.Role),
		Status:    status,
		Active:    user.Active,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
