package routes

import "github.com/go-chi/chi/v5"

func RegisterCommentRoutes(r chi.Router, deps RouterDependencies) {
	r.Route("/comments", func(r chi.Router) {

	})
}
