package user

import (
	"context"
	"errors"
	"fmt"
	"msg_app/internal/db"
	"time"
)

type AuthPayload struct {
	Username string `form:"username" binding:"required"`
}

func IsUserExist(d *db.Database, username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT EXISTS (SELECT 1 FROM users WHERE user_name = $1)`
	isExist := false

	err := d.Pool.QueryRow(ctx, query, username).Scan(&isExist)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return false, fmt.Errorf("database query timeout out after 3 seconds")
		}

		return false, fmt.Errorf("failed to check username exist: %v", err)
	}

	return isExist, nil
}

func CreateUser(d *db.Database, username string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `INSERT INTO users (user_name, is_online, in_group) VALUES($1, $2, $3) RETURNING id`
	var id int

	err := d.Pool.QueryRow(ctx, query, username, true, false).Scan(&id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return -1, fmt.Errorf("database query timeout out after 10 seconds")
		}

		return -1, fmt.Errorf("failed to insert new users data: %v", err)
	}

	return id, nil
}
