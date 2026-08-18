package application

import (
	"context"
	"fmt"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/handler/request"
	"github.com/gamee1910/social/internal/interfaces/http/handler/response"
	"github.com/gamee1910/social/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository repository.UserRepository
	jwtConfig      config.JWTConfig
}

func (authService *authService) Login(ctx context.Context, req request.UserLoginRequest) (*response.UserResponse, error) {
	user, err := authService.userRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		user.Password,
		[]byte(req.Password),
	)

	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	token, err := jwt.GenerateToken(
		user.ID, user.Username, user.Email, authService.jwtConfig.SecretKey, authService.jwtConfig.Expiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
		Token:     token,
	}, nil
}

func (authService *authService) Register(ctx context.Context, req request.UserCreationRequest) (*response.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user, err := authService.userRepository.Create(ctx, &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	})

	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(
		user.ID, user.Username, user.Email, authService.jwtConfig.SecretKey, authService.jwtConfig.Expiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
		Token:     token,
	}, nil
}

func NewAuthService(
	userRepository repository.UserRepository,
	jwtConfig config.JWTConfig,
) service.AuthService {
	return &authService{
		userRepository: userRepository,
		jwtConfig:      jwtConfig,
	}
}
