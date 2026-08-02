package store

import (
	"context"
	"database/sql"

	"github.com/gamee1910/social/internal/domain/entity"
)

type CommentsStore struct {
	db *sql.DB
}

func NewCommentsStore(db *sql.DB) *CommentsStore {
	return &CommentsStore{db: db}
}

func (store *CommentsStore) GetByPostId(ctx context.Context, postId int64) ([]entity.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.id
		FROM comments c JOIN users u on c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	rows, err := store.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

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
			return nil, err
		}

		comments = append(comments, cmt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
