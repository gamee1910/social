package application

import (
	"context"

	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/response"
	"github.com/gamee1910/social/pkg/logger"
)

type userService struct {
	userRepository repository.UserRepository
	logger         *logger.Logger
}

func NewUserService(
	userRepository repository.UserRepository,
	logger *logger.Logger,
) service.UserService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (us *userService) GetById(ctx context.Context, userID int64) (*response.UserResponse, error) {
	user, err := us.userRepository.GetById(ctx, userID)
	if err != nil {
		us.logger.Error("failed to get user by id", "userID", userID, "error", err)
		return nil, err
	}
	us.logger.Info("user retrieved successfully", "userID", user.ID)
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
