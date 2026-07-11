package main

import (
	"log"

	"github.com/gamee1910/social/internal/db"
	"github.com/gamee1910/social/internal/env"
	"github.com/gamee1910/social/internal/store"
)

func main() {
	cfg := applicationConfig{
		addr: env.GetString("ADDR", ":8080"),
		databaseConfig: databaseConfig{
			addr:              env.GetString("DB_ADDR", "postgres://admin:password@localhost/social?sslmode=disable"),
			maxOpenConnection: env.GetInt("DB_MAX_OPEN_CONS", 30),
			maxIdelConnection: env.GetInt("DB_MAX_IDEL_CONS", 30),
			maxIdelTime:       env.GetString("DB_MAX_IDEL_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.databaseConfig.addr,
		cfg.databaseConfig.maxOpenConnection,
		cfg.databaseConfig.maxIdelConnection,
		cfg.databaseConfig.maxIdelTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
