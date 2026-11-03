package store

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/heniums/habitx-server/db"
)

// New opens a connection pool and returns sqlc-backed queries.
func New(ctx context.Context) (*db.Queries, *pgxpool.Pool, error) {
	dsn := os.Getenv("HABITX_DB_DSN")
	if dsn == "" {
		dsn = "postgres://habitx:habitx@localhost:5433/habitx?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, err
	}
	return db.New(pool), pool, nil
}
