package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
	dbURL  string
}

func Init(logger *slog.Logger) *Database {
	return &Database{
		pool:   nil,
		logger: logger,
		dbURL:  os.Getenv("DB_URL"),
	}
}

func (d *Database) Run(ctx context.Context) error {
	if d.pool == nil {
		pool, err := pgxpool.New(ctx, d.dbURL)
		if err != nil {
			d.logger.Error("unable to create connection postgres pool", "error", err)
			return err
		}

		d.pool = pool
	}

	defer d.pool.Close()

	// handle context
	<-ctx.Done()
	return ctx.Err()
}
