package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}

func NewPostsStore(db *sql.DB) *PostsStore {
	return &PostsStore{db: db}
}

func (store *PostsStore) Create(ctx context.Context, post *entity.Post) error {
	query := `
		INSERT INTO posts(content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at, version
	`

	err := store.db.QueryRowContext(
		ctx, query, post.Content, post.Title, post.UserId, pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		return err
	}

	return nil
}

func (store *PostsStore) GetById(ctx context.Context, postId int64) (*entity.Post, error) {
	query := `
		SELECT id, user_id, title, content, tags, version, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	var post entity.Post

	err := store.db.QueryRowContext(ctx, query, postId).
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound

		default:
			return nil, err
		}
	}

	return &post, nil
}

func (store *PostsStore) Delete(ctx context.Context, postID int64) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("rollback failed: %v", err)
		}
	}()

	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM comments WHERE post_id = $1`,
		postID,
	)

	if err != nil {
		return fmt.Errorf("delete comments: %w", err)
	}

	res, err := tx.ExecContext(
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
		return ErrNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (store *PostsStore) Update(ctx context.Context, postId int64, post *entity.Post) (*entity.Post, error) {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction error: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("rollback failed: %v", err)
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
				return nil, ErrNotFound
			}

			return nil, ErrVersionConflict
		}
		return nil, fmt.Errorf("update post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return &updatedPost, nil
}
