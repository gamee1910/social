package service

import (
	"context"

	"github.com/gamee1910/social/internal/interfaces/http/transport/request"
	"github.com/gamee1910/social/internal/interfaces/http/transport/response"
)

type GetFeedInput struct {
	UserID int64
	Limit  int
	Offset int
	Sort   string
	Search string
	Tags   []string
}

type PostService interface {
	Create(ctx context.Context, req request.CreatePostRequest) (*response.PostResponse, error)
	GetById(ctx context.Context, postID int64) (*response.PostResponse, error)
	GetByIdWithComments(ctx context.Context, postID int64) (*response.PostResponse, error)
	Delete(ctx context.Context, postID int64) error
	Update(ctx context.Context, postID int64, req request.UpdatePostRequest) (*response.PostResponse, error)
	GetFeed(ctx context.Context, input GetFeedInput) ([]*response.PostWithMetaData, error)
}
