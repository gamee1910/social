package service

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/interfaces/http/transport/request"
)

type PostService interface {
	Create(ctx context.Context, req request.CreatePostRequest) (*entity.Post, error)

	GetById(ctx context.Context, postID int64) (*entity.Post, error)

	GetByIdWithComments(ctx context.Context, postID int64) (*entity.Post, error)

	Delete(ctx context.Context, postID int64) error

	Update(ctx context.Context, postID int64, req request.UpdatePostRequest) (*entity.Post, error)
}
