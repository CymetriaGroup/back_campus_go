package application_test

import (
	"context"
	"errors"
	"testing"

	"hexagonal-go-backend/internal/modules/users/application"
	"hexagonal-go-backend/internal/modules/users/domain"
	"hexagonal-go-backend/internal/modules/users/infrastructure/persistence/memory"
	"hexagonal-go-backend/internal/modules/users/infrastructure/security"
)

type noopMailer struct{}

func (noopMailer) Send(context.Context, application.MailMessage) error { return nil }

func TestCreateRejectsDuplicateEmail(t *testing.T) {
	service := application.NewUserService(memory.NewUserRepository(), security.NewBcryptHasher(4), noopMailer{})
	input := application.CreateUserInput{Name: "Jane Doe", Email: "jane@example.com", Password: "password123"}
	if _, err := service.Create(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Create(context.Background(), input); !errors.Is(err, domain.ErrEmailAlreadyExists) {
		t.Fatalf("expected duplicate email, got %v", err)
	}
}

func TestListCapsLimit(t *testing.T) {
	service := application.NewUserService(memory.NewUserRepository(), security.NewBcryptHasher(4), noopMailer{})
	for _, email := range []string{"a@example.com", "b@example.com"} {
		if _, err := service.Create(context.Background(), application.CreateUserInput{Name: "A User", Email: email, Password: "password123"}); err != nil {
			t.Fatal(err)
		}
	}
	users, total, err := service.List(context.Background(), application.UserFilter{Page: 1, Limit: 1000})
	if err != nil || len(users) != 2 || total != 2 {
		t.Fatalf("users=%d total=%d err=%v", len(users), total, err)
	}
}
