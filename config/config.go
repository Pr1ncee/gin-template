/*
Package config provides a structured application configuration loading using Viper.

This package defines typed configuration structures and loads configuration from a YAML file.
It supports environment variable overrides and provides default values for some fields.

The configuration includes settings for:
  - Application parameters (port, timeout, max threads)
  - Authentication(JWT settings & password cost)
  - PostgreSQL connection string
  - Redis connection settings (host, password, db and TTL)
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
APP_ORIGINS=*
APP_MAX_RPS=10

AUTH_PASSWORD_COST=10 // Min Cost is 4, Max Cost is 31
AUTH_JWT_SECRET_KEY=secret_key
AUTH_JWT_ACCESS_TOKEN_TTL=3600s
AUTH_JWT_REFRESH_TOKEN_TTL=36000s
AUTH_JWT_ISSUER=gin-box

POSTGRES_CONN_STRING=postgres://admin:admin@postgres:5432/gin-box

REDIS_HOST=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
TTL=60s

LOG_LEVEL=0 // 5 = fatal, 4 = panic, 3 = dpanic, 2 = error, 1 = warn, 0 = info, -1 = debug
LOG_LOKI_URL=http://loki:3100

// The settings below are ONLY used internally by GOOSE
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://admin:admin@localhost:5432/gin-box
GOOSE_MIGRATION_DIR=./migrations
```
*/
package config

import (
	"go.uber.org/zap/zapcore"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Auth     AuthConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Log      LogConfig
}

type AppConfig struct {
	Name           string        `mapstructure:"name"`
	Mode           string        `mapstructure:"mode"`
	Port           int           `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	MaxProcs       int           `mapstructure:"max_procs"`
	Origins        []string      `mapstructure:"origins"`
	MaxRPS         int           `mapstructure:"max_rps"`
}

type AuthConfig struct {
	PasswordCost    int           `mapstructure:"password_cost"`
	JWTSecret       string        `mapstructure:"jwt_secret"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

type PostgresConfig struct {
	ConnString string `mapstructure:"conn_string"`
}

type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	TTL      time.Duration `mapstructure:"ttl"`
}

type LogConfig struct {
	Level   zapcore.Level `mapstructure:"level"`
	LokiURL string        `mapstructure:"loki_url"`
}

// LoadConfigWithFile loads the configuration from a YAML file located at the given path with the provided filename.
// Environment variables override YAML file fields automatically using underscore notation.
func LoadConfigWithFile(path, filename string) (*Config, error) {
	viper.SetConfigName(filename)
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("APP_TIMEOUT", "30s")
	viper.SetDefault("APP_MAX_RPS", 10)

	viper.SetDefault("LOG_LEVEL", zapcore.InfoLevel)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	cfg.App.Name = viper.GetString("APP_NAME")
	cfg.App.Mode = viper.GetString("APP_MODE")
	cfg.App.Port = viper.GetInt("APP_PORT")
	cfg.App.ReadTimeout = viper.GetDuration("APP_READ_TIMEOUT")
	cfg.App.WriteTimeout = viper.GetDuration("APP_WRITE_TIMEOUT")
	cfg.App.IdleTimeout = viper.GetDuration("APP_IDLE_TIMEOUT")
	cfg.App.RequestTimeout = viper.GetDuration("APP_REQUEST_TIMEOUT")
	cfg.App.MaxProcs = viper.GetInt("APP_MAX_PROCS")
	cfg.App.Origins = viper.GetStringSlice("APP_ORIGINS")
	cfg.App.MaxRPS = viper.GetInt("APP_MAX_RPS")

	cfg.Auth.PasswordCost = viper.GetInt("AUTH_PASSWORD_COST")
	cfg.Auth.JWTSecret = viper.GetString("AUTH_JWT_SECRET_KEY")
	cfg.Auth.AccessTokenTTL = viper.GetDuration("AUTH_JWT_ACCESS_TOKEN_TTL")
	cfg.Auth.RefreshTokenTTL = viper.GetDuration("AUTH_JWT_REFRESH_TOKEN_TTL")

	cfg.Postgres.ConnString = viper.GetString("POSTGRES_CONN_STRING")

	cfg.Redis.Host = viper.GetString("REDIS_HOST")
	cfg.Redis.Password = viper.GetString("REDIS_PASSWORD")
	cfg.Redis.DB = viper.GetInt("REDIS_DB")
	cfg.Redis.TTL = viper.GetDuration("REDIS_TTL")

	cfg.Log.Level = zapcore.Level(viper.GetInt("LOG_LEVEL"))
	cfg.Log.LokiURL = viper.GetString("LOG_LOKI_URL")

	return &cfg, nil
}
