package repositories

import (
	db "GinBox/internal/postgresql"
	"context"
	"fmt"
	"go.uber.org/zap"
)

// IUserRepository defines the interface for basic methods for User entity (specifically in PostgreSQL database).
type IUserRepository interface {
	CreateUser(ctx context.Context, user db.User) (db.User, error)
	ListUsers(ctx context.Context, args db.ListUsersParams) ([]db.ListUsersRow, error)
	GetUserById(ctx context.Context, id int32) (db.GetUserByIDRow, error)
	GetUserWithPasswordById(ctx context.Context, id int32) (db.User, error)
	CheckUserExists(ctx context.Context, id int32) (bool, error)
	UpdateUserPartial(ctx context.Context, arg db.UpdateUserPartialParams) (db.User, error)
	DeleteUserById(ctx context.Context, id int32) error
}

// UserRepository implements IUserRepository interface and provides high-level usable methods.
type UserRepository struct {
	queries *db.Queries
	logger  *zap.Logger
}

// NewUserRepository initializes a new repository with all available queries.
func NewUserRepository(queries *db.Queries, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		queries: queries,
		logger:  logger,
	}
}

// CreateUser implements User creation functionality with the input schema.
func (r *UserRepository) CreateUser(ctx context.Context, user db.User) (db.User, error) {
	r.logger.Info("Creating new user", zap.String("email", user.Email))

	params := db.CreateUserParams{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Password:    user.Password,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	}

	createdUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		r.logger.Error("Failed to create user", zap.Error(err), zap.String("email", user.Email))
		return db.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Info(
		"User created successfully",
		zap.Int32("user_id", createdUser.ID), zap.String("email", createdUser.Email))

	return createdUser, nil
}

// ListUsers returns a list of Users with the corresponding input parameters
// such as filtering by role and searching by first and last names and email.
func (r *UserRepository) ListUsers(ctx context.Context, args db.ListUsersParams) ([]db.ListUsersRow, error) {
	r.logger.Info("Getting users with filters",
		zap.Int32("limit", args.Limit),
		zap.Int32("offset", args.Offset),
		zap.String("role", args.Column1),
		zap.String("searchTerm", args.Column2))

	users, err := r.queries.ListUsers(ctx, args)
	if err != nil {
		r.logger.Error("Failed to get users with filters",
			zap.Error(err),
			zap.Int32("limit", args.Limit),
			zap.Int32("offset", args.Offset),
			zap.String("role", args.Column2),
			zap.String("searchTerm", args.Column2))
		return nil, fmt.Errorf("failed to get users with filters: %w", err)
	}

	result := make([]db.ListUsersRow, len(users))
	for i, user := range users {
		result[i] = db.ListUsersRow{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		}
	}

	return result, nil
}

// GetUserById returns a User entity by input id.
func (r *UserRepository) GetUserById(ctx context.Context, id int32) (db.GetUserByIDRow, error) {
	r.logger.Info("Getting user by ID", zap.Int32("user_id", id))

	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		r.logger.Error("Failed to get user by ID", zap.Int32("user_id", id), zap.Error(err))
		return db.GetUserByIDRow{}, fmt.Errorf("failed to get user by ID: %w", err)
	}

	result := db.GetUserByIDRow{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return result, nil
}

// GetUserWithPasswordById returns a User entity by input id with password field in the response.
func (r *UserRepository) GetUserWithPasswordById(ctx context.Context, id int32) (db.User, error) {
	r.logger.Info("Getting user with password by ID", zap.Int32("user_id", id))

	user, err := r.queries.GetUserWithPasswordById(ctx, id)
	if err != nil {
		r.logger.Error("Failed to get user with password by ID", zap.Int32("user_id", id), zap.Error(err))
		return db.User{}, fmt.Errorf("failed to get user with password by email: %w", err)
	}

	return user, nil
}

// CheckUserExists returns a boolean response whether a User exists in the database based on input id.
func (r *UserRepository) CheckUserExists(ctx context.Context, id int32) (bool, error) {
	r.logger.Info("Checking if user exists", zap.Int32("user_id", id))

	result, err := r.queries.CheckUserExists(ctx, id)
	if err != nil {
		r.logger.Error("Failed to check if user exists", zap.Int32("user_id", id), zap.Error(err))
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return result, nil
}

// UpdateUserPartial partially updates a User.
func (r *UserRepository) UpdateUserPartial(ctx context.Context, arg db.UpdateUserPartialParams) (db.User, error) {
	r.logger.Info("Partially updating user", zap.Int32("user_id", arg.ID))

	updatedUser, err := r.queries.UpdateUserPartial(ctx, arg)
	if err != nil {
		r.logger.Error("Failed to partially update user", zap.Int32("user_id", arg.ID), zap.Error(err))
		return db.User{}, fmt.Errorf("failed to partially update user: %w", err)
	}

	return updatedUser, nil
}

// DeleteUserById deletes user by provided id.
func (r *UserRepository) DeleteUserById(ctx context.Context, id int32) error {
	r.logger.Info("Deleting user by ID", zap.Int32("user_id", id))

	err := r.queries.DeleteUserByID(ctx, id)
	if err != nil {
		r.logger.Error("Failed to delete user by ID", zap.Int32("user_id", id), zap.Error(err))
		return fmt.Errorf("failed to delete user by ID: %w", err)
	}

	return nil
}
