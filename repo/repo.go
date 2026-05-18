package repo

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	Pool *pgxpool.Pool
}

var NRA = errors.New("No rows affected prolem with entry")

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	connString := os.Getenv("DB_URL")
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		log.Println("DB is not connecting")
		return nil, err
	}
	return pool, nil
}
