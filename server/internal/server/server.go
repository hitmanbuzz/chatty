package server

import (
	"context"
	"fmt"
	"log/slog"
	"msg_app/internal/db"
	"msg_app/internal/startup"
	"msg_app/internal/storage"
	"msg_app/internal/ws"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	hostIP string
	logger *slog.Logger

	server   *gin.Engine
	w        *ws.Websocket
	database *db.Database
	storage  *storage.Storage
	start    *startup.Startup

	group_count uint
	user_count  uint
}

func Init(logger *slog.Logger) *Server {
	return &Server{
		hostIP:      os.Getenv("TCP_SERVER_IP"),
		server:      gin.Default(),
		w:           ws.Init(logger),
		database:    db.Init(logger),
		storage:     storage.InitStorage(logger),
		start:       startup.Init(logger),
		group_count: 0,
		user_count:  0,
		logger:      logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.hostIP,
		Handler: s.server,
	}

	if err := s.database.Connect(ctx); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	s.Exec()
	s.Routes()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server crashed: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := s.database.Run(gCtx); err != nil {
			return fmt.Errorf("database crashed: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server forced to shutdown with error: %w", err)
		}
		return nil
	})

	return g.Wait()
}

func (s *Server) Exec() {
	err := s.start.Exec(s.database, s.storage)
	if err != nil {
		s.logger.Error("failed to do startup proces", "error", err)
		return
	}
}

func (s *Server) Routes() {
	s.server.GET("/echo", s.w.Messaging)

	// s.server.POST("/auth", func(ctx *gin.Context) {
	// 	u := user.Init(s.logger, s.database)
	// 	u.AuthUser(ctx)
	// })
}
