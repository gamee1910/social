package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gamee1910/social/internal/domain/entity"
)

var (
	ErrNotFound = errors.New("resources not found")
)

type Storage struct {
	Users interface {
		Create(ctx context.Context, user *entity.User) error
	}

	Posts interface {
		Create(ctx context.Context, post *entity.Post) error
		GetById(ctx context.Context, postId int64) (*entity.Post, error)
		Delete(ctx context.Context, postId int64) error
		Update(ctx context.Context, postId int64, post *entity.Post) (*entity.Post, error)
	}

	Comments interface {
		GetByPostId(ctx context.Context, postId int64) ([]entity.Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}
