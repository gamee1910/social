package service

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
)

type CommentService interface {
	GetByPostId(ctx context.Context, postID int64) ([]entity.Comment, error)
}
