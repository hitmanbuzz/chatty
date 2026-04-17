package app

import (
	"context"
	"log/slog"
	"msg_app/internal/db"
	"msg_app/internal/server/tcp_server"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

type App struct {
	logger      *slog.Logger
	chat_server *tcp_server.Server
}

func Init() *App {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	return &App{
		logger:      logger,
		chat_server: tcp_server.Init(logger),
	}
}

func (a *App) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return a.chat_server.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		a.logger.Error("application shutdown with error")
	} else {
		a.logger.Error("application shutdown gracefully")
	}
}
