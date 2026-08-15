package entity

import "time"

type Comment struct {
	ID        int64
	PostId    int64
	UserId    int64
	Content   string
	User      User
	CreatedAt time.Time
}
