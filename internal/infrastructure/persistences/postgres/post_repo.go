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

func NewPostRepository(
	db *sql.DB,
	logger *logger.Logger,
) repository.PostRepository {
	return &postRepository{
		db:     db,
		logger: logger,
	}
}

func (repository *postRepository) Create(
	ctx context.Context,
	post *entity.Post,
) (*entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	repository.logger.Info("Creating post", map[string]any{
		"user_id": post.UserId,
		"title":   post.Title,
	})

	query := `
	INSERT INTO posts(content, title, user_id, tags)
	VALUES ($1, $2, $3, $4)
	RETURNING id, content, title, user_id, tags, created_at, updated_at, version
	`

	var result entity.Post

	err := repository.db.QueryRowContext(
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
		repository.logger.Error("Failed to create post", map[string]any{
			"user_id": post.UserId,
			"error":   err,
		})

		return nil, fmt.Errorf("create post error: %w", err)
	}

	repository.logger.Info("Post created successfully", map[string]any{
		"post_id": result.ID,
		"user_id": result.UserId,
	})

	return &result, nil
}

func (repository *postRepository) GetById(
	ctx context.Context,
	postId int64,
) (*entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	repository.logger.Info("Getting post by ID", map[string]any{
		"post_id": postId,
	})

	query := `
	SELECT id, user_id, title, content, tags, version, created_at, updated_at
	FROM posts
	WHERE id = $1
	`

	var post entity.Post

	err := repository.db.QueryRowContext(ctx, query, postId).
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
			repository.logger.Warn("Post not found", map[string]any{
				"post_id": postId,
			})

			return nil, domain.ErrNotFound
		}

		repository.logger.Error("Failed to get post", map[string]any{
			"post_id": postId,
			"error":   err,
		})

		return nil, fmt.Errorf("get post: %w", err)
	}

	repository.logger.Info("Post retrieved successfully", map[string]any{
		"post_id": post.ID,
		"user_id": post.UserId,
	})

	return &post, nil
}

func (repository *postRepository) Delete(
	ctx context.Context,
	postID int64,
) error {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	repository.logger.Info("Deleting post", map[string]any{
		"post_id": postID,
	})

	res, err := repository.db.ExecContext(
		ctx,
		`DELETE FROM posts WHERE id = $1`,
		postID,
	)
	if err != nil {
		repository.logger.Error("Failed to delete post", map[string]any{
			"post_id": postID,
			"error":   err,
		})

		return fmt.Errorf("delete post: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		repository.logger.Error("Failed to get affected rows", map[string]any{
			"post_id": postID,
			"error":   err,
		})

		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		repository.logger.Warn("Post not found for deletion", map[string]any{
			"post_id": postID,
		})

		return domain.ErrNotFound
	}

	repository.logger.Info("Post deleted successfully", map[string]any{
		"post_id": postID,
	})

	return nil
}

func (repository *postRepository) Update(
	ctx context.Context,
	postId int64,
	post *entity.Post,
) (*entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	repository.logger.Info("Updating post", map[string]any{
		"post_id": postId,
		"version": post.Version,
		"title":   post.Title,
	})

	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		repository.logger.Error("Failed to begin update transaction", map[string]any{
			"post_id": postId,
			"error":   err,
		})

		return nil, fmt.Errorf("begin transaction error: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			repository.logger.Error("Failed to rollback transaction", map[string]any{
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
				repository.logger.Error("Failed to check post existence", map[string]any{
					"post_id": postId,
					"error":   err,
				})

				return nil, fmt.Errorf("check post existence: %w", err)
			}

			if !exists {
				repository.logger.Warn("Post not found for update", map[string]any{
					"post_id": postId,
				})

				return nil, domain.ErrNotFound
			}

			repository.logger.Warn("Post version conflict", map[string]any{
				"post_id":       postId,
				"requested_ver": post.Version,
			})

			return nil, domain.ErrVersionConflict
		}

		repository.logger.Error("Failed to update post", map[string]any{
			"post_id": postId,
			"error":   err,
		})

		return nil, fmt.Errorf("update post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		repository.logger.Error("Failed to commit post update", map[string]any{
			"post_id": postId,
			"error":   err,
		})

		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	repository.logger.Info("Post updated successfully", map[string]any{
		"post_id":     updatedPost.ID,
		"user_id":     updatedPost.UserId,
		"old_version": post.Version,
		"new_version": updatedPost.Version,
	})

	return &updatedPost, nil
}

func (repository *postRepository) GetFeed(
	ctx context.Context,
	query repository.FeedQuery,
) ([]*entity.Post, error) {

	order := "DESC"
	if query.Sort == "asc" {
		order = "ASC"
	}

	sqlQuery := `
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username, COUNT(c.id) AS comments_count
		FROM posts p 
		    LEFT JOIN comments c ON c.post_id = p.id 
		    LEFT JOIN users u ON p.user_id = u.id
		    JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE ( 
		    p.user_id = $1 OR EXISTS ( 
		    	SELECT 1 FROM followers f WHERE f.user_id = $1 AND f.follower_id = p.user_id 
			) 
		) 
		AND (
			$4 = '' OR p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%' 
		)
		AND (
			$5::text[] = '{}' OR p.tags @> $5
		)
		GROUP BY p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username
		ORDER BY p.created_at ` + order + ` LIMIT $2 OFFSET $3
		`

	ctx, cancel := context.WithTimeout(
		ctx,
		config.QueryTimeoutDuration,
	)
	defer cancel()

	rows, err := repository.db.QueryContext(
		ctx,
		sqlQuery,
		query.UserID,
		query.Limit,
		query.Offset,
		query.Search,
		pq.Array(query.Tags),
	)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			repository.logger.Error(
				"failed to close posts rows",
				"error", err,
			)
		}
	}(rows)

	posts := make([]*entity.Post, 0, query.Limit)

	for rows.Next() {
		var post entity.Post

		err := rows.Scan(
			&post.ID,
			&post.UserId,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentsCount,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
