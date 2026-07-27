package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gamee1910/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config applicationConfig
	store  store.Storage
}

type applicationConfig struct {
	addr           string
	databaseConfig databaseConfig
	env            string
}

type databaseConfig struct {
	addr              string
	maxOpenConnection int
	maxIdelConnection int
	maxIdelTime       string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		//v1/posts/
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})

		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.addr)

	return server.ListenAndServe()
}
