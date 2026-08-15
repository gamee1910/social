package handler

import (
	"net/http"

	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/request"
	"github.com/gamee1910/social/internal/utils"
	"github.com/gamee1910/social/pkg/logger"
)

type UserHandler struct {
	userService service.UserService
	logger      *logger.Logger
}

const userIDKey string = "userID"

func NewUserHandler(userService service.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetUserById
//
// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags Users
// @Produce json
// @Param userID path int64 true "User ID"
// @Success 200 {object} response.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userID} [get]
func (handler *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := utils.GetIDFromParameter(userIDKey, r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	user, err := handler.userService.GetById(ctx, id)
	if err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, user); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

// Register
//
// @Summary Register
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body request.UserCreationRequest true "User registration request"
// @Success 201 {object} response.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/register [post]
func (handler *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req request.UserCreationRequest

	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ResponseValidationError(w, r, utils.FormatValidationErrors(err))
		return
	}

	user, err := handler.userService.Register(r.Context(), req)
	if err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, user); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

// Login
//
// @Summary Login
// @Description Login user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body request.UserLoginRequest true "User login request"
// @Success 200 {object} response.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/login [post]
func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.UserLoginRequest

	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ResponseValidationError(w, r, utils.FormatValidationErrors(err))
		return
	}

	user, err := handler.userService.Login(r.Context(), req)
	if err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, user); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}
