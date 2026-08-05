package handler

import (
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/pkg/logger"
)

type CommentHandler struct {
	commentService service.CommentService
	logger         *logger.Logger
}

func NewCommentHandler(commentService service.CommentService, logger *logger.Logger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		logger:         logger,
	}
}
