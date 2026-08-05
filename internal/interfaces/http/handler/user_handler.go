package handler

import (
	"net/http"

	"github.com/gamee1910/social/internal/domain/service"
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
