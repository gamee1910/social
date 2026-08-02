package service

import (
	"context"
	"errors"

	"github.com/gamee1910/social/internal/domain"
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
	user, err := us.userRepository.GetById(ctx, userID)
	if err != nil {
		return nil, translateUserError(err)
	}
	return user, nil
}

func translateUserError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return &domain.NotFoundError{Resource: "user"}
	default:
		return err
	}
}
