package handler

import (
	"net/http"

	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/utils"
)

type FollowerHandler struct {
	followerService service.FollowerService
}

func NewFollowerHandler(followerService service.FollowerService) *FollowerHandler {
	return &FollowerHandler{followerService: followerService}
}
func (h *FollowerHandler) FollowUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	followingID, err := utils.GetIDFromParameter(userIDKey, r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	// TODO: replace with authenticated user ID from JWT/context
	followerID := int64(1)

	err = h.followerService.FollowUser(
		ctx,
		followingID,
		followerID,
	)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FollowerHandler) UnfollowUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	followingID, err := utils.GetIDFromParameter(userIDKey, r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	// TODO: replace with authenticated user ID from JWT/context
	followerID := int64(1)

	err = h.followerService.UnfollowUser(
		ctx,
		followingID,
		followerID,
	)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
