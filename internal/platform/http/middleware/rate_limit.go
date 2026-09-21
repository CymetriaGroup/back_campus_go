package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count  int
	window time.Time
}

func RateLimit(max int, period time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	clients := map[string]visitor{}
	return func(c *gin.Context) {
		mu.Lock()
		v := clients[c.ClientIP()]
		now := time.Now()
		if now.Sub(v.window) >= period {
			v = visitor{window: now}
		}
		v.count++
		clients[c.ClientIP()] = v
		blocked := v.count > max
		mu.Unlock()
		if blocked {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": gin.H{"code": "RATE_LIMITED", "message": "Too many requests"}, "request_id": c.GetString("request_id")})
			return
		}
		c.Next()
	}
}
