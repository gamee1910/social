package repository

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
)

type UserRepository interface {
	GetById(ctx context.Context, userID int64) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
