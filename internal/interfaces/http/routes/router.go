package routes

import (
	"database/sql"
	"net/http"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/infrastructure/di"
	"github.com/gamee1910/social/internal/utils"
	"github.com/gamee1910/social/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterDependencies struct {
	Config    *config.Configuration
	Container *di.Container
	Logger    *logger.Logger
}

func SetupRouter(cfg *config.Configuration, db *sql.DB, container *di.Container, logger *logger.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			logger.Error("Health check: DB ping failed", "error", err)
			utils.InternalServerError(w, r, err)
			return
		}

		response := map[string]string{
			"status": "ok",
		}
		if err := utils.ResponseJSON(w, http.StatusOK, response); err != nil {
			logger.Error("Health check: failed to write response", "error", err)
		}
	})

	deps := RouterDependencies{
		Config:    cfg,
		Container: container,
		Logger:    logger,
	}

	r.Route("/api/v1", func(r chi.Router) {
		RegisterPostRoutes(r, deps)
		RegisterUserRoutes(r, deps)
		RegisterCommentRoutes(r, deps)
		RegisterFollowRoutes(r, deps)
	})

	return r
}
