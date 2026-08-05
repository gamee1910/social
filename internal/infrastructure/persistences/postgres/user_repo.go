package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/domain/repository"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (userRepo *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	query := `INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email, created_at`

	var result entity.User

	err := userRepo.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(
		&result.ID,
		&result.Username,
		&result.Email,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil

}

func (userRepo *userRepository) GetById(ctx context.Context, userID int64) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	query := `SELECT id, username, email FROM users WHERE id = $1`

	var user entity.User

	err := userRepo.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("[user store] - [GetById] - err: [%w]", err)
	}

	return &user, nil
}
