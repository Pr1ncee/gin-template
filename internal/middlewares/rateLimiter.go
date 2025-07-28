package middlewares

import (
	"GinBox/config"
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

type RateLimiter struct {
	rateLimiters sync.Map
	cfg          *config.Config
}

// NewRateLimiter initializes the middleware with the allowed requests per second
func NewRateLimiter(cfg *config.Config) *RateLimiter {
	return &RateLimiter{cfg: cfg, rateLimiters: sync.Map{}}
}

func (rl *RateLimiter) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid remote address"})
			return
		}

		val, exists := rl.rateLimiters.Load(ip)
		var limiter ratelimit.Limiter
		if !exists {
			limiter = ratelimit.New(rl.cfg.App.MaxRPS) // Limit in seconds
			rl.rateLimiters.Store(ip, limiter)
		} else {
			limiter = val.(ratelimit.Limiter)
		}

		limiter.Take()

		c.Next()
	}
}
