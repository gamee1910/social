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
		BadRequestError(w, r, err)
		return
	}

	user, err := h.service.UserService.GetById(ctx, id)
	if err != nil {
		HandleServiceError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, user); err != nil {
		InternalServerError(w, r, err)
		return
	}
}
