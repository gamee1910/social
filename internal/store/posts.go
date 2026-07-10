package store

import (
	"context"
	"database/sql"
)

type PostsStore struct {
	db *sql.DB
}

func (store *PostsStore) Create(ctx context.Context) error {
	return nil
}
