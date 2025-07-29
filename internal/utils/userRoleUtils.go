package utils

import (
	db "GinBox/internal/postgresql"
	"go.uber.org/zap"
)

// ParseUserRole converts plain string values into UserRole type with the corresponding value.
func ParseUserRole(role string, logger *zap.Logger) db.UserRole {
	// TODO think about this function, whether it's needed or not, pass logger here
	switch role {
	case "Admin":
		return db.UserRoleAdmin
	case "User":
		return db.UserRoleUser
	default:
		logger.Warn("unknown role", zap.String("role", role))
		return ""
	}
}
