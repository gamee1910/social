package repository

import "context"

type FollowerRepository interface {
	FollowUser(ctx context.Context, userID, followerID int64) error
	UnfollowUser(ctx context.Context, userID, followerID int64) error
}
