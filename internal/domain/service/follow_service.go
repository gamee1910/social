package service

import (
	"context"
)

type FollowService interface {
	FollowUser(ctx context.Context, userID, followerID int64) error
	UnfollowUser(ctx context.Context, userID, followerID int64) error
}
