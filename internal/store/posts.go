package store

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/gamee1910/social/internal/domain"
	"github.com/lib/pq"
)

type PostsStore struct {
	db *sql.DB
}

func (store *PostsStore) Create(ctx context.Context, post *domain.Post) error {
	query := `
		INSERT INTO posts(content, title, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := store.db.QueryRowContext(
		ctx, query, post.Content, post.Title, post.UserId, pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (store *PostsStore) GetById(ctx context.Context, postId int64) (*domain.Post, error) {
	query := `
		SELECT id, user_id, title, content, tags, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	var post domain.Post

	err := store.db.QueryRowContext(ctx, query, postId).
		Scan(
			&post.ID,
			&post.UserId,
			&post.Title,
			&post.Content,
			pq.Array(&post.Tags),
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

func (store *PostsStore) GetAll(ctx context.Context) ([]*domain.Post, error) {
	query := `SELECT id, user_id, title, content, tags, created_at, updated_at FROM posts`

	rows, err := store.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var posts []*domain.Post
	posts = make([]*domain.Post, 0, 100)

	for rows.Next() {
		var post domain.Post

		err := rows.Scan(
			&post.ID,
			&post.UserId,
			&post.Title,
			&post.Content,
			pq.Array(&post.Tags),
			&post.CreatedAt,
			&post.UpdatedAt,
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
