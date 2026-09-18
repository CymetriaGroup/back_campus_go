package services_test

import (
	"context"
	"errors"
	"testing"

	"hexagonal-go-backend/internal/adapters/repository/memory"
	"hexagonal-go-backend/internal/adapters/security"
	"hexagonal-go-backend/internal/core/domain"
	"hexagonal-go-backend/internal/core/ports"
	"hexagonal-go-backend/internal/core/services"
)

type noopMailer struct{}

func (noopMailer) Send(context.Context, ports.MailMessage) error { return nil }
func TestCreateRejectsDuplicateEmail(t *testing.T) {
	s := services.NewUserService(memory.NewUserRepository(), security.NewBcryptHasher(4), noopMailer{})
	in := ports.CreateUserInput{Name: "Jane Doe", Email: "jane@example.com", Password: "password123"}
	if _, err := s.Create(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create(context.Background(), in); !errors.Is(err, domain.ErrEmailAlreadyExists) {
		t.Fatalf("expected duplicate email, got %v", err)
	}
}
func TestListCapsLimit(t *testing.T) {
	s := services.NewUserService(memory.NewUserRepository(), security.NewBcryptHasher(4), noopMailer{})
	for _, email := range []string{"a@example.com", "b@example.com"} {
		if _, err := s.Create(context.Background(), ports.CreateUserInput{Name: "A User", Email: email, Password: "password123"}); err != nil {
			t.Fatal(err)
		}
	}
	users, total, err := s.List(context.Background(), ports.UserFilter{Page: 1, Limit: 1000})
	if err != nil || len(users) != 2 || total != 2 {
		t.Fatalf("users=%d total=%d err=%v", len(users), total, err)
	}
}
