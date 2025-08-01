package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// SetUpDBConn setups database connection with PostgreSQL using input connString and context.
func SetUpDBConn(ctx context.Context, connString string, logger *zap.Logger) *pgx.Conn {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		logger.Fatal("Error connecting to database", zap.Error(err))
	}
	return conn
}
