package http_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cache "hexagonal-go-backend/internal/adapters/cache/memory"
	httpadapter "hexagonal-go-backend/internal/adapters/handler/http"
	repository "hexagonal-go-backend/internal/adapters/repository/memory"
	"hexagonal-go-backend/internal/adapters/security"
	"hexagonal-go-backend/internal/config"
	"hexagonal-go-backend/internal/core/domain"
	"hexagonal-go-backend/internal/core/ports"
	"hexagonal-go-backend/internal/core/services"
)

type mailer struct{}

func (mailer) Send(context.Context, ports.MailMessage) error { return nil }
func TestLoginAndProtectedUsers(t *testing.T) {
	repo := repository.NewUserRepository()
	hasher := security.NewBcryptHasher(4)
	tokens := security.NewHMACTokenProvider("test-secret")
	users := services.NewUserService(repo, hasher, mailer{})
	_, err := users.Create(context.Background(), ports.CreateUserInput{Name: "Admin User", Email: "admin@example.com", Password: "password123", Role: domain.RoleAdmin})
	if err != nil {
		t.Fatal(err)
	}
	auth := services.NewAuthService(repo, hasher, tokens, cache.New(), time.Minute, time.Hour)
	cfg := config.Config{App: config.AppConfig{Env: "test"}, Security: config.SecurityConfig{RateLimit: 100}}
	router := httpadapter.NewRouter(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), httpadapter.NewUserHandler(users), httpadapter.NewAuthHandler(auth), tokens)
	login := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"email":"admin@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(login, req)
	if login.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", login.Code, login.Body.String())
	}
	var body struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	if err = json.Unmarshal(login.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	list := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+body.Data.AccessToken)
	router.ServeHTTP(list, req)
	if list.Code != http.StatusOK {
		t.Fatalf("list status=%d body=%s", list.Code, list.Body.String())
	}
	if list.Header().Get("X-Request-ID") == "" {
		t.Fatal("missing request id")
	}
}
