package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"hexagonal-go-backend/internal/modules/users/application"
	"hexagonal-go-backend/internal/modules/users/domain"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: make(map[string]domain.User)}
}

func (r *UserRepository) Create(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.users {
		if existing.Email == user.Email {
			return domain.ErrEmailAlreadyExists
		}
	}
	r.users[user.ID] = *user
	return nil
}

func (r *UserRepository) FindByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return clone(user), nil
}

func (r *UserRepository) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	email = strings.ToLower(strings.TrimSpace(email))
	for _, user := range r.users {
		if user.Email == email {
			return clone(user), nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (r *UserRepository) List(_ context.Context, filter application.UserFilter) ([]domain.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]domain.User, 0, len(r.users))
	query := strings.ToLower(filter.Search)
	for _, user := range r.users {
		if query != "" && !strings.Contains(strings.ToLower(user.Name+" "+user.Email), query) {
			continue
		}
		if filter.Role != "" && user.Role != filter.Role {
			continue
		}
		if filter.Active != nil && user.Active != *filter.Active {
			continue
		}
		all = append(all, user)
	}
	sort.Slice(all, func(i, j int) bool {
		less := all[i].CreatedAt.Before(all[j].CreatedAt)
		if filter.Sort == "name" {
			less = all[i].Name < all[j].Name
		}
		if filter.Order == "desc" {
			return !less
		}
		return less
	})
	total := len(all)
	start := (filter.Page - 1) * filter.Limit
	if start > total {
		start = total
	}
	end := start + filter.Limit
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func (r *UserRepository) Update(_ context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return domain.ErrUserNotFound
	}
	for id, existing := range r.users {
		if id != user.ID && existing.Email == user.Email {
			return domain.ErrEmailAlreadyExists
		}
	}
	r.users[user.ID] = *user
	return nil
}

func (r *UserRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return domain.ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}

func clone(user domain.User) *domain.User {
	copy := user
	return &copy
}
