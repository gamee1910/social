package service

import (
	"context"
)

type FollowerService interface {
	FollowUser(ctx context.Context, userID, followerID int64) error
	UnfollowUser(ctx context.Context, userID, followerID int64) error
}
