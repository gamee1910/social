package service

import (
	"context"

	"github.com/gamee1910/social/internal/domain/entity"
)

type UserService interface {
	GetById(ctx context.Context, userID int64) (*entity.User, error)
}
