package service

import (
	"context"

	"github.com/gamee1910/social/internal/interfaces/http/handler/response"
)

type UserService interface {
	GetById(ctx context.Context, userID int64) (*response.UserResponse, error)
}
