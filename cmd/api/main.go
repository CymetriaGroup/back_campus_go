package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	authapp "hexagonal-go-backend/internal/modules/auth/application"
	authhttp "hexagonal-go-backend/internal/modules/auth/delivery/http/v1"
	memorycache "hexagonal-go-backend/internal/modules/auth/infrastructure/cache/memory"
	authsecurity "hexagonal-go-backend/internal/modules/auth/infrastructure/security"
	usersapp "hexagonal-go-backend/internal/modules/users/application"
	usershttp "hexagonal-go-backend/internal/modules/users/delivery/http/v1"
	userdomain "hexagonal-go-backend/internal/modules/users/domain"
	usermail "hexagonal-go-backend/internal/modules/users/infrastructure/mail"
	memoryusers "hexagonal-go-backend/internal/modules/users/infrastructure/persistence/memory"
	usersecurity "hexagonal-go-backend/internal/modules/users/infrastructure/security"
	"hexagonal-go-backend/internal/platform/config"
	httpplatform "hexagonal-go-backend/internal/platform/http"
	"hexagonal-go-backend/internal/platform/logger"
	"hexagonal-go-backend/internal/platform/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg.App.LogLevel)
	userRepository := memoryusers.NewUserRepository()
	passwordHasher := usersecurity.NewBcryptHasher(0)
	mailer := usermail.NewLogMailer(log)
	tokenProvider := authsecurity.NewHMACTokenProvider(cfg.JWT.Secret)
	cache := memorycache.New()

	userService := usersapp.NewUserService(userRepository, passwordHasher, mailer)
	if err := seedAdmin(context.Background(), userService, cfg); err != nil {
		log.Error("failed to create development administrator", "error", err)
		os.Exit(1)
	}
	authService := authapp.NewAuthService(userRepository, passwordHasher, tokenProvider, cache, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	router := httpplatform.NewRouter(
		cfg,
		log,
		usershttp.NewController(userService),
		authhttp.NewController(authService),
		tokenProvider,
	)
	httpServer := server.New(router, cfg.App)

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("http server started", "address", httpServer.Addr, "environment", cfg.App.Env)
		serverErrors <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", "error", err)
			os.Exit(1)
		}
	case <-shutdownSignal.Done():
		log.Info("shutdown signal received")
	}

	shutdownContext, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(shutdownContext); err != nil {
		log.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	log.Info("http server stopped")
}

func seedAdmin(ctx context.Context, service usersapp.UserService, cfg config.Config) error {
	_, err := service.Create(ctx, usersapp.CreateUserInput{
		Name:     "Administrator",
		Email:    cfg.Security.AdminEmail,
		Password: cfg.Security.AdminPassword,
		Role:     userdomain.RoleAdmin,
	})
	return err
}
