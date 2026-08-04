package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gamee1910/social/internal/config"
	"github.com/gamee1910/social/internal/env"
	"github.com/gamee1910/social/internal/routes"
	"github.com/gamee1910/social/internal/service"
	"github.com/gamee1910/social/internal/store"
	"github.com/gamee1910/social/pkg/logger"
)

func main() {
	cfg := config.Config{
		Addr: env.GetString("ADDR", ":8080"),
		Database: config.DatabaseConfig{
			Addr:               env.GetString("DB_ADDR", "postgres://admin:password@localhost/social?sslmode=disable"),
			MaxOpenConnections: env.GetInt("DB_MAX_OPEN_CONS", 30),
			MaxIdleConnections: env.GetInt("DB_MAX_IDLE_CONS", 30),
			MaxIdleTime:        env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Env: env.GetString("ENV", "development"),
	}

	appLogger := logger.NewLogger(cfg.Env)
	defer appLogger.Sync()

	database, err := config.DatabaseConnection(
		cfg.Database.Addr,
		cfg.Database.MaxOpenConnections,
		cfg.Database.MaxIdleConnections,
		cfg.Database.MaxIdleTime,
	)

	if err != nil {
		appLogger.Fatal(err)
	}

	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			appLogger.Error("failed to close database", "error", err)
		}
	}(database)

	storage := store.NewStorage(database)

	svc := service.NewService(storage)

	handler := routes.NewHandler(cfg, svc, appLogger)

	appLogger.Infof(
		"starting server port=%s env=%s",
		cfg.Addr,
		cfg.Env,
	)

	log.Fatal(run(cfg, handler.Mount()))
}

func run(cfg config.Config, handler http.Handler) error {
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	return server.ListenAndServe()
}
