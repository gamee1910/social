package application

import (
	"context"

	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/pkg/logger"
)

type followService struct {
	followerRepository repository.FollowRepository
	logger             *logger.Logger
}

func NewFollowerService(
	followRepository repository.FollowRepository,
	logger *logger.Logger,
) service.FollowService {
	return &followService{
		followerRepository: followRepository,
		logger:             logger,
	}
}

func (followService *followService) FollowUser(
	ctx context.Context,
	userID, followerID int64,
) error {
	if err := followService.followerRepository.FollowUser(ctx, userID, followerID); err != nil {
		followService.logger.Error(
			"failed to follow user",
			"userID", userID,
			"followerID", followerID,
			"error", err,
		)
		return err
	}
	followService.logger.Info("user followed successfully", "userID", userID, "followerID", followerID)
	return nil
}

func (followService *followService) UnfollowUser(
	ctx context.Context,
	userID, followerID int64,
) error {
	if err := followService.followerRepository.UnfollowUser(ctx, userID, followerID); err != nil {
		followService.logger.Error(
			"failed to unfollow user",
			"userID", userID,
			"followerID", followerID,
			"error", err,
		)
		return err
	}
	followService.logger.Info("user unfollowed successfully", "userID", userID, "followerID", followerID)
	return nil
}
