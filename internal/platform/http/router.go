package http

import (
	"log/slog"
	"time"

	authapp "hexagonal-go-backend/internal/modules/auth/application"
	authhttp "hexagonal-go-backend/internal/modules/auth/delivery/http/v1"
	tenanthttp "hexagonal-go-backend/internal/modules/tenants/delivery/http/v1"
	usershttp "hexagonal-go-backend/internal/modules/users/delivery/http/v1"
	userdomain "hexagonal-go-backend/internal/modules/users/domain"
	"hexagonal-go-backend/internal/platform/config"
	"hexagonal-go-backend/internal/platform/http/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg config.Config, logger *slog.Logger, users *usershttp.Controller, auth *authhttp.Controller, tokens authapp.TokenProvider, tenants ...*tenanthttp.Controller) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middleware.RequestID(), middleware.Recovery(logger), middleware.Logger(logger), middleware.CORS(cfg.App.AllowedOrigins), middleware.SecurityHeaders(), middleware.BodyLimit(1<<20), middleware.RateLimit(cfg.Security.RateLimit, time.Minute))
	router.GET("/health", Health)
	router.GET("/ready", Ready)

	v1 := router.Group("/api/v1")
	authRoutes := v1.Group("/auth")
	authRoutes.POST("/login", auth.Login)
	authRoutes.POST("/refresh", auth.Refresh)
	authRoutes.POST("/logout", auth.Logout)

	if len(tenants) > 0 && tenants[0] != nil {
		tenants[0].RegisterRoutes(v1)
	}

	protected := v1.Group("")
	protected.Use(middleware.Authenticate(tokens))
	protected.GET("/users/:id", users.Get)

	admin := protected.Group("")
	admin.Use(middleware.RequireRole(userdomain.RoleAdmin))
	admin.POST("/users", users.Create)
	admin.GET("/users", users.List)
	admin.PUT("/users/:id", users.Update)
	admin.DELETE("/users/:id", users.Delete)
	return router
}

