package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gamee1910/social/internal/domain"
)

var (
	ErrNotFound = errors.New("resources not found")
)

type Storage struct {
	Users interface {
		Create(ctx context.Context, user *domain.User) error
	}

	Posts interface {
		Create(ctx context.Context, post *domain.Post) error
		GetById(ctx context.Context, postId int64) (*domain.Post, error)
		GetAll(ctx context.Context) ([]*domain.Post, error)
	}

	Comments interface {
		GetByPostId(ctx context.Context, postId int64) ([]domain.Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}
