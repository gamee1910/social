package routes

import (
	"net/http"

	"github.com/gamee1910/social/internal/utils"
)

const userIDKey string = "userID"

func (h *Handler) getUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := getIDFromParameter(userIDKey, r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	user, err := h.service.UserService.GetById(ctx, id)
	if err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, user); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}
