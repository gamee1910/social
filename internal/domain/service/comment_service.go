package service

import (
	"context"

	"github.com/gamee1910/social/internal/interfaces/http/handler/response"
)

type CommentService interface {
	GetByPostId(ctx context.Context, postID int64) ([]response.CommentResponse, error)
}
