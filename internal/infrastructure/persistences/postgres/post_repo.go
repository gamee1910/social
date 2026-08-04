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
	"github.com/lib/pq"
)

type postRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) repository.PostRepository {
	return &postRepository{db: db}
}

func (postRepo *postRepository) Create(ctx context.Context, post *entity.Post) error {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	query := `
		INSERT INTO posts(content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at, version
	`

	err := postRepo.db.QueryRowContext(
		ctx, query, post.Content, post.Title, post.UserId, pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		return fmt.Errorf("create post: [%w]", err)
	}

	return nil
}

func (postRepo *postRepository) GetById(ctx context.Context, postId int64) (*entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

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
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get post: [%w]", err)
	}

	return &post, nil
}

func (postRepo *postRepository) Delete(ctx context.Context, postID int64) error {

	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	res, err := postRepo.db.ExecContext(
		ctx,
		`DELETE FROM posts WHERE id = $1`,
		postID,
	)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (postRepo *postRepository) Update(ctx context.Context, postId int64, post *entity.Post) (*entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, config.QueryTimeoutDuration)
	defer cancel()

	tx, err := postRepo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction error: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			// rollback failure is non-critical; the original error already propagates
			_ = err
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
			err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`, postId).Scan(&exists)
			if err != nil {
				return nil, fmt.Errorf("check post existence: %w", err)
			}

			if !exists {
				return nil, domain.ErrNotFound
			}

			return nil, domain.ErrVersionConflict
		}
		return nil, fmt.Errorf("update post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &updatedPost, nil
}
