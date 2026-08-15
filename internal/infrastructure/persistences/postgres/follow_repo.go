package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/pkg/logger"
)

type followRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewFollowRepository(db *sql.DB, logger *logger.Logger) repository.FollowRepository {
	return &followRepository{db: db, logger: logger}
}

func (repository *followRepository) FollowUser(
	ctx context.Context,
	userID, followerID int64,
) error {
	ctx, cancel := context.WithTimeout(
		ctx,
		config.QueryTimeoutDuration,
	)
	defer cancel()

	query := `
        INSERT INTO followers(user_id, follower_id, created_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id, follower_id) DO NOTHING
    `

	_, err := repository.db.ExecContext(
		ctx,
		query,
		userID,
		followerID,
	)

	if err != nil {
		repository.logger.Error("follow user failed", "userID", userID, "followerID", followerID, "error", err)
		return fmt.Errorf("follow user: %w", err)
	}

	return nil
}

func (repository *followRepository) UnfollowUser(
	ctx context.Context,
	userID, followerID int64,
) error {
	ctx, cancel := context.WithTimeout(
		ctx,
		config.QueryTimeoutDuration,
	)
	defer cancel()

	query := `
        DELETE FROM followers
        WHERE user_id = $1
          AND follower_id = $2
    `

	result, err := repository.db.ExecContext(
		ctx,
		query,
		userID,
		followerID,
	)
	if err != nil {
		repository.logger.Error("unfollow user failed", "userID", userID, "followerID", followerID, "error", err)

		return fmt.Errorf("unfollow user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("follow relationship does not exist")
	}

	return nil
}
