package auth

import (
	"errors"
	"log/slog"
	"msg_app/internal/db"
	"msg_app/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthPayload struct {
	Username  string `form:"username" binding:"required"`
	Groupname string `form:"groupname" default:"default"`
}

func NewPayload(username string, groupname string) AuthPayload {
	return AuthPayload{
		Username:  username,
		Groupname: groupname,
	}
}

type Auth struct {
	logger   *slog.Logger
	database *db.Database
}

func NewAuth(logger *slog.Logger, database *db.Database) *Auth {
	return &Auth{
		logger:   logger,
		database: database,
	}
}

func (a *Auth) Authenticate(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	var payload AuthPayload

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	u := user.Init(a.logger, a.database)

	userid, err := u.CreateUser(reqCtx, payload.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusOK, gin.H{"status": "user already exist with this username"})
			return
		}

		a.logger.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	a.logger.Info("user created", "id", userid)
	ctx.JSON(http.StatusOK, gin.H{"status": "user created successfully"})

}
