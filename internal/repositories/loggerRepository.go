/*
Package repositories provides low-level functionality with particular tools such as Redis, PostgreSQL, etc.

Specifically, this file creates a custom Logger object via zap.
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

// SetUpLogger creates a zap.Logger object with custom parameters and settings and returns it
// implying using it further in different locations across the application to log actions.
func (a *AppLogger) SetUpLogger(cfg *config.Config) *zap.Logger {
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

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		cfg.Log.Level,
	)

	logger := zap.New(core, zap.AddCaller())

	if cfg.App.Mode != gin.ReleaseMode {
		logger = logger.WithOptions(
			zap.Development(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	}

	return logger
}
