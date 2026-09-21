package application

import (
	"context"
	"time"

	userdomain "hexagonal-go-backend/internal/modules/users/domain"
)

type TokenClaims struct {
	Subject   string
	Role      userdomain.Role
	ExpiresAt time.Time
	TokenID   string
}

type TokenProvider interface {
	Generate(TokenClaims) (string, error)
	Parse(string) (TokenClaims, error)
}

type TokenPair struct {
	AccessToken, RefreshToken         string
	AccessExpiresAt, RefreshExpiresAt time.Time
}

type AuthService interface {
	Login(context.Context, string, string) (TokenPair, error)
	Refresh(context.Context, string) (TokenPair, error)
	Logout(context.Context, string) error
}

type Cache interface {
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, []byte, time.Duration) error
	Delete(context.Context, string) error
}
