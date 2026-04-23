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

type SignupPayload struct {
	Username string `form:"username" binding:"required"`
}

type LoginPayload struct {
	Username  string `form:"username" binding:"required"`
	Groupname string `form:"groupname" default:"default"`
}

type Authenticate struct {
	logger   *slog.Logger
	database *db.Database
}

func NewLogin(logger *slog.Logger, database *db.Database) *Authenticate {
	return &Authenticate{
		logger:   logger,
		database: database,
	}
}

func (a *Authenticate) SignupUser(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	var payload SignupPayload

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	u := user.Init(a.logger, a.database)
	userID, err := u.CreateUser(reqCtx, payload.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusOK, gin.H{"status": "user already exist with this username"})
			return
		}

		a.logger.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	a.logger.Info("user signup successfully", "user id", userID, "username", payload.Username)
	ctx.JSON(http.StatusOK, gin.H{"status": "user created successfully"})
}

// func (a *Authenticate) LoginUser(ctx *gin.Context) {
// 	reqCtx := ctx.Request.Context()

// 	var payload LoginPayload

// 	if err := ctx.ShouldBind(&payload); err != nil {
// 		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
// 		return
// 	}

// 	u := user.Init(a.logger, a.database)

// 	userid, err := u.CreateUser(reqCtx, payload.Username)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			ctx.JSON(http.StatusOK, gin.H{"status": "user already exist with this username"})
// 			return
// 		}

// 		a.logger.Error(err.Error())
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
// 		return
// 	}

// 	a.logger.Info("user authenticated successfully", "id", userid, "username", payload.Username)
// 	ctx.JSON(http.StatusOK, gin.H{"status": "user created successfully"})
// }
