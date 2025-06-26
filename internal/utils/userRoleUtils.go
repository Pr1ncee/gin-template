package utils

import db "GinBox/internal/postgresql"

func ParseUserRole(role string) db.UserRole {
	// TODO think about this function, whether it's needed or not, pass logger here
	switch role {
	case "Admin":
		return db.UserRoleAdmin
	case "User":
		return db.UserRoleUser
	default:
		return ""
	}
}
