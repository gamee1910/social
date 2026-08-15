package entity

import "time"

type Follower struct {
	UserID     int64
	FollowerID int64
	CreatedAt  time.Time
}
