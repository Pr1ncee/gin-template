package internal

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

func SetUpDBConn(ctx context.Context, connString string) *pgx.Conn {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}
	return conn
}
