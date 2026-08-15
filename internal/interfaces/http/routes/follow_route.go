package routes

import "github.com/go-chi/chi/v5"

func RegisterFollowRoutes(r chi.Router, deps RouterDependencies) {
	r.Route("/follow", func(r chi.Router) {
		r.Route("/{followerID}", func(r chi.Router) {
			r.Put("/follow", deps.Container.FollowerHandler().FollowUser)
			r.Delete("/follow", deps.Container.FollowerHandler().UnfollowUser)
		})
	})
}
