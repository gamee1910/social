package repository

import "context"

type FollowRepository interface {
	FollowUser(ctx context.Context, userID, followerID int64) error
	UnfollowUser(ctx context.Context, userID, followerID int64) error
}
