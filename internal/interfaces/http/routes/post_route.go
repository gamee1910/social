package routes

import "github.com/go-chi/chi/v5"

func RegisterPostRoutes(r chi.Router, deps RouterDependencies) {
	r.Route("/posts", func(r chi.Router) {
		r.Post("/", deps.Container.PostHandler().CreatePostHandler)
		r.Route("/{postID}", func(r chi.Router) {
			r.Get("/", deps.Container.PostHandler().GetPostHandler)
			r.Delete("/", deps.Container.PostHandler().DeletePostHandler)
			r.Patch("/", deps.Container.PostHandler().UpdatePostHandler)
		})
	})
}
