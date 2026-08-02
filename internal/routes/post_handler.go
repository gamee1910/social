package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/dto"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePostRequest

	if err := config.ReadJSON(w, r, &req); err != nil {
		BadRequestError(w, r, err)
		return
	}

	if err := config.Validate.Struct(req); err != nil {
		ResponseValidationError(w, r, config.FormatValidationErrors(err))
		return
	}

	post, err := h.service.PostsService.Create(r.Context(), req)
	if err != nil {
		InternalServerError(w, r, err)
		return
	}

	if err := config.ResponseJSON(w, http.StatusCreated, post); err != nil {
		InternalServerError(w, r, err)
	}
}

func (h *Handler) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postId, err := getIDFromParameter(r)
	if err != nil {
		BadRequestError(w, r, err)
		return
	}

	post, err := h.service.PostsService.GetByIdWithComments(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			NotFoundError(w, r, err)
		default:
			InternalServerError(w, r, err)
		}
		return
	}

	if err := config.ResponseJSON(w, http.StatusOK, post); err != nil {
		InternalServerError(w, r, err)
		return
	}
}

func (h *Handler) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postId, err := getIDFromParameter(r)
	if err != nil {
		BadRequestError(w, r, err)
		return
	}

	if err := h.service.PostsService.Delete(ctx, postId); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			NotFoundError(w, r, err)
			return
		}

		InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID, err := getIDFromParameter(r)
	if err != nil {
		BadRequestError(w, r, err)
		return
	}

	var req dto.UpdatePostRequest
	if err := config.ReadJSON(w, r, &req); err != nil {
		BadRequestError(w, r, err)
		return
	}

	if req.Title == nil && req.Content == nil && len(req.Tags) == 0 {
		BadRequestError(w, r, errors.New("at least one field must be provided for update"))
		return
	}

	if err := config.Validate.Struct(req); err != nil {
		ResponseValidationError(w, r, config.FormatValidationErrors(err))
		return
	}

	post, err := h.service.PostsService.Update(ctx, postID, req)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			NotFoundError(w, r, err)

		case errors.Is(err, domain.ErrVersionConflict):
			ConflictError(w, r, err)

		default:
			InternalServerError(w, r, err)
		}
		return
	}

	if err := config.ResponseJSON(w, http.StatusOK, post); err != nil {
		InternalServerError(w, r, err)
		return
	}
}

func getIDFromParameter(r *http.Request) (int64, error) {
	id := chi.URLParam(r, "postId")

	postID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid post id: %w", err)
	}

	return postID, nil
}
