package services

import (
	"context"
	"errors"
	"time"

	"hexagonal-go-backend/internal/core/domain"
	"hexagonal-go-backend/internal/core/ports"
)

type authService struct {
	users                 ports.UserRepository
	hasher                ports.PasswordHasher
	tokens                ports.TokenProvider
	cache                 ports.Cache
	accessTTL, refreshTTL time.Duration
}

func NewAuthService(u ports.UserRepository, h ports.PasswordHasher, t ports.TokenProvider, c ports.Cache, access, refresh time.Duration) ports.AuthService {
	return &authService{u, h, t, c, access, refresh}
}
func (s *authService) Login(ctx context.Context, email, password string) (ports.TokenPair, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil || !u.Active || s.hasher.Compare(u.PasswordHash, password) != nil {
		return ports.TokenPair{}, domain.ErrUnauthorized
	}
	return s.issue(ctx, u)
}
func (s *authService) Refresh(ctx context.Context, token string) (ports.TokenPair, error) {
	claims, err := s.tokens.Parse(token)
	if err != nil {
		return ports.TokenPair{}, domain.ErrInvalidToken
	}
	if _, err = s.cache.Get(ctx, "refresh:"+claims.TokenID); err != nil {
		return ports.TokenPair{}, domain.ErrInvalidToken
	}
	_ = s.cache.Delete(ctx, "refresh:"+claims.TokenID)
	u, err := s.users.FindByID(ctx, claims.Subject)
	if err != nil || !u.Active {
		return ports.TokenPair{}, domain.ErrUnauthorized
	}
	return s.issue(ctx, u)
}
func (s *authService) Logout(ctx context.Context, token string) error {
	claims, err := s.tokens.Parse(token)
	if err != nil {
		return domain.ErrInvalidToken
	}
	return s.cache.Delete(ctx, "refresh:"+claims.TokenID)
}
func (s *authService) issue(ctx context.Context, u *domain.User) (ports.TokenPair, error) {
	now := time.Now().UTC()
	accessExp, refreshExp := now.Add(s.accessTTL), now.Add(s.refreshTTL)
	accessID, refreshID := newID()+"a", newID()+"r"
	a, err := s.tokens.Generate(ports.TokenClaims{Subject: u.ID, Role: u.Role, ExpiresAt: accessExp, TokenID: accessID})
	if err != nil {
		return ports.TokenPair{}, err
	}
	r, err := s.tokens.Generate(ports.TokenClaims{Subject: u.ID, Role: u.Role, ExpiresAt: refreshExp, TokenID: refreshID})
	if err != nil {
		return ports.TokenPair{}, err
	}
	if err = s.cache.Set(ctx, "refresh:"+refreshID, []byte(u.ID), s.refreshTTL); err != nil {
		return ports.TokenPair{}, err
	}
	return ports.TokenPair{AccessToken: a, RefreshToken: r, AccessExpiresAt: accessExp, RefreshExpiresAt: refreshExp}, nil
}

var _ = errors.Is
