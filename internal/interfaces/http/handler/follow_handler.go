package handler

import (
	"net/http"

	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/utils"
)

type FollowerHandler struct {
	followerService service.FollowerService
}

func NewFollowerHandler(
	followerService service.FollowerService,
) *FollowerHandler {
	return &FollowerHandler{
		followerService: followerService,
	}
}

// FollowUser
//
// @Summary Follow a user
// @Description Follow another user by user ID
// @Tags Followers
// @Produce json
// @Param userID path int64 true "User ID to follow"
// @Success 204 "Successfully followed user"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /users/{userID}/follow [post]
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

// UnfollowUser
//
// @Summary Unfollow a user
// @Description Stop following another user by user ID
// @Tags Followers
// @Produce json
// @Param userID path int64 true "User ID to unfollow"
// @Success 204 "Successfully unfollowed user"
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /users/{userID}/follow [delete]
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
