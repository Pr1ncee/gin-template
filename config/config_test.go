package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
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

	require.Equal(t, "debug", cfg.App.Mode)
	require.Equal(t, 8080, cfg.App.Port)
	require.Equal(t, 30*time.Second, cfg.App.ReadTimeout)
	require.Equal(t, 30*time.Second, cfg.App.WriteTimeout)
	require.Equal(t, 30*time.Second, cfg.App.IdleTimeout)
	require.Equal(t, 30*time.Second, cfg.App.RequestTimeout)
	require.Equal(t, 1000, cfg.App.MaxProcs)

	require.Equal(t, "postgres://admin:admin@localhost:5432/gin-box", cfg.Postgres.ConnString)

	require.Equal(t, logrus.InfoLevel, cfg.Log.Level)
}
