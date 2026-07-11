package store

import (
	"context"
	"database/sql"

	"github.com/gamee1910/social/internal/domain"
)

type PostStore interface {
	Create(context.Context, *domain.Post) error
}
type UserStore interface {
	Create(context.Context, *domain.User) error
}

type Storage struct {
	Posts PostStore
	Users UserStore
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
	}
}
