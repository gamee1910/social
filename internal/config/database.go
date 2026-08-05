package config

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

const QueryTimeoutDuration = time.Second * 5

func DatabaseConnection(cfg *Configuration) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Database.DatabaseDSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	db.SetConnMaxIdleTime(cfg.Database.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
