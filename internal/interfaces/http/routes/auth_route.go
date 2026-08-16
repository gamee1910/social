package routes

import "github.com/go-chi/chi/v5"

func RegisterAuthRoutes(r chi.Router, deps RouterDependencies) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", deps.Container.AuthHandler().Register)
		r.Post("/login", deps.Container.AuthHandler().Login)
	})
}
