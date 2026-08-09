package application

import (
	"context"

	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/pkg/logger"
)

type followerService struct {
	followerRepository repository.FollowerRepository
	logger             *logger.Logger
}

func NewFollowerService(
	followerRepository repository.FollowerRepository,
	logger *logger.Logger,
) service.FollowerService {
	return &followerService{
		followerRepository: followerRepository,
		logger:             logger,
	}
}

func (f *followerService) FollowUser(
	ctx context.Context,
	userID, followerID int64,
) error {
	if err := f.followerRepository.FollowUser(ctx, userID, followerID); err != nil {
		f.logger.Error(
			"failed to follow user",
			"userID", userID,
			"followerID", followerID,
			"error", err,
		)
		return err
	}
	f.logger.Info("user followed successfully", "userID", userID, "followerID", followerID)
	return nil
}

func (f *followerService) UnfollowUser(
	ctx context.Context,
	userID, followerID int64,
) error {
	if err := f.followerRepository.UnfollowUser(ctx, userID, followerID); err != nil {
		f.logger.Error(
			"failed to unfollow user",
			"userID", userID,
			"followerID", followerID,
			"error", err,
		)
		return err
	}
	f.logger.Info("user unfollowed successfully", "userID", userID, "followerID", followerID)
	return nil
}
