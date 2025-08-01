package config

import (
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoadConfigWithFile(t *testing.T) {
	testFile := "config.test.env"
	testPath := "./testconfig"
	fullPath := filepath.Join(testPath, testFile)

	_, err := os.Stat(fullPath)
	require.NoError(t, err, "Config file not found")

	cfg, err := LoadConfigWithFile(testPath, "config.test")
	require.NoError(t, err)

	require.Equal(t, "gin-box", cfg.App.Name)
	require.Equal(t, "debug", cfg.App.Mode)
	require.Equal(t, 8080, cfg.App.Port)
	require.Equal(t, 30*time.Second, cfg.App.ReadTimeout)
	require.Equal(t, 30*time.Second, cfg.App.WriteTimeout)
	require.Equal(t, 30*time.Second, cfg.App.IdleTimeout)
	require.Equal(t, 30*time.Second, cfg.App.RequestTimeout)
	require.Equal(t, 1000, cfg.App.MaxProcs)
	require.Equal(t, []string{"*"}, cfg.App.Origins)
	require.Equal(t, 10, cfg.App.MaxRPS)

	require.Equal(t, 10, cfg.Auth.PasswordCost)
	require.Equal(t, "secret_key", cfg.Auth.JWTSecret)
	require.Equal(t, 3600*time.Second, cfg.Auth.AccessTokenTTL)
	require.Equal(t, 36000*time.Second, cfg.Auth.RefreshTokenTTL)

	require.Equal(t, "postgres://admin:admin@postgres:5432/gin-box", cfg.Postgres.ConnString)

	require.Equal(t, zapcore.InfoLevel, cfg.Log.Level)
	require.Equal(t, "http://loki:3100", cfg.Log.LokiURL)

	require.Equal(t, "redis:6379", cfg.Redis.Host)
	require.Equal(t, "", cfg.Redis.Password)
	require.Equal(t, 0, cfg.Redis.DB)
	require.Equal(t, 60*time.Second, cfg.Redis.TTL)
}
