package user

import (
	"context"
	"errors"
	"log/slog"
	"msg_app/internal/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type authPayload struct {
	Username string `form:"username" binding:"required"`
}

type User struct {
	logger   *slog.Logger
	database *db.Database
}

func Init(logger *slog.Logger, database *db.Database) *User {
	return &User{
		logger:   logger,
		database: database,
	}
}

func (u *User) AuthUser(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	var payload authPayload

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	id, err := createUser(reqCtx, u.database, payload.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusOK, gin.H{"status": "user already exist with this username"})
			return
		}

		u.logger.Error(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	u.logger.Info("user created", "id", id)
	ctx.JSON(http.StatusOK, gin.H{"status": "user created successfully"})
}

func createUser(pctx context.Context, d *db.Database, username string) (int, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (user_name, in_group)
		VALUES($1, $2)
		ON CONFLICT (user_name) DO NOTHING
		RETURNING id
	`
	var id int
	err := d.Pool.QueryRow(ctx, query, username, false).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}
