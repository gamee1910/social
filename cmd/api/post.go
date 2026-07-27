package main

import (
	"errors"
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

func (app *application) createPostHandler(response http.ResponseWriter, request *http.Request) {
	var postRequest CreatePostRequest

	if err := httpx.ReadJSON(response, request, &postRequest); err != nil {
		httpx.BadRequestError(response, request, err)
		return
	}

	if err := httpx.Validate.Struct(postRequest); err != nil {
		formatError := httpx.FormatValidationErrors(err)
		httpx.ResponseValidationError(response, request, formatError)
		return
	}

	post := &domain.Post{
		Title:   postRequest.Title,
		Content: postRequest.Content,
		Tags:    postRequest.Tags,
		//TODO: Change after auth
		UserId: 1,
	}

	ctx := request.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		httpx.InternalServerError(response, request, err)
		return
	}

	if err := httpx.WriteJSON(response, http.StatusCreated, post); err != nil {
		httpx.InternalServerError(response, request, err)
		return
	}
}

func (app *application) getPostHandler(response http.ResponseWriter, request *http.Request) {
	context := request.Context()

	idParam := chi.URLParam(request, "postId")
	postId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		httpx.InternalServerError(response, request, err)
		return
	}

	post, err := app.store.Posts.GetById(context, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			httpx.NotFoundError(response, request, err)
		default:
			httpx.InternalServerError(response, request, err)
		}
		return
	}

	comments, err := app.store.Comments.GetByPostId(context, postId)
	if err != nil {
		httpx.InternalServerError(response, request, err)
		return
	}

	post.Comment = comments

	if err := httpx.WriteJSON(response, http.StatusCreated, post); err != nil {
		httpx.InternalServerError(response, request, err)
		return
	}
}

func (app *application) getAllPostHandler(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	posts, err := app.store.Posts.GetAll(ctx)
	if err != nil {
		httpx.InternalServerError(response, request, err)
		return
	}

	if err := httpx.WriteJSON(response, http.StatusOK, posts); err != nil {
		httpx.InternalServerError(response, request, err)
		return
	}
}
