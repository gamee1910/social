package application

import (
	"context"

	"github.com/gamee1910/social/internal/domain"
	"github.com/gamee1910/social/internal/domain/entity"
	"github.com/gamee1910/social/internal/domain/repository"
	"github.com/gamee1910/social/internal/domain/service"
	"github.com/gamee1910/social/internal/interfaces/http/transport/request"
	"github.com/gamee1910/social/internal/interfaces/http/transport/response"
	"github.com/gamee1910/social/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository repository.UserRepository
	logger         *logger.Logger
}

func NewUserService(
	userRepository repository.UserRepository, logger *logger.Logger,
) service.UserService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (userService *userService) GetById(ctx context.Context, userID int64) (*response.UserResponse, error) {

	user, err := userService.userRepository.GetById(ctx, userID)
	if err != nil {
		userService.logger.Error("failed to get user by id", "userID", userID, "error", err)
		return nil, err
	}

	userService.logger.Info("user retrieved successfully", "userID", user.ID)

	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
	}, nil
}

func (userService *userService) Register(ctx context.Context, req request.UserCreationRequest) (*response.UserResponse, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user, err := userService.userRepository.Create(ctx, &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	})

	if err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
	}, nil

}

func (userService *userService) Login(ctx context.Context, req request.UserLoginRequest) (*response.UserResponse, error) {
	user, err := userService.userRepository.GetUserByEmail(ctx, req.Email)
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

	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.In(response.VietnamLocation).Format(response.VietNamTimeFormat),
	}, nil
}
