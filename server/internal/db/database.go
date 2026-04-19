package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool   *pgxpool.Pool
	logger *slog.Logger
	dbURL  string
}

func Init(logger *slog.Logger) *Database {
	return &Database{
		Pool:   nil,
		logger: logger,
		dbURL:  os.Getenv("DB_URL"),
	}
}

func (d *Database) Run(ctx context.Context) error {
	if d.Pool == nil {
		pool, err := pgxpool.New(ctx, d.dbURL)
		if err != nil {
			d.logger.Error("unable to create connection postgres pool", "error", err)
			return err
		}

		d.Pool = pool
	}

	defer d.Pool.Close()

	// handle context
	<-ctx.Done()
	return ctx.Err()
}
