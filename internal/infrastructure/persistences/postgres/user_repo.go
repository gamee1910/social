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

func (repository *userRepository) Create(
	ctx context.Context,
	user *entity.User,
) (*entity.User, error) {

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	repository.logger.Info(
		"Creating user",
		"username", user.Username,
		"email", user.Email,
	)

	query := `
	INSERT INTO users(username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, username, email, created_at
	`

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
		repository.logger.Error(
			"Failed to create user",
			"username", user.Username,
			"email", user.Email,
			"error", err,
		)

		return nil, err
	}

	repository.logger.Info(
		"Successfully created user",
		"userID", result.ID,
		"username", result.Username,
		"email", result.Email,
	)

	return &result, nil
}

func (repository *userRepository) GetById(
	ctx context.Context,
	userID int64,
) (*entity.User, error) {

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	repository.logger.Info(
		"Getting user by ID",
		"userID", userID,
	)

	query := `SELECT id, username, email FROM users WHERE id = $1`

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
			repository.logger.Warn(
				"User not found",
				"userID", userID,
			)

			return nil, domain.ErrNotFound
		}

		repository.logger.Error(
			"Failed to get user by ID",
			"userID", userID,
			"error", err,
		)

		return nil, fmt.Errorf(
			"[user store] - [GetById] - err: [%w]",
			err,
		)
	}

	repository.logger.Info(
		"Successfully retrieved user",
		"userID", user.ID,
		"username", user.Username,
	)

	return &user, nil
}
