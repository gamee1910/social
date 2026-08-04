package routes

import (
	"errors"
	"net/http"

	"github.com/gamee1910/social/internal/dto"
	"github.com/gamee1910/social/internal/utils"
)

const postIDKey string = "postId"

func (h *Handler) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePostRequest

	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ResponseValidationError(w, r, utils.FormatValidationErrors(err))
		return
	}

	post, err := h.service.PostsService.Create(r.Context(), req)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusCreated, post); err != nil {
		utils.InternalServerError(w, r, err)
	}
}

func (h *Handler) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postId, err := getIDFromParameter(postIDKey, r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	post, err := h.service.PostsService.GetByIdWithComments(ctx, postId)
	if err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, post); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

func (h *Handler) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postId, err := getIDFromParameter(postIDKey, r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.service.PostsService.Delete(ctx, postId); err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := getIDFromParameter(postIDKey, r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	var req dto.UpdatePostRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if req.Title == nil && req.Content == nil && len(req.Tags) == 0 {
		utils.BadRequestError(w, r, errors.New("at least one field must be provided for update"))
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ResponseValidationError(w, r, utils.FormatValidationErrors(err))
		return
	}

	post, err := h.service.PostsService.Update(ctx, id, req)
	if err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, post); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}
