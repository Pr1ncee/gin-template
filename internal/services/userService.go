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
	log "github.com/sirupsen/logrus"
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
	logger   *log.Logger
	userRepo repositories.IUserRepository
	cfg      config.Config
}

func NewUserService(queries *db.Queries, logger *log.Logger, cfg config.Config, userRepo repositories.IUserRepository) *UserService {
	return &UserService{queries: queries, logger: logger, cfg: cfg, userRepo: userRepo}
}

func (u *UserService) CreateUser(ctx context.Context, data requests.CreateUserRequest) (db.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost) // TODO move to env vars

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
		Role:        utils.ParseUserRole(data.Role),
	}

	createdUser, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
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
		return db.User{}, http.StatusInternalServerError, err
	}
	if !exists {
		return db.User{}, http.StatusNotFound, errors.New("user not found")
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
			return db.User{}, http.StatusInternalServerError, err
		}
		errPass := bcrypt.CompareHashAndPassword(user.Password, []byte(data.CurrentPassword))
		if errPass != nil {
			return db.User{}, http.StatusForbidden, errors.New("current password doesn't match")
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
		params.Password = hashedPassword
	}
	if data.Role != "" {
		role := utils.ParseUserRole(data.Role)
		params.Role = db.NullUserRole{UserRole: role, Valid: true}
	}

	updatedUser, err := u.userRepo.UpdateUserPartial(ctx, params)
	if err != nil {
		return db.User{}, http.StatusInternalServerError, err
	}

	updatedUser.Password = nil
	return updatedUser, http.StatusAccepted, nil
}

func (u *UserService) DeleteUser(ctx context.Context, id int32) (int, error) {
	exists, err := u.userRepo.CheckUserExists(ctx, id)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !exists {
		return http.StatusNotFound, errors.New("user not found")
	}

	err = u.userRepo.DeleteUserById(ctx, id)
	if err != nil {
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
	if err != nil {
		return bytes.Buffer{}, "", http.StatusInternalServerError, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"ID", "First Name", "Last Name", "Email", "Phone Number", "Role", "Created At", "Updated At"}
	if err := writer.Write(header); err != nil {
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
			utils.TimestamptzToString(user.CreatedAt, "2006-01-02 15:04:05"),
			utils.TimestamptzToString(user.UpdatedAt, "2006-01-02 15:04:05"),
		}

		if err := writer.Write(record); err != nil {
			return bytes.Buffer{}, "", http.StatusInternalServerError, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return bytes.Buffer{}, "", http.StatusInternalServerError, err
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("users_export_%s.csv", timestamp)

	return buf, filename, http.StatusOK, nil
}
