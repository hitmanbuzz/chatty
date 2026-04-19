package server

import (
	"context"
	"fmt"
	"log/slog"
	"msg_app/internal/db"
	"msg_app/internal/user"
	"msg_app/internal/ws"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	hostIP string
	logger *slog.Logger

	server   *gin.Engine
	w        *ws.Websocket
	database *db.Database

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
		database:    db.Init(logger),
		groups:      make([]uint, 0),
		users:       make([]uint, 0),
		group_count: 0,
		user_count:  0,
		logger:      logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.Routes()

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

	go s.database.Run(ctx)

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

func (s *Server) Routes() {
	s.server.GET("/echo", s.w.Messaging)

	s.server.POST("/auth", func(ctx *gin.Context) {
		var payload user.AuthPayload

		if err := ctx.ShouldBind(&payload); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		s.logger.Info("received from client", "username", payload.Username)

		isExist, err := user.IsUserExist(s.database, payload.Username)
		if err != nil {
			s.logger.Error(err.Error())
		} else {
			switch isExist {
			case true:
				ctx.JSON(http.StatusOK, gin.H{"status": "user already exist with this username"})
			case false:
				id, err := user.CreateUser(s.database, payload.Username)
				if err != nil {
					s.logger.Error(err.Error())
					ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
					return
				}

				s.logger.Info("user created", "id", id)
				ctx.JSON(http.StatusOK, gin.H{"status": "user created successfully"})
			}
		}
	})
}
