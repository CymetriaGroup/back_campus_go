package ent

import (
	"context"
	"errors"
	"fmt"
	"strings"

	entclient "hexagonal-go-backend/internal/ent"
	entuser "hexagonal-go-backend/internal/ent/user"
	"hexagonal-go-backend/internal/modules/users/application"
	"hexagonal-go-backend/internal/modules/users/domain"

	"entgo.io/ent/dialect/sql"
)

type UserRepository struct {
	client *entclient.Client
}

func NewUserRepository(client *entclient.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.client.User.Create().
		SetID(user.ID).
		SetName(user.Name).
		SetEmail(user.Email).
		SetPasswordHash(user.PasswordHash).
		SetRole(entuser.Role(user.Role)).
		SetActive(user.Active).
		SetCreatedAt(user.CreatedAt).
		SetUpdatedAt(user.UpdatedAt).
		Save(ctx)
	return translateError(err)
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, translateError(err)
	}
	return toDomain(user), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.client.User.Query().
		Where(entuser.EmailEQ(strings.ToLower(strings.TrimSpace(email)))).
		Only(ctx)
	if err != nil {
		return nil, translateError(err)
	}
	return toDomain(user), nil
}

func (r *UserRepository) List(ctx context.Context, filter application.UserFilter) ([]domain.User, int, error) {
	query := r.client.User.Query()
	if search := strings.TrimSpace(filter.Search); search != "" {
		query = query.Where(entuser.Or(entuser.NameContainsFold(search), entuser.EmailContainsFold(search)))
	}
	if filter.Role != "" {
		query = query.Where(entuser.RoleEQ(entuser.Role(filter.Role)))
	}
	if filter.Active != nil {
		query = query.Where(entuser.ActiveEQ(*filter.Active))
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	order := sql.OrderAsc()
	if filter.Order == "desc" {
		order = sql.OrderDesc()
	}
	if filter.Sort == "name" {
		query = query.Order(entuser.ByName(order))
	} else {
		query = query.Order(entuser.ByCreatedAt(order))
	}

	users, err := query.
		Offset((filter.Page - 1) * filter.Limit).
		Limit(filter.Limit).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	result := make([]domain.User, 0, len(users))
	for _, user := range users {
		result = append(result, *toDomain(user))
	}
	return result, total, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := r.client.User.UpdateOneID(user.ID).
		SetName(user.Name).
		SetEmail(user.Email).
		SetPasswordHash(user.PasswordHash).
		SetRole(entuser.Role(user.Role)).
		SetActive(user.Active).
		SetUpdatedAt(user.UpdatedAt).
		Save(ctx)
	return translateError(err)
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return translateError(r.client.User.DeleteOneID(id).Exec(ctx))
}

func translateError(err error) error {
	if err == nil {
		return nil
	}
	if entclient.IsNotFound(err) {
		return domain.ErrUserNotFound
	}
	if entclient.IsConstraintError(err) {
		return domain.ErrEmailAlreadyExists
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return fmt.Errorf("users repository: %w", err)
}

var _ application.UserRepository = (*UserRepository)(nil)
