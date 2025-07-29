package main

import (
	"GinBox/config"
	_ "GinBox/docs"
	"GinBox/internal/api"
	v1 "GinBox/internal/api/v1"
	"GinBox/internal/handlers"
	db "GinBox/internal/postgresql"
	"GinBox/internal/repositories"
	"GinBox/internal/services"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

// @title Gin Box API
// @version 1.0
// @description This is GinBox application intended to save your time at every project's setup routine!

// @contact.name   API Support
// @contact.url    https://linkedin.com/in/andrew-strongin-191b56243
// @contact.email  andrewvs0707@gmail.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token like: **Bearer YOUR_TOKEN_HERE**
func main() {
	cfg, err := config.LoadConfigWithFile("./", ".env")
	if err != nil {
		log.Fatalf("failed to load config %v", err)
	}
	runtime.GOMAXPROCS(cfg.App.MaxProcs)

	initLogger := repositories.NewAppLogger(cfg)
	logger := initLogger.SetUpLogger(cfg)
	logger.Info(
		"Starting Gin Box API server...",
		zap.String("Environment", cfg.App.Mode),
		zap.Int("Max processes", cfg.App.MaxProcs))

	gin.SetMode(cfg.App.Mode)

	ctx := context.Background()
	conn := db.SetUpDBConn(ctx, cfg.Postgres.ConnString, logger)
	defer func(dbConn *pgx.Conn, context context.Context) {
		if err := dbConn.Close(context); err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
		} else {
			logger.Info("Database connection closed successfully")
		}
	}(conn, ctx)
	queries := db.New(conn)

	redisRepo := repositories.NewRedisRepository(cfg)
	logger.Info("Database connections established successfully")

	userRepo := repositories.NewUserRepository(queries, logger)
	logger.Info("Repositories initialized")

	authService := services.NewAuthService(*cfg, logger)
	userService := services.NewUserService(userRepo, authService, logger, *cfg)
	logger.Info("Services initialized")

	healthCheckHandler := handlers.NewHealthCheckHandler(logger)

	userHandler := handlers.NewUserHandler(userService, logger, cfg.App.RequestTimeout)
	logger.Info("Handlers initialized")

	router := gin.New()
	apiV1 := v1.NewApiV1(userHandler, authService, healthCheckHandler, redisRepo, cfg, logger)
	logger.Info("ApiV1 initialized")
	apiGroup := api.NewApiGroup(userHandler, apiV1, cfg, logger)
	apiGroup.InitRouterGroups(router)
	logger.Info("Routes initialized")

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		Handler:           router,
		ReadTimeout:       cfg.App.ReadTimeout,
		WriteTimeout:      cfg.App.WriteTimeout,
		IdleTimeout:       cfg.App.IdleTimeout,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("Server starting on port %d", cfg.App.Port))
		logger.Info(fmt.Sprintf("Swagger documentation available at: http://localhost:%d/swagger/index.html", cfg.App.Port))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Info(fmt.Sprintf("Received signal: %v. Shutting down server...", sig))

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Server exited gracefully")
}
