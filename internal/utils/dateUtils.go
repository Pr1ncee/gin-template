package utils

import (
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func TimestamptzToString(ts pgtype.Timestamptz, layout string, logger *zap.Logger) string {
	if !ts.Valid {
		logger.Warn("timestamp not valid")
		return ""
	}
	return ts.Time.Format(layout)
}
