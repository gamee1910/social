package main

import (
	"database/sql"
	"log"

	"github.com/gamee1910/social/internal/db"
	"github.com/gamee1910/social/internal/env"
	"github.com/gamee1910/social/internal/store"
)

const version = "0.0.1"

func main() {
	cfg := applicationConfig{
		addr: env.GetString("ADDR", ":8080"),
		databaseConfig: databaseConfig{
			addr:              env.GetString("DB_ADDR", "postgres://admin:password@localhost/social?sslmode=disable"),
			maxOpenConnection: env.GetInt("DB_MAX_OPEN_CONS", 30),
			maxIdelConnection: env.GetInt("DB_MAX_IDEL_CONS", 30),
			maxIdelTime:       env.GetString("DB_MAX_IDEL_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	database, err := db.New(
		cfg.databaseConfig.addr,
		cfg.databaseConfig.maxOpenConnection,
		cfg.databaseConfig.maxIdelConnection,
		cfg.databaseConfig.maxIdelTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer func(database *sql.DB) {
		err := database.Close()
		if err != nil {
			log.Println("database connection pool established")
		}
	}(database)

	storage := store.NewStorage(database)

	app := &application{
		config: cfg,
		store:  storage,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
