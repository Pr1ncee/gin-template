package services

import (
	"GinBox/config"
	"GinBox/internal/handlers/requests"
	db "GinBox/internal/postgresql"
	"GinBox/internal/repositories"
	"GinBox/internal/utils"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type IUserService interface {
	CreateUser(ctx context.Context, data requests.CreateUserRequest) (db.User, error)
	ListUsers(ctx context.Context, args db.ListUsersParams) ([]db.ListUsersRow, error)
	GetUser(ctx context.Context, id int32) (db.GetUserByIDRow, error)
	UpdateUser(ctx context.Context, id int32, data requests.UpdateUserRequest) (db.User, int, error)
	DeleteUser(ctx context.Context, id int32) (int, error)
	ExportUsers(ctx context.Context, args db.ListUsersParams) (bytes.Buffer, string, int, error)
}

type UserService struct {
	queries  *db.Queries
	logger   *zap.Logger
	userRepo repositories.IUserRepository
	cfg      config.Config
}

func NewUserService(queries *db.Queries, logger *zap.Logger, cfg config.Config, userRepo repositories.IUserRepository) *UserService {
	return &UserService{queries: queries, logger: logger, cfg: cfg, userRepo: userRepo}
}

func (u *UserService) CreateUser(ctx context.Context, data requests.CreateUserRequest) (db.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), u.cfg.App.PasswordCost)

	phoneNumber := pgtype.Text{}
	if data.PhoneNumber != "" {
		phoneNumber = pgtype.Text{String: data.PhoneNumber, Valid: true}
	}

	user := db.User{
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		Email:       data.Email,
		Password:    hashedPassword,
		PhoneNumber: phoneNumber,
		Role:        utils.ParseUserRole(data.Role, u.logger),
	}

	createdUser, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		u.logger.Error("Failed to create user", zap.Error(err))
		return createdUser, err
	}
	createdUser.Password = nil
	return createdUser, nil
}

func (u *UserService) ListUsers(ctx context.Context, args db.ListUsersParams) ([]db.ListUsersRow, error) {
	users, err := u.userRepo.ListUsers(ctx, db.ListUsersParams{
		Limit:   args.Limit,
		Offset:  args.Offset,
		Column1: args.Column1,
		Column2: args.Column2,
	})
	return users, err
}

func (u *UserService) GetUser(ctx context.Context, id int32) (db.GetUserByIDRow, error) {
	user, err := u.userRepo.GetUserById(ctx, id)
	return user, err
}

func (u *UserService) UpdateUser(ctx context.Context, id int32, data requests.UpdateUserRequest) (db.User, int, error) {
	exists, err := u.userRepo.CheckUserExists(ctx, id)
	if err != nil {
		u.logger.Error("Failed to check user", zap.Error(err))
		return db.User{}, http.StatusInternalServerError, err
	}
	if !exists {
		errMsg := "user not found"
		u.logger.Warn(errMsg)
		return db.User{}, http.StatusNotFound, errors.New(errMsg)
	}

	params := db.UpdateUserPartialParams{
		ID: id,
	}

	if data.FirstName != "" {
		params.FirstName = pgtype.Text{String: data.FirstName, Valid: true}
	}
	if data.LastName != "" {
		params.LastName = pgtype.Text{String: data.LastName, Valid: true}
	}
	if data.Email != "" {
		params.Email = pgtype.Text{String: data.Email, Valid: true}
	}
	if data.PhoneNumber != "" {
		params.PhoneNumber = pgtype.Text{String: data.PhoneNumber, Valid: true}
	}
	if data.CurrentPassword != "" && data.NewPassword != "" {
		user, err := u.userRepo.GetUserWithPasswordById(ctx, params.ID) // TODO Move to Auth service
		if err != nil {
			u.logger.Error("Failed to get user with password", zap.Error(err))
			return db.User{}, http.StatusInternalServerError, err
		}
		errPass := bcrypt.CompareHashAndPassword(user.Password, []byte(data.CurrentPassword))
		if errPass != nil {
			errMsg := "current password doesn't match"
			u.logger.Error(errMsg)
			return db.User{}, http.StatusForbidden, errors.New(errMsg)
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
		params.Password = hashedPassword
	}
	if data.Role != "" {
		role := utils.ParseUserRole(data.Role, u.logger)
		params.Role = db.NullUserRole{UserRole: role, Valid: true}
	}

	updatedUser, err := u.userRepo.UpdateUserPartial(ctx, params)
	if err != nil {
		u.logger.Error("Failed to update user", zap.Error(err))
		return db.User{}, http.StatusInternalServerError, err
	}

	updatedUser.Password = nil
	return updatedUser, http.StatusAccepted, nil
}

func (u *UserService) DeleteUser(ctx context.Context, id int32) (int, error) {
	exists, err := u.userRepo.CheckUserExists(ctx, id)
	if err != nil {
		u.logger.Error("Failed to check user", zap.Error(err))
		return http.StatusInternalServerError, err
	}
	if !exists {
		errMsg := "user not found"
		u.logger.Warn(errMsg)
		return http.StatusNotFound, errors.New(errMsg)
	}

	err = u.userRepo.DeleteUserById(ctx, id)
	if err != nil {
		u.logger.Error("Failed to delete user", zap.Error(err))
		return http.StatusInternalServerError, err
	}
	return http.StatusAccepted, nil
}

func (u *UserService) ExportUsers(ctx context.Context, args db.ListUsersParams) (bytes.Buffer, string, int, error) {
	users, err := u.userRepo.ListUsers(ctx, db.ListUsersParams{
		Limit:   args.Limit,
		Offset:  args.Offset,
		Column1: args.Column1,
		Column2: args.Column2,
	})
	u.logger.Info("Got users to export", zap.Int("numOfUsers", len(users)))
	if err != nil {
		u.logger.Error("Failed to list users", zap.Error(err))
		return bytes.Buffer{}, "", http.StatusInternalServerError, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"ID", "First Name", "Last Name", "Email", "Phone Number", "Role", "Created At", "Updated At"}
	if err := writer.Write(header); err != nil {
		u.logger.Error("Failed to write header", zap.Error(err))
		return bytes.Buffer{}, "", http.StatusInternalServerError, err
	}

	for _, user := range users {
		record := []string{
			fmt.Sprintf("%d", user.ID),
			user.FirstName,
			user.LastName,
			user.Email,
			user.PhoneNumber.String,
			string(user.Role),
			utils.TimestamptzToString(user.CreatedAt, "2006-01-02 15:04:05", u.logger),
			utils.TimestamptzToString(user.UpdatedAt, "2006-01-02 15:04:05", u.logger),
		}

		if err := writer.Write(record); err != nil {
			u.logger.Error("Failed to write record", zap.Error(err))
			return bytes.Buffer{}, "", http.StatusInternalServerError, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		u.logger.Error("Failed to flush record", zap.Error(err))
		return bytes.Buffer{}, "", http.StatusInternalServerError, err
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("users_export_%s.csv", timestamp)
	u.logger.Info(
		"Exporting users to csv",
		zap.String("filename", filename),
		zap.Int("fileSize", len(buf.Bytes())))

	return buf, filename, http.StatusOK, nil
}
