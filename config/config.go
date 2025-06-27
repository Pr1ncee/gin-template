/*
Package config provides a structured application configuration loading using Viper.

This package defines typed configuration structures and loads configuration from a YAML file.
It supports environment variable overrides and provides default values for some fields.

The configuration includes settings for:
  - Application parameters (port, timeout, max threads)
  - PostgreSQL connection string
  - Logging level

# Sample Configuration (.env)

```
APP_MODE=debug
APP_PORT=8080
APP_READ_TIMEOUT=30s
APP_WRITE_TIMEOUT=30s
APP_IDLE_TIMEOUT=30s
APP_REQUEST_TIMEOUT=30s
APP_MAX_PROCS=1000
APP_PASSWORD_COST=10 # Min Cost is 4, Max Cost is 31

POSTGRES_CONN_STRING=postgres://admin:admin@localhost:5432/gin-box

LOG_LEVEL=0 # 5 = fatal, 4 = panic, 3 = dpanic, 2 = error, 1 = warn, 0 = info, -1 = debug
*/
package config

import (
	"go.uber.org/zap/zapcore"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Log      LogConfig
}

type AppConfig struct {
	Mode           string        `mapstructure:"mode"`
	Port           int           `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	MaxProcs       int           `mapstructure:"max_procs"`
	PasswordCost   int           `mapstructure:"password_cost"`
}

type PostgresConfig struct {
	ConnString string `mapstructure:"conn_string"`
}

type LogConfig struct {
	Level zapcore.Level `mapstructure:"level"`
}

// LoadConfigWithFile loads the configuration from a YAML file located at the given path with the provided filename.
// Environment variables override YAML file fields automatically using underscore notation.
func LoadConfigWithFile(path, filename string) (*Config, error) {
	viper.SetConfigName(filename)
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("APP_TIMEOUT", "30s")
	viper.SetDefault("LOG_LEVEL", zapcore.InfoLevel)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	cfg.App.Mode = viper.GetString("APP_MODE")
	cfg.App.Port = viper.GetInt("APP_PORT")
	cfg.App.ReadTimeout = viper.GetDuration("APP_READ_TIMEOUT")
	cfg.App.WriteTimeout = viper.GetDuration("APP_WRITE_TIMEOUT")
	cfg.App.IdleTimeout = viper.GetDuration("APP_IDLE_TIMEOUT")
	cfg.App.RequestTimeout = viper.GetDuration("APP_REQUEST_TIMEOUT")
	cfg.App.PasswordCost = viper.GetInt("APP_PASSWORD_COST")

	cfg.App.MaxProcs = viper.GetInt("APP_MAX_PROCS")

	cfg.Postgres.ConnString = viper.GetString("POSTGRES_CONN_STRING")

	cfg.Log.Level = zapcore.Level(viper.GetInt("LOG_LEVEL"))

	return &cfg, nil
}
