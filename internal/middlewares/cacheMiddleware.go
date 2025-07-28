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

type CacheMiddleware struct {
	repo   *repositories.RedisRepository
	cfg    *config.Config
	logger *zap.Logger
}

func NewCacheMiddleware(repo *repositories.RedisRepository, cfg *config.Config, logger *zap.Logger) *CacheMiddleware {
	return &CacheMiddleware{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (m *CacheMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		ctx := context.Background()
		cacheKey := m.GenerateCacheKey(c.Request)

		cached, err := m.repo.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			m.logger.Info("cache hit", zap.String("key", cacheKey))
			c.Data(http.StatusOK, "application/json", []byte(cached))
			c.Abort()
			return
		}

		writer := &utils.ResponseBodyWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		if c.Writer.Status() >= http.StatusOK || c.Writer.Status() <= http.StatusIMUsed {
			m.logger.Info("setting cache", zap.String("key", cacheKey))
			_ = m.repo.Set(ctx, cacheKey, writer.Body.String(), int64(m.cfg.Redis.TTL.Seconds()))
		}
	}
}

func (m *CacheMiddleware) GenerateCacheKey(req *http.Request) string {
	hash := sha1.New()
	io.WriteString(hash, req.Method+"|"+req.URL.String())
	return "cache:" + hex.EncodeToString(hash.Sum(nil))
}
