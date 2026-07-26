package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(response http.ResponseWriter, request *http.Request) {
	var postRequest CreatePostRequest

	if err := readJSON(response, request, &postRequest); err != nil {
		app.badRequestError(response, request, err)
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
		app.internalServerError(response, request, err)
		return
	}

	if err := writeJSON(response, http.StatusCreated, post); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}

func (app *application) getPostHandler(response http.ResponseWriter, request *http.Request) {
	context := request.Context()

	idParam := chi.URLParam(request, "postId")
	postId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(response, request, err)
		return
	}

	post, err := app.store.Posts.GetById(context, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(response, request, err)
		default:
			app.internalServerError(response, request, err)
		}
		return
	}

	if err := writeJSON(response, http.StatusCreated, post); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}

func (app *application) getAllPostHandler(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	posts, err := app.store.Posts.GetAll(ctx)
	if err != nil {
		app.internalServerError(response, request, err)
		return
	}

	if err := writeJSON(response, http.StatusOK, posts); err != nil {
		app.internalServerError(response, request, err)
		return
	}
}
