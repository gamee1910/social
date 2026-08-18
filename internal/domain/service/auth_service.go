package service

import (
	"context"

	"github.com/gamee1910/social/internal/interfaces/http/handler/request"
	"github.com/gamee1910/social/internal/interfaces/http/handler/response"
)

type AuthService interface {
	Login(ctx context.Context, req request.UserLoginRequest) (*response.UserResponse, error)
	Register(ctx context.Context, req request.UserCreationRequest) (*response.UserResponse, error)
}
