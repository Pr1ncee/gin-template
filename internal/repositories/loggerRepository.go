/*
Package repositories provides low-level functionality with particular tools such as Redis, PostgreSQL, etc.

Specifically, this file creates a custom Logger object via zap and connects to Loki.
*/
package repositories

import (
	"GinBox/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// AppLogger represents an object that should be used across the whole application to log actions.
type AppLogger struct {
	cfg *config.Config
}

// NewAppLogger initializes the repository with input config.
func NewAppLogger(cfg *config.Config) *AppLogger {
	return &AppLogger{cfg: cfg}
}

// SetUpLogger creates a zap.Logger object to log actions across the application.
// It also creates a connection with Loki to send all zap logs there and visualize them with Grafana.
func (a *AppLogger) SetUpLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Console logger
	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		a.cfg.Log.Level,
	)

	// Loki client
	lokiClient := NewLokiClient(
		a.cfg.Log.LokiURL,
		map[string]string{
			"app": a.cfg.App.Name,
			"env": a.cfg.App.Mode,
		},
	)
	lokiSyncer := NewLokiSyncer(lokiClient)

	lokiCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(lokiSyncer),
		a.cfg.Log.Level,
	)

	core := zapcore.NewTee(consoleCore, lokiCore)

	logger := zap.New(core, zap.AddCaller())

	if a.cfg.App.Mode != gin.ReleaseMode {
		logger = logger.WithOptions(
			zap.Development(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	}

	return logger
}
