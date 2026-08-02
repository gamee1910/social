package routes

import (
	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/service"
	"github.com/gamee1910/social/pkg/logger"
)

type Handler struct {
	config  config.Config
	service *service.Service
	logger  *logger.Logger
}

func NewHandler(cfg config.Config, service *service.Service, logger *logger.Logger) *Handler {
	return &Handler{
		config:  cfg,
		service: service,
		logger:  logger,
	}
}
