package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"hexagonal-go-backend/internal/core/domain"
	"hexagonal-go-backend/internal/core/ports"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]domain.User
}

func NewUserRepository() *UserRepository { return &UserRepository{users: make(map[string]domain.User)} }
func (r *UserRepository) Create(_ context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, v := range r.users {
		if v.Email == u.Email {
			return domain.ErrEmailAlreadyExists
		}
	}
	r.users[u.ID] = *u
	return nil
}
func (r *UserRepository) FindByID(_ context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return clone(u), nil
}
func (r *UserRepository) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	email = strings.ToLower(strings.TrimSpace(email))
	for _, u := range r.users {
		if u.Email == email {
			return clone(u), nil
		}
	}
	return nil, domain.ErrUserNotFound
}
func (r *UserRepository) List(_ context.Context, f ports.UserFilter) ([]domain.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]domain.User, 0, len(r.users))
	q := strings.ToLower(f.Search)
	for _, u := range r.users {
		if q != "" && !strings.Contains(strings.ToLower(u.Name+" "+u.Email), q) {
			continue
		}
		if f.Role != "" && u.Role != f.Role {
			continue
		}
		if f.Active != nil && u.Active != *f.Active {
			continue
		}
		all = append(all, u)
	}
	sort.Slice(all, func(i, j int) bool {
		less := all[i].CreatedAt.Before(all[j].CreatedAt)
		if f.Sort == "name" {
			less = all[i].Name < all[j].Name
		}
		if f.Order == "desc" {
			return !less
		}
		return less
	})
	total := len(all)
	start := (f.Page - 1) * f.Limit
	if start > total {
		start = total
	}
	end := start + f.Limit
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}
func (r *UserRepository) Update(_ context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[u.ID]; !ok {
		return domain.ErrUserNotFound
	}
	for id, v := range r.users {
		if id != u.ID && v.Email == u.Email {
			return domain.ErrEmailAlreadyExists
		}
	}
	r.users[u.ID] = *u
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
func clone(u domain.User) *domain.User { v := u; return &v }
