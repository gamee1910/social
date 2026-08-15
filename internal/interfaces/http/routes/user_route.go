package routes

import "github.com/go-chi/chi/v5"

func RegisterUserRoutes(r chi.Router, deps RouterDependencies) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/register", deps.Container.UserHandler().Register)
		r.Post("/login", deps.Container.UserHandler().Login)
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", deps.Container.UserHandler().GetUserById)
		})
	})
}
