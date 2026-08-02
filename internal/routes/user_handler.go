package routes

import (
	"net/http"

	"github.com/gamee1910/social/internal/config"
)

var userID string = "userID"

func (h *Handler) getUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getIDFromParameter(userID, r)
	if err != nil {
		BadRequestError(w, r, err)
		return
	}

	user, err := h.service.UserService.GetById(ctx, userID)
	if err != nil {
		HandleServiceError(w, r, err)
		return
	}

	if err := config.ResponseJSON(w, http.StatusOK, user); err != nil {
		InternalServerError(w, r, err)
		return
	}
}
