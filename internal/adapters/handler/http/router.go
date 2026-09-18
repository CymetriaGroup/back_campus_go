package http

import (
	"hexagonal-go-backend/internal/adapters/handler/http/middleware"
	"hexagonal-go-backend/internal/config"
	"hexagonal-go-backend/internal/core/domain"
	"hexagonal-go-backend/internal/core/ports"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, l *slog.Logger, users *UserHandler, auth *AuthHandler, tokens ports.TokenProvider) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Recovery(l), middleware.Logger(l), middleware.CORS(cfg.App.AllowedOrigins), middleware.SecurityHeaders(), middleware.BodyLimit(1<<20), middleware.RateLimit(cfg.Security.RateLimit, time.Minute))
	r.GET("/health", Health)
	r.GET("/ready", Ready)
	v1 := r.Group("/api/v1")
	a := v1.Group("/auth")
	a.POST("/login", auth.Login)
	a.POST("/refresh", auth.Refresh)
	a.POST("/logout", auth.Logout)
	protected := v1.Group("")
	protected.Use(middleware.Authenticate(tokens))
	protected.GET("/users/:id", users.Get)
	admin := protected.Group("")
	admin.Use(middleware.RequireRole(domain.RoleAdmin))
	admin.POST("/users", users.Create)
	admin.GET("/users", users.List)
	admin.PUT("/users/:id", users.Update)
	admin.DELETE("/users/:id", users.Delete)
	return r
}
