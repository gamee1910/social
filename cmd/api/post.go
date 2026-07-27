package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/httpx"
	"github.com/gamee1910/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostRequest struct {
	Title   string   `json:"title" validate:"required,max=200"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostRequest struct {
	Title   *string  `json:"title" validate:"omitempty,min=1,max=500"`
	Content *string  `json:"content" validate:"omitempty,min=1,max=50000"`
	Tags    []string `json:"tags" validate:"max=20,dive,min=1,max=50"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postRequest CreatePostRequest

	if err := httpx.ReadJSON(w, r, &postRequest); err != nil {
		httpx.BadRequestError(w, r, err)
		return
	}

	if err := httpx.Validate.Struct(postRequest); err != nil {
		formatError := httpx.FormatValidationErrors(err)
		httpx.ResponseValidationError(w, r, formatError)
		return
	}

	post := &domain.Post{
		Title:   postRequest.Title,
		Content: postRequest.Content,
		Tags:    postRequest.Tags,
		//TODO: Change after auth
		UserId: 1,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		httpx.InternalServerError(w, r, err)
		return
	}

	if err := httpx.WriteJSON(w, http.StatusCreated, post); err != nil {
		httpx.InternalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postId, err := getIDFromParameter(r)
	if err != nil {
		httpx.BadRequestError(w, r, err)
		return
	}

	post, err := app.store.Posts.GetById(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			httpx.NotFoundError(w, r, err)
			return
		default:
			httpx.InternalServerError(w, r, err)
			return
		}
	}

	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)
	if err != nil {
		httpx.InternalServerError(w, r, err)
		return
	}

	post.Comment = comments

	if err := httpx.WriteJSON(w, http.StatusCreated, post); err != nil {
		httpx.InternalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postId, err := getIDFromParameter(r)
	if err != nil {
		httpx.BadRequestError(w, r, err)
		return
	}

	if err := app.store.Posts.Delete(ctx, postId); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httpx.NotFoundError(w, r, err)
			return
		}

		httpx.InternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID, err := getIDFromParameter(r)
	if err != nil {
		httpx.BadRequestError(w, r, err)
		return
	}

	var updateRequest UpdatePostRequest
	if err := httpx.ReadJSON(w, r, &updateRequest); err != nil {
		httpx.BadRequestError(w, r, err)
		return
	}

	if updateRequest.Content == nil && updateRequest.Title == nil && len(updateRequest.Tags) == 0 {
		httpx.BadRequestError(w, r, errors.New("at least one field must be provided for update"))
		return
	}

	if err := httpx.Validate.Struct(updateRequest); err != nil {
		formatError := httpx.FormatValidationErrors(err)
		httpx.ResponseValidationError(w, r, formatError)
		return
	}

	exitingPost, err := app.store.Posts.GetById(ctx, postID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			httpx.NotFoundError(w, r, errors.New("post not found"))
			return
		default:
			httpx.InternalServerError(w, r, err)
			return
		}
	}

	// TODO: verify user owns this post
	//userID := r.Context().Value("userID").(int64)
	//if existingPost.UserId != userID {
	//	httpx.ForbiddenError(w, r, errors.New("cannot update post owned by another user"))
	//	return
	//}
	updatePost := &domain.Post{
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

	updatedPost, err := app.store.Posts.Update(ctx, postID, updatePost)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			httpx.NotFoundError(w, r, errors.New("post not found"))
			return
		default:
			httpx.InternalServerError(w, r, err)
			return
		}
	}

	comments, err := app.store.Comments.GetByPostId(ctx, postID)
	if err != nil {
		comments = []domain.Comment{}
	}
	updatedPost.Comment = comments

	if err := httpx.WriteJSON(w, http.StatusOK, updatedPost); err != nil {
		httpx.InternalServerError(w, r, err)
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
