package services

import (
	"GinBox/config"
	"GinBox/internal/appErrors"
	"GinBox/internal/handlers/requests"
	db "GinBox/internal/postgresql"
	"GinBox/internal/repositories"
	"GinBox/internal/utils"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"time"
)

type IUserService interface {
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	Login(ctx context.Context, request requests.LoginRequest) (string, string, error)
	CreateUser(ctx context.Context, data requests.CreateUserRequest) (db.User, error)
	ListUsers(ctx context.Context, args db.ListUsersParams) ([]db.ListUsersRow, error)
	GetUser(ctx context.Context, id int32) (db.GetUserByIDRow, error)
	UpdateUser(ctx context.Context, id int32, data requests.UpdateUserRequest) (db.User, error)
	DeleteUser(ctx context.Context, id int32) error
	ExportUsers(ctx context.Context, args db.ListUsersParams) (bytes.Buffer, string, error)
}

type UserService struct {
	queries     *db.Queries
	logger      *zap.Logger
	userRepo    repositories.IUserRepository
	authService AuthService
	cfg         config.Config
}

func NewUserService(queries *db.Queries, logger *zap.Logger, cfg config.Config, userRepo repositories.IUserRepository, authService AuthService) *UserService {
	return &UserService{queries: queries, logger: logger, cfg: cfg, userRepo: userRepo, authService: authService}
}

func (u *UserService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := u.authService.ParseToken(refreshToken)
	if err != nil {
		u.logger.Error("Error parsing the token", zap.Error(err))
		return "", appErrors.ErrInvalidToken
	}
	user, err := u.GetUser(ctx, claims.UserID)
	if err != nil {
		u.logger.Error("User not found", zap.Error(err))
		return "", appErrors.ErrUserNotFound
	}

	newRefreshToken, err := u.authService.GenerateAccessToken(UserClaims{UserID: user.ID, Email: user.Email, Role: user.Role})
	if err != nil {
		u.logger.Error("Error generating access token", zap.Error(err))
		return "", appErrors.ErrFailedToGenerateAccessToken
	}
	return newRefreshToken, nil
}

func (u *UserService) Login(ctx context.Context, data requests.LoginRequest) (string, string, error) {
	user, err := u.GetUserWithPasswordById(ctx, data.ID)
	if err != nil {
		u.logger.Warn("User not found", zap.Error(err))
		return "", "", appErrors.ErrUserNotFound
	}
	errPass := u.authService.ValidatePassword(user.Password, []byte(data.Password))
	if errPass != nil {
		u.logger.Error("Password does not match", zap.Int32("user_id", data.ID))
		return "", "", appErrors.ErrPasswordDoesNotMatch
	}

	userClaims := UserClaims{UserID: user.ID, Email: user.Email, Role: user.Role}
	accessToken, err := u.authService.GenerateAccessToken(userClaims)
	if err != nil {
		u.logger.Warn("Failed to generate access token", zap.Error(err))
		return "", "", appErrors.ErrFailedToGenerateAccessToken
	}
	refreshToken, err := u.authService.GenerateRefreshToken(userClaims)
	if err != nil {
		u.logger.Warn("Failed to generate refresh token", zap.Error(err))
		return "", "", appErrors.ErrFailedToGenerateRefreshToken
	}
	return accessToken, refreshToken, nil
}

func (u *UserService) CreateUser(ctx context.Context, data requests.CreateUserRequest) (db.User, error) {
	hashedPassword, err := u.authService.GeneratePassword([]byte(data.Password))

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
		return createdUser, appErrors.ErrFailedToCreateUser
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

func (u *UserService) GetUserWithPasswordById(ctx context.Context, id int32) (db.User, error) {
	user, err := u.userRepo.GetUserWithPasswordById(ctx, id)
	return user, err
}

func (u *UserService) UpdateUser(ctx context.Context, id int32, data requests.UpdateUserRequest) (db.User, error) {
	exists, err := u.userRepo.CheckUserExists(ctx, id)
	if err != nil {
		u.logger.Error("Failed to check user", zap.Error(err))
		return db.User{}, err
	}
	if !exists {
		u.logger.Warn("User not found")
		return db.User{}, appErrors.ErrUserNotFound
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
		user, err := u.userRepo.GetUserWithPasswordById(ctx, params.ID)
		if err != nil {
			u.logger.Error("Failed to get user with password", zap.Error(err))
			return db.User{}, err
		}
		errPass := u.authService.ValidatePassword(user.Password, []byte(data.NewPassword))
		if errPass != nil {
			u.logger.Error("Current password doesn't match")
			return db.User{}, appErrors.ErrPasswordDoesNotMatch
		}
		hashedPassword, err := u.authService.GeneratePassword([]byte(data.NewPassword))
		if err != nil {
			u.logger.Error("Failed to generate password", zap.Error(err))
			return db.User{}, err
		}
		params.Password = hashedPassword
	}
	if data.Role != "" {
		role := utils.ParseUserRole(data.Role, u.logger)
		params.Role = db.NullUserRole{UserRole: role, Valid: true}
	}

	updatedUser, err := u.userRepo.UpdateUserPartial(ctx, params)
	if err != nil {
		u.logger.Error("Failed to update user", zap.Error(err))
		return db.User{}, err
	}

	updatedUser.Password = nil
	return updatedUser, nil
}

func (u *UserService) DeleteUser(ctx context.Context, id int32) error {
	exists, err := u.userRepo.CheckUserExists(ctx, id)
	if err != nil {
		u.logger.Error("Failed to check user", zap.Error(err))
		return err
	}
	if !exists {
		u.logger.Warn("User not found")
		return appErrors.ErrUserNotFound
	}

	err = u.userRepo.DeleteUserById(ctx, id)
	if err != nil {
		u.logger.Error("Failed to delete user", zap.Error(err))
		return err
	}
	return nil
}

func (u *UserService) ExportUsers(ctx context.Context, args db.ListUsersParams) (bytes.Buffer, string, error) {
	users, err := u.userRepo.ListUsers(ctx, db.ListUsersParams{
		Limit:   args.Limit,
		Offset:  args.Offset,
		Column1: args.Column1,
		Column2: args.Column2,
	})
	u.logger.Info("Got users to export", zap.Int("numOfUsers", len(users)))
	if err != nil {
		u.logger.Error("Failed to list users", zap.Error(err))
		return bytes.Buffer{}, "", err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	header := []string{"ID", "First Name", "Last Name", "Email", "Phone Number", "Role", "Created At", "Updated At"}
	if err := writer.Write(header); err != nil {
		u.logger.Error("Failed to write header", zap.Error(err))
		return bytes.Buffer{}, "", err
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
			return bytes.Buffer{}, "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		u.logger.Error("Failed to flush record", zap.Error(err))
		return bytes.Buffer{}, "", err
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("users_export_%s.csv", timestamp)
	u.logger.Info(
		"Exporting users to csv",
		zap.String("filename", filename),
		zap.Int("fileSize", len(buf.Bytes())))

	return buf, filename, nil
}
