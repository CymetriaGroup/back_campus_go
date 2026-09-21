package application

import (
	"context"

	"hexagonal-go-backend/internal/modules/users/domain"
)

type UserFilter struct {
	Page, Limit int
	Search      string
	Role        domain.Role
	Active      *bool
	Sort, Order string
}

type CreateUserInput struct {
	Name, Email, Password string
	Role                  domain.Role
}

type UpdateUserInput struct {
	Name, Email string
	Active      *bool
	Role        domain.Role
}

type UserRepository interface {
	Create(context.Context, *domain.User) error
	FindByID(context.Context, string) (*domain.User, error)
	FindByEmail(context.Context, string) (*domain.User, error)
	List(context.Context, UserFilter) ([]domain.User, int, error)
	Update(context.Context, *domain.User) error
	Delete(context.Context, string) error
}

type UserService interface {
	Create(context.Context, CreateUserInput) (*domain.User, error)
	FindByID(context.Context, string) (*domain.User, error)
	List(context.Context, UserFilter) ([]domain.User, int, error)
	Update(context.Context, string, UpdateUserInput) (*domain.User, error)
	Delete(context.Context, string) error
}

type PasswordHasher interface {
	Hash(string) (string, error)
	Compare(string, string) error
}

type MailMessage struct{ To, Subject, Body string }

type Mailer interface {
	Send(context.Context, MailMessage) error
}
