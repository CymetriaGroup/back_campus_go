package middleware

import (
	"net/http"
	"strings"

	auth "hexagonal-go-backend/internal/modules/auth/application"
	userdomain "hexagonal-go-backend/internal/modules/users/domain"

	"github.com/gin-gonic/gin"
)

func Authenticate(tokens auth.TokenProvider) gin.HandlerFunc {
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
		c.Set("subject", claims.Subject)
		c.Set("role", string(claims.Role))
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

func unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Missing or invalid token"}, "request_id": c.GetString("request_id")})
}
