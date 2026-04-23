package group

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"msg_app/internal/db"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type GroupData struct {
	Id         int32  `db:"id"`
	GroupName  string `db:"group_name"`
	TotalUsers int32  `db:"total_users"`
	OwnerId    int32  `db:"owner_id"`
}

type Group struct {
	logger   *slog.Logger
	database *db.Database
}

func Init(logger *slog.Logger, database *db.Database) *Group {
	return &Group{
		logger:   logger,
		database: database,
	}
}

func (g *Group) CreateGroup(pctx context.Context, owner_id int32, groupname string) (int32, error) {
	ctx, cancel := context.WithTimeout(pctx, 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO groups (group_name, owner_id)
		VALUES($1, $2)
		ON CONFLICT (group_name) DO UPDATE
		SET group_name = groups.group_name
		RETURNING id
	`
	var groupid int32
	err := g.database.Pool.QueryRow(ctx, query, groupname, owner_id).Scan(&groupid)

	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23503" {
				return -1, fmt.Errorf("invalid references: the user id doesn't exist")
			}
		}

		if errors.Is(err, pgx.ErrNoRows) {
			g.logger.Warn("group with that name already exist")
			return groupid, nil
		}

		return -1, fmt.Errorf("failed to insert default group: %w", err)
	}

	g.logger.Info("new group created", "id", groupid)

	return groupid, nil
}

func (g *Group) GetGroupName(pctx context.Context, groupID int32) (string, error) {
	ctx, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT group_name from groups
		WHERE id = $1
	`

	var groupName string
	err := g.database.Pool.QueryRow(ctx, query, groupID).Scan(&groupName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("no group found, group id: %d", groupID)
		}
	}

	return groupName, nil
}

// Return (groupID, groupName, err)
func (g *Group) FetchUserGroup(pctx context.Context, userID int32) (int32, string, error) {
	ctx, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT group_id FROM group_members
		WHERE user_id = $1
	`

	var groupID int32
	err := g.database.Pool.QueryRow(ctx, query, userID).Scan(&groupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			var username string
			_ = g.database.Pool.QueryRow(ctx, "SELECT user_name from users WHERE id = $1", userID).Scan(&username)
			return -1, "", fmt.Errorf("user is not part of this group [ user id: %d | username: %s ]", userID, username)
		}

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23503" {
				return -1, "", fmt.Errorf("invalid references: the user id doesn't exist")
			}
		}
	}

	groupName, err := g.GetGroupName(ctx, groupID)
	if err != nil {
		return groupID, "", err
	}

	return groupID, groupName, nil
}
