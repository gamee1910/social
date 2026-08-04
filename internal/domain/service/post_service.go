package service

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/dto"
)

type PostService interface {
	Create(ctx context.Context, req dto.CreatePostRequest) (*entity.Post, error)

	GetById(ctx context.Context, postID int64) (*entity.Post, error)

	GetByIdWithComments(ctx context.Context, postID int64) (*entity.Post, error)

	Delete(ctx context.Context, postID int64) error

	Update(ctx context.Context, postID int64, req dto.UpdatePostRequest) (*entity.Post, error)
}
