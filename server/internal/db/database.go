package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Connect(dsn string) error {
	var err error
	pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	return nil
}

func Close() {
	pool.Close()
}

func GetDB() *pgxpool.Pool {
	return pool
}
