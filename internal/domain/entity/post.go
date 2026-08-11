package entity

import "time"

type Post struct {
	ID        int64
	Content   string
	Title     string
	UserId    int64
	Tags      []string
	Comment   []Comment
	User      User
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time

	CommentsCount int
}
