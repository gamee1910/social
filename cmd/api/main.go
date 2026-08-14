package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/infrastructure/di"
	"github.com/gamee1910/social/internal/interfaces/http/routes"
	"github.com/gamee1910/social/pkg/logger"

	_ "github.com/gamee1910/social/docs"
)

// @title Social API
// @version 1.0
// @description Social network REST API
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg := config.Load()
	log := logger.NewLogger(cfg.Application.Env)
	defer log.Sync()

	db, err := config.DatabaseConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer closeDB(db, log)

	container := di.NewContainer(cfg, db, log)

	router := routes.SetupRouter(cfg, db, container, log)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  time.Minute,
	}

	go func() {
		log.Infof(
			"starting application [%s] port [%s] env [%s]",
			cfg.Application.Name,
			cfg.Server.Port,
			cfg.Application.Env,
		)

		var err error

		if cfg.Server.TLS.Mode == "enabled" {
			err = server.ListenAndServeTLS(
				cfg.Server.TLS.CertFile,
				cfg.Server.TLS.KeyFile,
			)
		} else {
			err = server.ListenAndServe()
		}

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	gracefulShutdown(server, log)
}

func gracefulShutdown(server *http.Server, log *logger.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Info("server stopped")
}

func closeDB(db *sql.DB, log *logger.Logger) {
	if err := db.Close(); err != nil {
		log.Error("failed to close database", "error", err)
	}
}
