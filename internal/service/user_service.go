package service

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/store"
)

type UserService struct {
	userRepository store.UserRepository
}

func NewUserService(userRepository store.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) GetById(ctx context.Context, userID int64) (*entity.User, error) {
	return us.userRepository.GetById(ctx, userID)
}
