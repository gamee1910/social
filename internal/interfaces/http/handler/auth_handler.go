package handler

import (
	"net/http"

	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/request"
	"github.com/gamee1910/social/internal/utils"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
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
// @Router /auth/register [post]
func (handler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req request.UserCreationRequest

	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ResponseValidationError(w, r, utils.FormatValidationErrors(err))
		return
	}

	user, err := handler.authService.Register(r.Context(), req)
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
// @Router /auth/login [post]
func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.UserLoginRequest

	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ResponseValidationError(w, r, utils.FormatValidationErrors(err))
		return
	}

	user, err := handler.authService.Login(r.Context(), req)
	if err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, user); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}
