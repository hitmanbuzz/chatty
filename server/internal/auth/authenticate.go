package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"msg_app/internal/db"
	"msg_app/internal/session"
	"msg_app/internal/storage"
	"msg_app/internal/user"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthPayload struct {
	Username string `json:"username" binding:"required"`
}

type Authenticate struct {
	logger   *slog.Logger
	database *db.Database
	sess     *session.Session
}

func NewAuth(logger *slog.Logger, database *db.Database) *Authenticate {
	return &Authenticate{
		logger:   logger,
		database: database,
	}
}

// TODO: handle username length limit
func (a *Authenticate) Register(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	var payload AuthPayload

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := user.Init(a.logger, a.database)
	userID, err := u.CreateUser(reqCtx, payload.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusFound, gin.H{"status": "user already exist with this username"})
			return
		}

		a.logger.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	a.logger.Info("user signup successfully", "user id", userID, "username", payload.Username)
	ctx.JSON(http.StatusOK, gin.H{"status": "user created successfully"})
}

// TODO: handle username length limit
func (a *Authenticate) LoginUser(ctx *gin.Context, store *storage.Storage) {
	if store == nil {
		a.logger.Error("storage is nil in authentication")
		return
	}

	var payload AuthPayload

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.Lock()
	u, found := store.Users[payload.Username]
	if !found {
		store.Unlock()
		ctx.JSON(http.StatusNotFound, gin.H{"status": fmt.Sprintf("user with username %s not found", payload.Username)})
		return
	}

	if u.IsOnline {
		store.Unlock()
		ctx.JSON(http.StatusConflict, gin.H{"error": "user is already connected to another device"})
		return
	}

	store.Users[u.Username].IsOnline = true
	store.Unlock()

	session := sessions.Default(ctx)
	session.Clear()

	session.Set("userID", u.UserID)
	session.Set("username", u.Username)

	if err := session.Save(); err != nil {
		store.Lock()
		store.Users[u.Username].IsOnline = false
		store.Unlock()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session"})
		return
	}

	a.logger.Info("user authenticated successfully", "id", u.UserID, "username", u.Username)
	ctx.JSON(http.StatusOK, gin.H{"status": "login successfully", "user": u.Username})
}

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := sessions.Default(ctx)
		userID := s.Get("userID")

		if userID == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "you must be logged in to access this area",
			})
			return
		}

		ctx.Set("currUserID", userID)
		ctx.Next()
	}
}
