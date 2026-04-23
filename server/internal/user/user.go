package user

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

type UserData struct {
	Id       int32  `db:"id"`
	Username string `db:"user_name"`
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

func (u *User) CreateUser(pctx context.Context, username string) (int32, error) {
	ctx, cancel := context.WithTimeout(pctx, 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (user_name)
		VALUES($1)
		ON CONFLICT (user_name) DO NOTHING
		RETURNING id
	`
	var userID int32
	err := u.database.Pool.QueryRow(ctx, query, username).Scan(&userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = u.database.Pool.QueryRow(ctx, "SELECT id FROM users WHERE user_name = $1", username).Scan(&userID)
			if err != nil {
				// could check if the user exist here aswell but I doubt it is needed for this project
				return -1, err
			}

			u.logger.Info("user already exist", "username", username)
			return userID, nil
		}

		return -1, err
	}

	u.logger.Info("new user created", "id", userID)

	return userID, nil
}

func (u *User) LoadAllUsers(pctx context.Context) (map[string]UserData, error) {
	ctx, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	query := `SELECT id, user_name FROM users`

	rows, err := u.database.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[UserData])
	if err != nil {
		u.logger.Error("failed to collect all users rows from db")
		return nil, err
	}

	usersMap := arrToMap(users)

	return usersMap, nil
}

func (u *User) JoinGroup(pctx context.Context, userId int32, groupId int32) error {
	ctx, cancel := context.WithTimeout(pctx, 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO group_members (user_id, group_id)
		VALUES($1, $2)
		ON CONFLICT (group_id, user_id) DO NOTHING
	`

	cTag, err := u.database.Pool.Exec(ctx, query, userId, groupId)
	if err != nil {
		if pgError, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgError.Code == "23503" {
				return fmt.Errorf("cannot join group: either the user (%d) or the group (%d) does not exist", userId, groupId)
			}
		}

		return err
	}

	if cTag.RowsAffected() == 0 {
		u.logger.Warn("user is already in the group", "user id", userId, "group id", groupId)
		return nil
	}

	u.logger.Info("user joined group", "user id", userId, "group id", groupId)
	return nil
}

func arrToMap(users []UserData) map[string]UserData {
	result := make(map[string]UserData)

	for _, user := range users {
		result[user.Username] = UserData{
			Id:       user.Id,
			Username: user.Username,
		}
	}

	return result
}
