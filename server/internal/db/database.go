package db

import (
	"context"
	"log/slog"
	"os"
	"time"

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
		config, err := pgxpool.ParseConfig(d.dbURL)
		if err != nil {
			return err
		}

		// limit queries to 5 seconds
		// https://www.postgresql.org/docs/current/runtime-config-client.html
		config.ConnConfig.RuntimeParams["statement_timeout"] = "5000"

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			d.logger.Error("unable to create connection postgres pool", "error", err)
			return err
		}

		d.Pool = pool
		d.logger.Info("database connection", "status", "success")
	}

	<-ctx.Done()
	d.logger.Info("database shutdown initiated")

	closed := make(chan struct{})
	go func() {
		d.Pool.Close()
		close(closed)
	}()

	select {
	case <-closed:
		d.logger.Info("database connection pool closed gracefully")
	case <-time.After(5 * time.Second):
		d.logger.Warn("database connection pool forced to close due to timeout")
	}

	return nil
}
