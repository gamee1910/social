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
	"github.com/gamee1910/social/pkg/logger"
	"github.com/lib/pq"
)

type userRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewUserRepository(
	db *sql.DB,
	logger *logger.Logger,
) repository.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (repository *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	query := "INSERT INTO users(username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email, created_at"

	var result entity.User

	err := repository.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&result.ID,
		&result.Username,
		&result.Email,
		&result.CreatedAt,
	)

	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				switch pqErr.Constraint {
				case "users_username_key":
					return nil, domain.ErrUsernameAlreadyExists

				case "users_email_key":
					return nil, domain.ErrEmailAlreadyExists
				}
			}
		}

		return nil, err
	}

	return &result, nil
}

func (repository *userRepository) GetById(ctx context.Context, userID int64) (*entity.User, error) {

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	query := "SELECT id, username, email FROM users WHERE id = $1"

	var user entity.User

	err := repository.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			repository.logger.Warn("user not found", "userID", userID)
			return nil, domain.ErrNotFound
		}

		repository.logger.Error("failed to get user by ID", "userID", userID, "error", err)

		return nil, fmt.Errorf("[user repository] - [get by id] - err: [%w]", err)
	}

	return &user, nil
}

func (repository *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	query := "SELECT id, username, email, password, created_at FROM users WHERE email = $1"

	var result entity.User

	err := repository.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&result.ID,
		&result.Username,
		&result.Email,
		&result.Password,
		&result.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			repository.logger.Warn("user not found", "email", email)
			return nil, domain.ErrNotFound
		}

		repository.logger.Error("Failed to get user by ID", "email", email, "error", err)
		return nil, fmt.Errorf("[user repository] - [Login] - err: [%w]", err)
	}

	return &result, nil
}
