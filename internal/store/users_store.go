package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gamee1910/social/internal/domain/entity"
)

type UsersStore struct {
	db *sql.DB
}

func NewUsersStore(db *sql.DB) *UsersStore {
	return &UsersStore{db: db}
}

func (us *UsersStore) Create(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := us.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil

}

func (us *UsersStore) GetById(ctx context.Context, userID int64) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `SELECT id, username, email FROM users WHERE id = $1`

	var user entity.User

	err := us.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("user not found: [%d]", userID)
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("[user store] - [GetById] - err: [%w]", err)
	}

	return &user, nil
}
