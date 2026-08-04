package application

import (
	"context"
	"errors"

	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
)

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) service.UserService {
	return &userService{userRepository: userRepository}
}

func (us *userService) GetById(ctx context.Context, userID int64) (*entity.User, error) {
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
	case errors.Is(err, domain.ErrVersionConflict):
		return &domain.ConflictError{Resource: "user"}
	default:
		return err
	}
}
