package routes

import (
	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/store"
	"github.com/gamee1910/social/pkg/logger"
)

type Handler struct {
	config config.Config
	store  store.Storage
	logger *logger.Logger
}

func NewHandler(cfg config.Config, store store.Storage, logger *logger.Logger) *Handler {
	return &Handler{
		config: cfg,
		store:  store,
		logger: logger,
	}
}
