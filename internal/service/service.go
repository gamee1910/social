package service

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/dto"
	"github.com/gamee1910/social/internal/store"
)

type PostServiceInterface interface {
	Create(ctx context.Context, req dto.CreatePostRequest) (*entity.Post, error)
	GetById(ctx context.Context, postID int64) (*entity.Post, error)
	GetByIdWithComments(ctx context.Context, postID int64) (*entity.Post, error)
	Delete(ctx context.Context, postID int64) error
	Update(ctx context.Context, postID int64, req dto.UpdatePostRequest) (*entity.Post, error)
}

type CommentServiceInterface interface {
	GetByPostId(ctx context.Context, postID int64) ([]entity.Comment, error)
}

type UserServiceInterface interface {
	GetById(ctx context.Context, userID int64) (*entity.User, error)
}

type Service struct {
	PostsService   PostServiceInterface
	CommentService CommentServiceInterface
	UserService    UserServiceInterface
}

func NewService(storage *store.Storage) *Service {
	return &Service{
		PostsService:   NewPostService(storage.Posts, storage.Comments),
		CommentService: NewCommentService(storage.Comments),
		UserService:    NewUserService(storage.Users),
	}
}
