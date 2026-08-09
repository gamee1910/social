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

type postRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewPostRepository(db *sql.DB, logger *logger.Logger) repository.PostRepository {
	return &postRepository{
		db:     db,
		logger: logger,
	}
}

func (postRepo *postRepository) Create(
	ctx context.Context,
	post *entity.Post,
) (*entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	postRepo.logger.Info("Creating post", map[string]any{
		"user_id": post.UserId,
		"title":   post.Title,
	})

	query := `
	INSERT INTO posts(content, title, user_id, tags)
	VALUES ($1, $2, $3, $4)
	RETURNING id, content, title, user_id, tags, created_at, updated_at, version
	`

	var result entity.Post

	err := postRepo.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserId,
		pq.Array(post.Tags),
	).Scan(
		&result.ID,
		&result.Content,
		&result.Title,
		&result.UserId,
		pq.Array(&result.Tags),
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.Version,
	)

	if err != nil {
		postRepo.logger.Error("Failed to create post", map[string]any{
			"user_id": post.UserId,
			"error":   err,
		})

		return nil, fmt.Errorf("create post error: %w", err)
	}

	postRepo.logger.Info("Post created successfully", map[string]any{
		"post_id": result.ID,
		"user_id": result.UserId,
	})

	return &result, nil
}

func (postRepo *postRepository) GetById(
	ctx context.Context,
	postId int64,
) (*entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	postRepo.logger.Info("Getting post by ID", map[string]any{
		"post_id": postId,
	})

	query := `
	SELECT id, user_id, title, content, tags, version, created_at, updated_at
	FROM posts
	WHERE id = $1
	`

	var post entity.Post

	err := postRepo.db.QueryRowContext(ctx, query, postId).
		Scan(
			&post.ID,
			&post.UserId,
			&post.Title,
			&post.Content,
			pq.Array(&post.Tags),
			&post.Version,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			postRepo.logger.Warn("Post not found", map[string]any{
				"post_id": postId,
			})

			return nil, domain.ErrNotFound
		}

		postRepo.logger.Error("Failed to get post", map[string]any{
			"post_id": postId,
			"error":   err,
		})

		return nil, fmt.Errorf("get post: %w", err)
	}

	postRepo.logger.Info("Post retrieved successfully", map[string]any{
		"post_id": post.ID,
		"user_id": post.UserId,
	})

	return &post, nil
}

func (postRepo *postRepository) Delete(
	ctx context.Context,
	postID int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	postRepo.logger.Info("Deleting post", map[string]any{
		"post_id": postID,
	})

	res, err := postRepo.db.ExecContext(
		ctx,
		`DELETE FROM posts WHERE id = $1`,
		postID,
	)
	if err != nil {
		postRepo.logger.Error("Failed to delete post", map[string]any{
			"post_id": postID,
			"error":   err,
		})

		return fmt.Errorf("delete post: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		postRepo.logger.Error("Failed to get affected rows", map[string]any{
			"post_id": postID,
			"error":   err,
		})

		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		postRepo.logger.Warn("Post not found for deletion", map[string]any{
			"post_id": postID,
		})

		return domain.ErrNotFound
	}

	postRepo.logger.Info("Post deleted successfully", map[string]any{
		"post_id": postID,
	})

	return nil
}

func (postRepo *postRepository) Update(
	ctx context.Context,
	postId int64,
	post *entity.Post,
) (*entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	postRepo.logger.Info("Updating post", map[string]any{
		"post_id": postId,
		"version": post.Version,
		"title":   post.Title,
	})

	tx, err := postRepo.db.BeginTx(ctx, nil)
	if err != nil {
		postRepo.logger.Error("Failed to begin update transaction", map[string]any{
			"post_id": postId,
			"error":   err,
		})

		return nil, fmt.Errorf("begin transaction error: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			postRepo.logger.Error("Failed to rollback transaction", map[string]any{
				"post_id": postId,
				"error":   err,
			})
		}
	}()

	query := `
		UPDATE posts SET title = $1, content = $2, tags = $3, version = version + 1, updated_at = NOW()
		WHERE id = $4 AND version = $5
		RETURNING id, user_id, title, content, tags, version, created_at, updated_at
	`

	var updatedPost entity.Post

	err = tx.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		postId,
		post.Version,
	).Scan(
		&updatedPost.ID,
		&updatedPost.UserId,
		&updatedPost.Title,
		&updatedPost.Content,
		pq.Array(&updatedPost.Tags),
		&updatedPost.Version,
		&updatedPost.CreatedAt,
		&updatedPost.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var exists bool

			err = tx.QueryRowContext(
				ctx,
				`SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`,
				postId,
			).Scan(&exists)

			if err != nil {
				postRepo.logger.Error("Failed to check post existence", map[string]any{
					"post_id": postId,
					"error":   err,
				})

				return nil, fmt.Errorf("check post existence: %w", err)
			}

			if !exists {
				postRepo.logger.Warn("Post not found for update", map[string]any{
					"post_id": postId,
				})

				return nil, domain.ErrNotFound
			}

			postRepo.logger.Warn("Post version conflict", map[string]any{
				"post_id":       postId,
				"requested_ver": post.Version,
			})

			return nil, domain.ErrVersionConflict
		}

		postRepo.logger.Error("Failed to update post", map[string]any{
			"post_id": postId,
			"error":   err,
		})

		return nil, fmt.Errorf("update post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		postRepo.logger.Error("Failed to commit post update", map[string]any{
			"post_id": postId,
			"error":   err,
		})

		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	postRepo.logger.Info("Post updated successfully", map[string]any{
		"post_id":     updatedPost.ID,
		"user_id":     updatedPost.UserId,
		"old_version": post.Version,
		"new_version": updatedPost.Version,
	})

	return &updatedPost, nil
}
