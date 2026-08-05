package application

import (
	"context"

	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/response"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) service.UserService {
	return &userService{userRepository: userRepository}
}

func (us *userService) GetById(ctx context.Context, userID int64) (*response.UserResponse, error) {
	user, err := us.userRepository.GetById(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}
