package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/dto"
	"github.com/gamee1910/social/internal/store"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postRequest dto.CreatePostRequest

	if err := config.ReadJSON(w, r, &postRequest); err != nil {
		domain.BadRequestError(w, r, err)
		return
	}

	if err := config.Validate.Struct(postRequest); err != nil {
		formatError := config.FormatValidationErrors(err)
		domain.ResponseValidationError(w, r, formatError)
		return
	}

	post := &entity.Post{
		Title:   postRequest.Title,
		Content: postRequest.Content,
		Tags:    postRequest.Tags,
		//TODO: Change after auth
		UserId: 1,
	}

	ctx := r.Context()

	if err := h.store.Posts.Create(ctx, post); err != nil {
		domain.InternalServerError(w, r, err)
		return
	}

	if err := config.WriteJSON(w, http.StatusCreated, post); err != nil {
		domain.InternalServerError(w, r, err)
		return
	}
}

func (h *Handler) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postId, err := getIDFromParameter(r)
	if err != nil {
		domain.BadRequestError(w, r, err)
		return
	}

	post, err := h.store.Posts.GetById(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			domain.NotFoundError(w, r, err)
			return
		default:
			domain.InternalServerError(w, r, err)
			return
		}
	}

	comments, err := h.store.Comments.GetByPostId(r.Context(), post.ID)
	if err != nil {
		domain.InternalServerError(w, r, err)
		return
	}

	post.Comment = comments

	if err := config.WriteJSON(w, http.StatusCreated, post); err != nil {
		domain.InternalServerError(w, r, err)
		return
	}
}

func (h *Handler) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postId, err := getIDFromParameter(r)
	if err != nil {
		domain.BadRequestError(w, r, err)
		return
	}

	if err := h.store.Posts.Delete(ctx, postId); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			domain.NotFoundError(w, r, err)
			return
		}

		domain.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID, err := getIDFromParameter(r)
	if err != nil {
		domain.BadRequestError(w, r, err)
		return
	}

	var updateRequest dto.UpdatePostRequest
	if err := config.ReadJSON(w, r, &updateRequest); err != nil {
		domain.BadRequestError(w, r, err)
		return
	}

	if updateRequest.Content == nil && updateRequest.Title == nil && len(updateRequest.Tags) == 0 {
		domain.BadRequestError(w, r, errors.New("at least one field must be provided for update"))
		return
	}

	if err := config.Validate.Struct(updateRequest); err != nil {
		formatError := config.FormatValidationErrors(err)
		domain.ResponseValidationError(w, r, formatError)
		return
	}

	exitingPost, err := h.store.Posts.GetById(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			domain.NotFoundError(w, r, errors.New("post not found"))
			return
		default:
			domain.InternalServerError(w, r, err)
			return
		}
	}

	// TODO: verify user owns this post
	//userID := r.Context().Value("userID").(int64)
	//if existingPost.UserId != userID {
	//	httpx.ForbiddenError(w, r, errors.New("cannot update post owned by another user"))
	//	return
	//}
	updatePost := &entity.Post{
		ID:      postID,
		Content: exitingPost.Content,
		Title:   exitingPost.Title,
		UserId:  exitingPost.UserId,
		Tags:    exitingPost.Tags,
	}

	if updateRequest.Title != nil {
		updatePost.Title = *updateRequest.Title
	}
	if updateRequest.Content != nil {
		updatePost.Content = *updateRequest.Content
	}
	if len(updateRequest.Tags) > 0 {
		updatePost.Tags = updateRequest.Tags
	}

	updatedPost, err := h.store.Posts.Update(ctx, postID, updatePost)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			domain.NotFoundError(w, r, errors.New("post not found"))
			return
		default:
			domain.InternalServerError(w, r, err)
			return
		}
	}

	comments, err := h.store.Comments.GetByPostId(ctx, postID)
	if err != nil {
		comments = []entity.Comment{}
	}
	updatedPost.Comment = comments

	if err := config.WriteJSON(w, http.StatusOK, updatedPost); err != nil {
		domain.InternalServerError(w, r, err)
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
