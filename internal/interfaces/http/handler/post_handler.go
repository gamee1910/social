package handler

import (
	"errors"
	"net/http"

	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/request"

	"github.com/gamee1910/social/internal/utils"
	"github.com/gamee1910/social/pkg/logger"
)

type PostHandler struct {
	postService service.PostService
	logger      *logger.Logger
}

func NewPostHandler(postService service.PostService, logger *logger.Logger) *PostHandler {
	return &PostHandler{
		postService: postService,
		logger:      logger,
	}
}

const postIDKey string = "postID"

func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var req request.CreatePostRequest

	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.BadRequestError(w, r, err)
		return
	}

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

	//TODO: userID := getUserIDFromContext(r.Context())
	userID := int64(1)

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
