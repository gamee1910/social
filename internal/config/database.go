package config

import (
	"context"
	"database/sql"
	"time"
)

const QueryTimeoutDuration = time.Second * 5

func DatabaseConnection(addr string, maxOpenConnection, maxIdleConnection int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConnection)
	db.SetMaxIdleConns(maxIdleConnection)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(time.Duration(duration))

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeoutDuration)

	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil

}
