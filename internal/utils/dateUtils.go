package utils

import "github.com/jackc/pgx/v5/pgtype"

func TimestamptzToString(ts pgtype.Timestamptz, layout string) string {
	if !ts.Valid {
		return ""
	}
	return ts.Time.Format(layout)
}
