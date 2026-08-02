package config

import (
	"context"
	"database/sql"
	"time"
)

func DatabaseConnection(add string, maxOpenConnection, maxIdelConnection int, maxIdelTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", add)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConnection)
	db.SetMaxIdleConns(maxIdelConnection)

	duration, err := time.ParseDuration(maxIdelTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(time.Duration(duration))

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancle()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil

}
