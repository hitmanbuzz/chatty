package session

import (
	"database/sql"
	"log/slog"
	"msg_app/internal/db"
	"msg_app/internal/util"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"
)

type Session struct {
	logger *slog.Logger
	sqlDB  *sql.DB
	ginCtx *gin.Engine
}

func Init(logger *slog.Logger, database *db.Database, ginCtx *gin.Engine) *Session {
	if logger == nil || database == nil || ginCtx == nil {
		logger.Error("logger or database or gin is nil")
		return nil
	}

	sqlDB := stdlib.OpenDBFromPool(database.Pool)
	return &Session{
		logger: logger,
		sqlDB:  sqlDB,
		ginCtx: ginCtx,
	}
}

func (s *Session) HandleSession() {
	defer s.sqlDB.Close()

	secret_key := os.Getenv("SESSION_KEY")
	cookie_name := os.Getenv("CLIENT_COOKIE")

	store, err := postgres.NewStore(s.sqlDB, []byte(secret_key))
	if err != nil {
		s.logger.Error("failed to create postgres session key store", "error", err.Error())
		return
	}

	// modify it before production
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   util.DAY_SECS * 30, // 30 days
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	s.ginCtx.Use(sessions.Sessions(cookie_name, store))
}
