package service

import (
	"context"

	"github.com/gamee1910/social/internal/interfaces/http/transport/request"
	"github.com/gamee1910/social/internal/interfaces/http/transport/response"
)

type UserService interface {
	GetById(ctx context.Context, userID int64) (*response.UserResponse, error)
	Register(ctx context.Context, req request.UserCreationRequest) (*response.UserResponse, error)
	Login(ctx context.Context, req request.UserLoginRequest) (*response.UserResponse, error)
}
