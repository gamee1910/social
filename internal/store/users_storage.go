package store

import (
	"context"
	"database/sql"

	"github.com/gamee1910/social/internal/domain/entity"
)

type UsersStore struct {
	db *sql.DB
}

func (users *UsersStore) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := users.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil

}
