package postgres

import (
	"context"
	"database/sql"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/pkg/logger"
)

type commentRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewCommentRepository(db *sql.DB, logger *logger.Logger) repository.CommentRepository {
	return &commentRepository{
		db:     db,
		logger: logger,
	}
}

func (commentRepo *commentRepository) GetByPostId(
	ctx context.Context,
	postId int64,
) ([]entity.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	commentRepo.logger.Info(
		"Getting comments by post ID",
		"postID", postId,
	)

	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.id
		FROM comments c JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	rows, err := commentRepo.db.QueryContext(ctx, query, postId)
	if err != nil {
		commentRepo.logger.Error(
			"Failed to get comments by post ID",
			"postID", postId,
			"error", err,
		)
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			commentRepo.logger.Error(
				"Failed to close comment rows",
				"postID", postId,
				"error", err,
			)
		}
	}()

	comments := make([]entity.Comment, 0)

	for rows.Next() {
		var cmt entity.Comment
		cmt.User = entity.User{}

		err := rows.Scan(
			&cmt.ID,
			&cmt.PostId,
			&cmt.UserId,
			&cmt.Content,
			&cmt.CreatedAt,
			&cmt.User.Username,
			&cmt.User.ID,
		)

		if err != nil {
			commentRepo.logger.Error(
				"Failed to scan comment",
				"postID", postId,
				"error", err,
			)
			return nil, err
		}

		comments = append(comments, cmt)
	}

	if err := rows.Err(); err != nil {
		commentRepo.logger.Error(
			"Error while iterating comments",
			"postID", postId,
			"error", err,
		)
		return nil, err
	}

	commentRepo.logger.Info(
		"Successfully retrieved comments",
		"postID", postId,
		"count", len(comments),
	)

	return comments, nil
}
