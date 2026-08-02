package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/domain/entity"
)

var QueryTimeoutDuration = time.Second * 5

var (
	ErrNotFound        = domain.ErrNotFound
	ErrVersionConflict = domain.ErrVersionConflict
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetById(ctx context.Context, userID int64) (*entity.User, error)
}

type PostRepository interface {
	Create(ctx context.Context, post *entity.Post) error
	GetById(ctx context.Context, postId int64) (*entity.Post, error)
	Delete(ctx context.Context, postId int64) error
	Update(ctx context.Context, postId int64, post *entity.Post) (*entity.Post, error)
}

type CommentRepository interface {
	GetByPostId(ctx context.Context, postId int64) ([]entity.Comment, error)
}

type Storage struct {
	Users    UserRepository
	Posts    PostRepository
	Comments CommentRepository
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Posts:    NewPostsStore(db),
		Users:    NewUsersStore(db),
		Comments: NewCommentsStore(db),
	}
}
