package application

import (
	"context"
	"time"

	"hexagonal-go-backend/internal/modules/auth/domain"
	users "hexagonal-go-backend/internal/modules/users/application"
	userdomain "hexagonal-go-backend/internal/modules/users/domain"
	"hexagonal-go-backend/internal/platform/identifier"
)

type authService struct {
	users                 users.UserRepository
	hasher                users.PasswordHasher
	tokens                TokenProvider
	cache                 Cache
	accessTTL, refreshTTL time.Duration
}

func NewAuthService(userRepository users.UserRepository, hasher users.PasswordHasher, tokens TokenProvider, cache Cache, accessTTL, refreshTTL time.Duration) AuthService {
	return &authService{users: userRepository, hasher: hasher, tokens: tokens, cache: cache, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (s *authService) Login(ctx context.Context, email, password string) (TokenPair, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil || !u.Active || s.hasher.Compare(u.PasswordHash, password) != nil {
		return TokenPair{}, domain.ErrUnauthorized
	}
	return s.issue(ctx, u)
}

func (s *authService) Refresh(ctx context.Context, token string) (TokenPair, error) {
	claims, err := s.tokens.Parse(token)
	if err != nil {
		return TokenPair{}, domain.ErrInvalidToken
	}
	if _, err = s.cache.Get(ctx, "refresh:"+claims.TokenID); err != nil {
		return TokenPair{}, domain.ErrInvalidToken
	}
	_ = s.cache.Delete(ctx, "refresh:"+claims.TokenID)
	u, err := s.users.FindByID(ctx, claims.Subject)
	if err != nil || !u.Active {
		return TokenPair{}, domain.ErrUnauthorized
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

func (s *authService) issue(ctx context.Context, user *userdomain.User) (TokenPair, error) {
	now := time.Now().UTC()
	accessExp, refreshExp := now.Add(s.accessTTL), now.Add(s.refreshTTL)
	accessID, refreshID := identifier.New(), identifier.New()
	accessToken, err := s.tokens.Generate(TokenClaims{Subject: user.ID, Role: user.Role, ExpiresAt: accessExp, TokenID: accessID})
	if err != nil {
		return TokenPair{}, err
	}
	refreshToken, err := s.tokens.Generate(TokenClaims{Subject: user.ID, Role: user.Role, ExpiresAt: refreshExp, TokenID: refreshID})
	if err != nil {
		return TokenPair{}, err
	}
	if err = s.cache.Set(ctx, "refresh:"+refreshID, []byte(user.ID), s.refreshTTL); err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken, AccessExpiresAt: accessExp, RefreshExpiresAt: refreshExp}, nil
}
