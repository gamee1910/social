package store

import (
	"context"
	"database/sql"
)

type UsersStore struct {
	db *sql.DB
}

func (users *UsersStore) Create(ctx context.Context) error {
	return nil
}
