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
	Posts PostStore
	Users UserStore
}

type PostStore interface {
	Create(context.Context, *domain.Post) error
	GetById(context.Context, int64) (*domain.Post, error)
	GetAll(context.Context) ([]*domain.Post, error)
}
type UserStore interface {
	Create(context.Context, *domain.User) error
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
	}
}
