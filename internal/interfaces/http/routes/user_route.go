package routes

import "github.com/go-chi/chi/v5"

func RegisterUserRoutes(r chi.Router, deps RouterDependencies) {
	r.Route("/users", func(r chi.Router) {
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", deps.Container.UserHandler().GetUserById)
			r.Put("/follow", deps.Container.FollowerHandler().FollowUser)
			r.Delete("/follow", deps.Container.FollowerHandler().UnfollowUser)
		})
	})
}
