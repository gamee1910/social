package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/service"
	"github.com/gamee1910/social/pkg/logger"
	"github.com/go-chi/chi/v5"
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

func getIDFromParameter(value string, r *http.Request) (int64, error) {
	id := chi.URLParam(r, value)

	valueFromParam, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %w", err)
	}

	return valueFromParam, nil
}
