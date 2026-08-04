package handler

import (
	"errors"
	"net/http"

	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/dto"
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

const postIDKey string = "postId"

func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePostRequest

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

// func (h *Handler) getUserById(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	id, err := getIDFromParameter(userIDKey, r)
// 	if err != nil {
// 		utils.BadRequestError(w, r, err)
// 		return
// 	}

// 	user, err := h.service.UserService.GetById(ctx, id)
// 	if err != nil {
// 		utils.HandleServiceError(w, r, err)
// 		return
// 	}

// 	if err := utils.ResponseJSON(w, http.StatusOK, user); err != nil {
// 		utils.InternalServerError(w, r, err)
// 		return
// 	}
// }
