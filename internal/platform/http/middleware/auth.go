package middleware

import (
	"net/http"
	"strings"

	auth "hexagonal-go-backend/internal/modules/auth/application"
	usersapp "hexagonal-go-backend/internal/modules/users/application"
	userdomain "hexagonal-go-backend/internal/modules/users/domain"

	"github.com/gin-gonic/gin"
)

func Authenticate(tokens auth.TokenProvider, userRepo usersapp.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			unauthorized(c)
			return
		}
		claims, err := tokens.Parse(parts[1])
		if err != nil {
			unauthorized(c)
			return
		}

		userID := claims.Subject
		role := string(claims.Role)
		tenantID := claims.TenantID

		// Lectura únicamente: el usuario local lo provisiona el flujo de login.
		if userRepo != nil {
			ctx := c.Request.Context()
			user, err := userRepo.FindByAuthID(ctx, claims.Subject)
			if err != nil {
				// Tokens locales (HMAC) llevan el ID interno como subject.
				user, err = userRepo.FindByID(ctx, claims.Subject)
			}
			if err != nil || user == nil || !user.Active {
				unauthorized(c)
				return
			}
			userID = user.ID
			role = string(user.Role)
			tenantID = user.TenantID
		}

		c.Set("user_id", userID)
		c.Set("auth_id", claims.Subject)
		c.Set("subject", userID)
		c.Set("role", role)
		c.Set("tenant_id", tenantID)
		c.Next()
	}
}

func RequireRole(role userdomain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != string(role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "Insufficient permissions"}, "request_id": c.GetString("request_id")})
			return
		}
		c.Next()
	}
}

func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("tenant_id") == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "TENANT_REQUIRED", "message": "Tenant context required"}, "request_id": c.GetString("request_id")})
			return
		}
		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Missing or invalid token"}, "request_id": c.GetString("request_id")})
}
