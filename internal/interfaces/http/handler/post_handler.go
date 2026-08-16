package handler

import (
	"errors"
	"net/http"

	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/middleware"
	"github.com/gamee1910/social/internal/interfaces/http/transport/request"

	"github.com/gamee1910/social/pkg/utils"
	"github.com/gamee1910/social/pkg/logger"
)

type PostHandler struct {
	postService service.PostService
	logger      *logger.Logger
}

func NewPostHandler(
	postService service.PostService,
	logger *logger.Logger,
) *PostHandler {
	return &PostHandler{
		postService: postService,
		logger:      logger,
	}
}

const postIDKey string = "postID"

// CreatePostHandler
//
// @Summary Create a post
// @Description Create a new post
// @Tags Posts
// @Accept json
// @Produce json
// @Param request body request.CreatePostRequest true "Create post request"
// @Success 201 {object} response.PostResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /posts/ [post]
func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var req request.CreatePostRequest

	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	claims := middleware.GetUserClaims(r)
	req.UserID = claims.UserID

	if err := utils.Validate.Struct(req); err != nil {
		utils.ResponseValidationError(w, r, utils.FormatValidationErrors(err))
		return
	}

	post, err := h.postService.Create(r.Context(), req)
	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusCreated, post); err != nil {
		utils.InternalServerError(w, r, err)
	}
}

// GetPostHandler
//
// @Summary Get a post
// @Description Get a post by ID
// @Tags Posts
// @Produce json
// @Param postID path int64 true "Post ID"
// @Success 200 {object} response.PostResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /posts/{postID}/ [get]
func (h *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postId, err := utils.GetIDFromParameter(postIDKey, r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	post, err := h.postService.GetByIdWithComments(ctx, postId)
	if err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, post); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

// DeletePostHandler
//
// @Summary Delete a post
// @Description Delete a post by ID
// @Tags Posts
// @Param postID path int64 true "Post ID"
// @Success 204
// @Failure 404 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /posts/{postID}/ [delete]
func (h *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postId, err := utils.GetIDFromParameter(postIDKey, r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	if err := h.postService.Delete(ctx, postId); err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePostHandler
//
// @Summary Update a post
// @Description Update a post by ID
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path int64 true "Post ID"
// @Param request body request.UpdatePostRequest true "Update post request"
// @Success 200 {object} response.PostResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /posts/{postID}/ [patch]
func (h *PostHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := utils.GetIDFromParameter(postIDKey, r)
	if err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	var req request.UpdatePostRequest
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

	post, err := h.postService.Update(ctx, id, req)
	if err != nil {
		utils.HandleServiceError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, post); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}

// GetFeed
//
// @Summary Get feed
// @Description Get paginated posts for the current user
// @Tags Posts
// @Produce json
// @Param limit query int false "Number of posts to return"
// @Param offset query int false "Number of posts to skip"
// @Param sort query string false "Sort order" Enums(asc,desc)
// @Param search query string false "Search keyword"
// @Param tags query []string false "Filter by tags"
// @Success 200 {array} response.PostResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posts/feed [get]
func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	req := request.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
	}

	if err := req.Parse(r); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

	userID := middleware.GetUserClaims(r).UserID

	posts, err := h.postService.GetFeed(r.Context(), service.GetFeedInput{
		UserID: userID,
		Limit:  req.Limit,
		Offset: req.Offset,
		Sort:   req.Sort,
		Search: req.Search,
		Tags:   req.Tags,
	})

	if err != nil {
		utils.InternalServerError(w, r, err)
		return
	}

	if err := utils.ResponseJSON(w, http.StatusOK, posts); err != nil {
		utils.InternalServerError(w, r, err)
		return
	}
}
