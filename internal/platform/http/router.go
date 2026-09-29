package http

import (
	"log/slog"
	"time"

	authapp "hexagonal-go-backend/internal/modules/auth/application"
	authhttp "hexagonal-go-backend/internal/modules/auth/delivery/http/v1"
	courseshttp "hexagonal-go-backend/internal/modules/courses/delivery/http/v1"
	usershttp "hexagonal-go-backend/internal/modules/users/delivery/http/v1"
	userdomain "hexagonal-go-backend/internal/modules/users/domain"
	"hexagonal-go-backend/internal/platform/config"
	"hexagonal-go-backend/internal/platform/http/middleware"

	_ "hexagonal-go-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title                      Go Hexagonal API
// @version                    1.0
// @description                API REST construida con Go, Gin y Arquitectura Hexagonal.

// @contact.name              API Support
// @contact.email             soporte@empresa.com

// @host                      localhost:8080
// @BasePath                  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Escribe "Bearer " seguido de tu token JWT

func NewRouter(cfg config.Config, logger *slog.Logger, users *usershttp.Controller, auth *authhttp.Controller, tokens authapp.TokenProvider, courses ...*courseshttp.Controller) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(middleware.RequestID(), middleware.Recovery(logger), middleware.Logger(logger), middleware.CORS(cfg.App.AllowedOrigins), middleware.SecurityHeaders(), middleware.BodyLimit(1<<20), middleware.RateLimit(cfg.Security.RateLimit, time.Minute))
	router.GET("/health", Health)
	router.GET("/ready", Ready)

	v1 := router.Group("/api/v1")
	authRoutes := v1.Group("/auth")
	authRoutes.POST("/login", auth.Login)
	authRoutes.POST("/refresh", auth.Refresh)
	authRoutes.POST("/logout", auth.Logout)

	protected := v1.Group("")
	protected.Use(middleware.Authenticate(tokens))
	protected.GET("/users/:id", users.Get)

	admin := protected.Group("")
	admin.Use(middleware.RequireRole(userdomain.RoleAdmin))
	admin.POST("/users", users.Create)
	admin.GET("/users", users.List)
	admin.PUT("/users/:id", users.Update)
	admin.DELETE("/users/:id", users.Delete)
	if len(courses) > 0 && courses[0] != nil {
		catalog := v1.Group("/courses")
		catalog.POST("", courses[0].CreateCourse)
		catalog.GET("/:id", courses[0].GetCourse)
		catalog.GET("/categories", courses[0].ListCategories)
		catalog.GET("/templates", courses[0].ListTemplates)
		catalog.GET("/templates/:id", courses[0].GetTemplate)
		catalog.GET("/templates/:id/versions", courses[0].ListVersions)
		catalog.GET("/versions/:id/syllabus", courses[0].GetSyllabus)
		catalog.POST("/categories", courses[0].CreateCategory)
		catalog.POST("/templates", courses[0].CreateTemplate)
		catalog.POST("/templates/:id/versions", courses[0].CreateVersion)
		catalog.PUT("/versions/:id/syllabus", courses[0].SetSyllabus)
	}
	return router
}
