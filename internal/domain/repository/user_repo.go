package repository

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetById(ctx context.Context, userID int64) (*entity.User, error)
}
