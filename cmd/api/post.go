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
		_ = writeJSON(response, http.StatusBadRequest, err.Error())
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
		_ = writeJSON(response, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(response, http.StatusCreated, post); err != nil {
		_ = responseJSONError(response, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostHandler(response http.ResponseWriter, request *http.Request) {
	context := request.Context()

	idParam := chi.URLParam(request, "postId")
	postId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		_ = writeJSON(response, http.StatusInternalServerError, err.Error())
		return
	}

	post, err := app.store.Posts.GetById(context, postId)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			_ = responseJSONError(response, http.StatusNotFound, err.Error())
		default:
			_ = responseJSONError(response, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(response, http.StatusCreated, post); err != nil {
		_ = responseJSONError(response, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getAllPostHandler(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	posts, err := app.store.Posts.GetAll(ctx)
	if err != nil {
		_ = responseJSONError(response, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(response, http.StatusOK, posts); err != nil {
		_ = responseJSONError(response, http.StatusInternalServerError, err.Error())
		return
	}
}
