package store

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gamee1910/social/internal/domain"
)

type CommentsStore struct {
	db *sql.DB
}

func (store *CommentsStore) GetByPostId(ctx context.Context, postId int64) ([]domain.Comment, error) {
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

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Get Post By Id error date time [%s]", time.DateTime)
		}
	}(rows)

	var comments []domain.Comment

	for rows.Next() {
		var cmt domain.Comment
		cmt.User = domain.User{}

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

	return comments, nil
}
