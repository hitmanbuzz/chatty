package app

import (
	"context"
	"log/slog"
	"msg_app/internal/server"
	"os"
	"os/signal"
)

type App struct {
	logger      *slog.Logger
	chat_server *server.Server
}

func Init() *App {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	return &App{
		logger:      logger,
		chat_server: server.Init(logger),
	}
}

func (a *App) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := a.chat_server.Run(ctx); err != nil {
		a.logger.Error("application shutdown with error")
	} else {
		a.logger.Info("application shutdown gracefully")
	}
}
