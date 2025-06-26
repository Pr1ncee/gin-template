package repositories

import (
	db "GinBox/internal/postgresql"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user db.User) (db.User, error)
	ListUsers(ctx context.Context, args db.ListUsersParams) ([]db.ListUsersRow, error)
	GetUserById(ctx context.Context, id int32) (db.GetUserByIDRow, error)
	GetUserWithPasswordById(ctx context.Context, id int32) (db.User, error)
	CheckUserExists(ctx context.Context, id int32) (bool, error)
	UpdateUserPartial(ctx context.Context, arg db.UpdateUserPartialParams) (db.User, error)
	DeleteUserById(ctx context.Context, id int32) error
}

type UserRepository struct {
	queries *db.Queries
	logger  *log.Logger
}

func NewUserRepository(queries *db.Queries, logger *log.Logger) *UserRepository {
	return &UserRepository{
		queries: queries,
		logger:  logger,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user db.User) (db.User, error) {
	r.logger.WithField("email", user.Email).Info("Creating new user")

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
		r.logger.WithError(err).WithField("email", user.Email).Error("Failed to create user")
		return db.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.WithFields(log.Fields{
		"user_id": createdUser.ID,
		"email":   createdUser.Email,
	}).Info("User created successfully")

	return createdUser, nil
}

func (r *UserRepository) ListUsers(ctx context.Context, args db.ListUsersParams) ([]db.ListUsersRow, error) {
	r.logger.WithFields(log.Fields{
		"limit":      args.Limit,
		"offset":     args.Offset,
		"role":       args.Column1,
		"searchTerm": args.Column2,
	}).Info("Getting users with filters")

	users, err := r.queries.ListUsers(ctx, args)
	if err != nil {
		r.logger.WithError(err).WithFields(log.Fields{
			"limit":      args.Limit,
			"offset":     args.Offset,
			"role":       args.Column1,
			"searchTerm": args.Column2,
		}).Error("Failed to get users with filters")
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

	r.logger.WithFields(log.Fields{
		"count":      len(result),
		"limit":      args.Limit,
		"offset":     args.Offset,
		"role":       args.Column1,
		"searchTerm": args.Column2,
	}).Info("Users retrieved with filters")

	return result, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int32) (db.GetUserByIDRow, error) {
	r.logger.WithField("user_id", id).Info("Getting user by ID")

	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		r.logger.WithError(err).WithField("user_id", id).Error("Failed to get user by ID")
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

	r.logger.WithFields(log.Fields{
		"user_id": id,
		"email":   result.Email,
	}).Info("User retrieved by ID")

	return result, nil
}

func (r *UserRepository) GetUserWithPasswordById(ctx context.Context, id int32) (db.User, error) {
	r.logger.WithField("ID", id).Info("Getting user with password by ID")

	user, err := r.queries.GetUserWithPasswordById(ctx, id)
	if err != nil {
		r.logger.WithError(err).WithField("ID", id).Error("Failed to get user with password by ID")
		return db.User{}, fmt.Errorf("failed to get user with password by email: %w", err)
	}

	r.logger.WithField("user_id", user.ID).Info("User with password retrieved by ID")

	return user, nil
}

func (r *UserRepository) CheckUserExists(ctx context.Context, id int32) (bool, error) {
	r.logger.WithField("user_id", id).Info("Checking if user exists")

	result, err := r.queries.CheckUserExists(ctx, id)
	if err != nil {
		r.logger.WithError(err).WithField("user_id", id).Error("Failed to check if user exists")
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	r.logger.WithFields(log.Fields{
		"user_id": id,
		"exists":  result,
	}).Info("User existence check completed")

	return result, nil
}

func (r *UserRepository) UpdateUserPartial(ctx context.Context, arg db.UpdateUserPartialParams) (db.User, error) {
	r.logger.WithField("user_id", arg.ID).Info("Partially updating user")

	updatedUser, err := r.queries.UpdateUserPartial(ctx, arg)
	if err != nil {
		r.logger.WithError(err).WithField("user_id", arg.ID).Error("Failed to partially update user")
		return db.User{}, fmt.Errorf("failed to partially update user: %w", err)
	}

	r.logger.WithFields(log.Fields{
		"user_id": updatedUser.ID,
		"email":   updatedUser.Email,
	}).Info("User partially updated successfully")

	return updatedUser, nil
}

func (r *UserRepository) DeleteUserById(ctx context.Context, id int32) error {
	r.logger.WithField("user_id", id).Info("Deleting user by ID")

	err := r.queries.DeleteUserByID(ctx, id)
	if err != nil {
		r.logger.WithError(err).WithField("user_id", id).Error("Failed to delete user by ID")
		return fmt.Errorf("failed to delete user by ID: %w", err)
	}

	r.logger.WithField("user_id", id).Info("User deleted successfully by ID")
	return nil
}
