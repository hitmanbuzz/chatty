package server

import (
	"context"
	"fmt"
	"log/slog"
	"msg_app/internal/ws"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	hostIP string
	logger *slog.Logger

	server *gin.Engine
	w      *ws.Websocket

	groups      []uint
	users       []uint
	group_count uint
	user_count  uint
}

func Init(logger *slog.Logger) *Server {
	return &Server{
		hostIP:      os.Getenv("TCP_SERVER_IP"),
		server:      gin.Default(),
		w:           ws.Init(logger),
		groups:      make([]uint, 0),
		users:       make([]uint, 0),
		group_count: 0,
		user_count:  0,
		logger:      logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.server.GET("/echo", s.w.Echo)

	server := &http.Server{
		Addr:    s.hostIP,
		Handler: s.server,
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server crashed: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server forced to shutdown with error: %w", err)
		}

		return nil
	}
}
