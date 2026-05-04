package server

import (
	"context"
	"fmt"
	"log/slog"
	"msg_app/internal/auth"
	"msg_app/internal/db"
	"msg_app/internal/session"
	"msg_app/internal/startup"
	"msg_app/internal/storage"
	"msg_app/internal/ws"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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
	sess     *session.Session
}

func Init(logger *slog.Logger) *Server {
	server := gin.Default()
	storage := storage.InitStorage(logger)
	database := db.Init(logger)
	sess := session.Init(logger, database, server) // session
	start := startup.Init(logger, storage)         // startup
	w := ws.Init(logger)                           // websocket

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	return &Server{
		hostIP:   os.Getenv("SERVER_IP"),
		server:   server,
		w:        w,
		database: database,
		storage:  storage,
		start:    start,
		sess:     sess,
		logger:   logger,
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

	if s.sess == nil {
		return fmt.Errorf("session manager is nil")
	}

	s.sess.HandleSession()
	s.logger.Info("session handling running in background")

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
	s.server.GET("/user/:username", func(ctx *gin.Context) {
		username := ctx.Param("username")

		isUserOnline := s.storage.IsUserOnline(username)
		ctx.JSON(http.StatusOK, gin.H{
			"is_online": isUserOnline,
		})
	})

	s.server.POST("/register", func(ctx *gin.Context) {
		a := auth.NewAuth(s.logger, s.database)
		a.Register(ctx)
	})

	s.server.POST("/login", func(ctx *gin.Context) {
		a := auth.NewAuth(s.logger, s.database)
		a.Login(ctx, s.storage)
	})

	protected := s.server.Group("/api")
	protected.Use(auth.AuthRequired())
	{
		protected.GET("/chat")
	}

	// s.server.GET("/echo", s.w.Messaging)

	// s.server.POST("/auth", func(ctx *gin.Context) {
	// 	a := auth.NewAuth(s.logger, s.database)
	// 	a.Authenticate(ctx)
	// })
}
