package session

import (
	"database/sql"
	"log/slog"
	"msg_app/internal/db"

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
	if database == nil || ginCtx == nil {
		logger.Error("database or gin is nil")
		return nil
	}

	sqlDB := stdlib.OpenDBFromPool(database.Pool)
	return &Session{
		logger: logger,
		sqlDB:  sqlDB,
		ginCtx: ginCtx,
	}
}

func (s *Session) StoreSession() {
	defer s.sqlDB.Close()

	store, err := postgres.NewStore(s.sqlDB, []byte("secret-key"))
	if err != nil {
		s.logger.Error("failed to create postgres session key store", "error", err.Error())
		return
	}

	store.Options(sessions.Options{
		MaxAge: 86400 * 30, // 30 days
		Path:   "/",
	})

	s.ginCtx.Use(sessions.Sessions("auth_session", store))
}
