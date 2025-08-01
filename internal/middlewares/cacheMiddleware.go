package middlewares

import (
	"GinBox/config"
	"GinBox/internal/repositories"
	"GinBox/internal/utils"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"go.uber.org/zap"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CacheMiddleware represents cache middleware with Redis a in-memory storage.
// It also accepts logger and config as input fields
type CacheMiddleware struct {
	repo   repositories.IRedisRepository
	cfg    *config.Config
	logger *zap.Logger
}

// NewCacheMiddleware creates a new CacheMiddleware object.
func NewCacheMiddleware(repo repositories.IRedisRepository, cfg *config.Config, logger *zap.Logger) *CacheMiddleware {
	return &CacheMiddleware{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

// Handle implements the primary functionality of caching endpoints' responses.
// This middleware ignores all endpoints with HTTP GET method.
func (m *CacheMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		ctx := context.Background()
		cacheKey := m.generateCacheKey(c.Request)

		cached, err := m.repo.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			m.logger.Info("Cache hit", zap.String("key", cacheKey))
			c.Data(http.StatusOK, "application/json", []byte(cached))
			c.Abort()
			return
		}

		writer := &utils.ResponseBodyWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		if c.Writer.Status() >= http.StatusOK && c.Writer.Status() <= http.StatusIMUsed {
			m.logger.Info("Setting cache", zap.String("key", cacheKey))
			_ = m.repo.Set(ctx, cacheKey, writer.Body.String(), int64(m.cfg.Redis.TTL.Seconds()))
		}
	}
}

// generateCacheKey implements a method for generating a cache key in a hash format to be saved to Redis.
func (m *CacheMiddleware) generateCacheKey(req *http.Request) string {
	hash := sha1.New()
	io.WriteString(hash, req.Method+"|"+req.URL.String())
	return "cache:" + hex.EncodeToString(hash.Sum(nil))
}
