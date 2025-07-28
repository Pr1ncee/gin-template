package handlers

import (
	"GinBox/internal/handlers/requests"
	"GinBox/internal/handlers/responses"
	db "GinBox/internal/postgresql"
	"GinBox/internal/services"
	"GinBox/internal/utils"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	*BaseHandler
	userService services.IUserService
	logger      *zap.Logger
	timeout     time.Duration
}

func NewUserHandler(userService services.IUserService, logger *zap.Logger, timeout time.Duration) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(logger, timeout),
		userService: userService,
		logger:      logger,
		timeout:     timeout,
	}
}

// RefreshToken godoc
// @Summary Refresh user's token
// @Description Refresh user's token and return a new access token in the response
// @Tags users
// @Accept json
// @Produce json
// @Param user body requests.RefreshTokenRequest true "Refresh token"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req requests.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	ctx, cancel := h.GetContextWithTimeout()
	defer cancel()

	accessToken, err := h.userService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		h.SendError(c, http.StatusUnauthorized, "Error refreshing token", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

// Login godoc
// @Summary Log a user in
// @Description Login a user and return a pair of access and refresh token in the response
// @Tags users
// @Accept json
// @Produce json
// @Param user body requests.LoginRequest true "User information"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	ctx, cancel := h.GetContextWithTimeout()
	defer cancel()

	accessToken, refreshToken, err := h.userService.Login(ctx, req)
	if err != nil {
		h.SendError(c, http.StatusUnauthorized, "Failed to login", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user and return user object in the response
// @Tags users
// @Accept json
// @Produce json
// @Param user body requests.CreateUserRequest true "User information"
// @Success 201 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req requests.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	ctx, cancel := h.GetContextWithTimeout()
	defer cancel()

	createdUser, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	h.SendSuccess(c, http.StatusCreated, createdUser, "User created successfully")
}

// ListUsers godoc
// @Summary List all users
// @Description List all users with pagination params
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param role query string false "User Role" default()
// @Param searchTerm query string false "first_name" or "last_name" or "email" default()
// @Success 200 {object} responses.PaginationResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	role := utils.ParseUserRole(c.DefaultQuery("role", ""), h.logger)
	query := c.DefaultQuery("searchTerm", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		h.SendError(c, http.StatusBadRequest, "Invalid page parameter", err)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		h.SendError(c, http.StatusBadRequest, "Invalid limit parameter (1-100)", err)
		return
	}

	ctx, cancel := h.GetContextWithTimeout()
	defer cancel()

	offset := (page - 1) * limit
	users, err := h.userService.ListUsers(ctx, db.ListUsersParams{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Column1: string(role),
		Column2: query,
	})
	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed to retrieve users", err)
		return
	}

	response := responses.PaginationResponse{
		Data:  users,
		Page:  page,
		Limit: limit,
	}

	c.JSON(http.StatusOK, response)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a user by his ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	ctx, cancel := h.GetContextWithTimeout()
	defer cancel()

	user, err := h.userService.GetUser(ctx, int32(id))
	if err != nil {
		h.SendError(c, http.StatusNotFound, "User not found", err)
		return
	}

	h.SendSuccess(c, http.StatusOK, user, "")
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body requests.UpdateUserRequest true "Updated user information"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	var req requests.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	ctx, cancel := h.GetContextWithTimeout()
	defer cancel()

	updatedUser, err := h.userService.UpdateUser(ctx, int32(id), req)
	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed to update the user", err)
		return
	}
	h.SendSuccess(c, http.StatusOK, updatedUser, "User updated successfully")
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		h.SendError(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	ctx, cancel := h.GetContextWithTimeout()
	defer cancel()

	err = h.userService.DeleteUser(ctx, int32(id))
	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed to delete the user", err)
	}

	h.SendSuccess(c, http.StatusOK, nil, "User deleted successfully")
}

// ExportUsers godoc
// @Summary Export users
// @Description Export users with filter parameters
// @Tags users
// @Produce text/csv
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param role query string false "User Role" default()
// @Param searchTerm query string false "first_name" or "last_name" or "email" default()
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Security BearerAuth
// @Router /users/export [get]
func (h *UserHandler) ExportUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	role := utils.ParseUserRole(c.DefaultQuery("role", ""), h.logger)
	query := c.DefaultQuery("searchTerm", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		h.SendError(c, http.StatusBadRequest, "Invalid page parameter", err)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		h.SendError(c, http.StatusBadRequest, "Invalid limit parameter (1-100)", err)
		return
	}

	ctx, cancel := h.GetContextWithTimeout()
	defer cancel()

	offset := (page - 1) * limit
	data, filename, err := h.userService.ExportUsers(ctx, db.ListUsersParams{
		Limit:   int32(limit),
		Offset:  int32(offset),
		Column1: string(role),
		Column2: query,
	})
	if err != nil {
		h.SendError(c, http.StatusInternalServerError, "Failed to export users", err)
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Length", fmt.Sprintf("%d", data.Len()))

	c.Data(http.StatusOK, "text/csv", data.Bytes())
}
