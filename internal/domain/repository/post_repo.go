package repository

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
)

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) (*entity.Post, error)
	GetById(ctx context.Context, postId int64) (*entity.Post, error)
	Delete(ctx context.Context, postId int64) error
	Update(ctx context.Context, postId int64, post *entity.Post) (*entity.Post, error)
}
