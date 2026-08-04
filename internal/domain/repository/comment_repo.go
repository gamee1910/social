package repository

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
)

type CommentRepository interface {
	GetByPostId(ctx context.Context, postId int64) ([]entity.Comment, error)
}
