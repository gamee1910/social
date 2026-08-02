package routes

import (
	"net/http"

	"github.com/gamee1910/social/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const version = "0.0.1"

func (h *Handler) Mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", h.healthCheck)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", h.createPostHandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Get("/", h.getPostHandler)
				r.Delete("/", h.deletePostHandler)
				r.Patch("/", h.updatePostHandler)
			})

		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", h.getUserById)
			})
		})
	})

	return r
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"env":     h.config.Env,
		"version": version,
	}

	if err := config.WriteJSON(w, http.StatusOK, data); err != nil {
		InternalServerError(w, r, err)
	}
}
