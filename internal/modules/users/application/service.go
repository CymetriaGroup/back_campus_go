package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"hexagonal-go-backend/internal/modules/users/domain"
)

type userService struct {
	repo   UserRepository
	hasher PasswordHasher
	mailer Mailer
}

func NewUserService(repo UserRepository, hasher PasswordHasher, mailer Mailer) UserService {
	return &userService{repo: repo, hasher: hasher, mailer: mailer}
}

func (s *userService) Create(ctx context.Context, in CreateUserInput) (*domain.User, error) {
	in.Name, in.Email = strings.TrimSpace(in.Name), strings.ToLower(strings.TrimSpace(in.Email))
	if len(in.Name) < 2 || in.Email == "" || len(in.Password) < 8 {
		return nil, domain.ErrInvalidInput
	}
	if _, err := s.repo.FindByEmail(ctx, in.Email); err == nil {
		return nil, domain.ErrEmailAlreadyExists
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	hash, err := s.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}
	if in.Role == "" {
		in.Role = domain.RoleUser
	}
	now := time.Now().UTC()
	u := &domain.User{ID: newID(), Name: in.Name, Email: in.Email, PasswordHash: hash, Role: in.Role, Active: true, CreatedAt: now, UpdatedAt: now}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	if s.mailer != nil {
		_ = s.mailer.Send(ctx, MailMessage{To: u.Email, Subject: "Welcome", Body: "Your account was created."})
	}
	return u, nil
}

func (s *userService) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) List(ctx context.Context, filter UserFilter) ([]domain.User, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return s.repo.List(ctx, filter)
}

func (s *userService) Update(ctx context.Context, id string, in UpdateUserInput) (*domain.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		u.Name = name
	}
	if email := strings.ToLower(strings.TrimSpace(in.Email)); email != "" && email != u.Email {
		if _, findErr := s.repo.FindByEmail(ctx, email); findErr == nil {
			return nil, domain.ErrEmailAlreadyExists
		} else if !errors.Is(findErr, domain.ErrUserNotFound) {
			return nil, findErr
		}
		u.Email = email
	}
	if in.Active != nil {
		u.Active = *in.Active
	}
	if in.Role != "" {
		u.Role = in.Role
	}
	u.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func newID() string { return time.Now().UTC().Format("20060102150405.000000000") }
