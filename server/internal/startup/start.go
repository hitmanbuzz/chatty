package startup

import (
	"context"
	"log/slog"
	"msg_app/internal/db"
	"msg_app/internal/group"
	"msg_app/internal/storage"
	"msg_app/internal/user"
)

const (
	DEFAULT_SERVER_NAME  = "main"
	DEFAULT_SERVER_OWNER = "owner"
)

type Startup struct {
	logger *slog.Logger
}

func Init(logger *slog.Logger) *Startup {
	return &Startup{
		logger: logger,
	}
}

func (s *Startup) Exec(database *db.Database, store *storage.Storage) error {
	u := user.Init(s.logger, database)
	g := group.Init(s.logger, database)

	ctx := context.Background()

	userId, err := u.CreateUser(ctx, DEFAULT_SERVER_OWNER)
	if err != nil {
		return err
	}

	groupId, err := g.CreateGroup(ctx, userId, DEFAULT_SERVER_NAME)
	if err != nil {
		return err
	}

	err = u.JoinGroup(ctx, userId, groupId)
	if err != nil {
		return err
	}

	err = store.InsertUser(userId, DEFAULT_SERVER_OWNER, groupId, DEFAULT_SERVER_NAME)
	if err != nil {
		return err
	}

	s.logger.Info("startup executed succesfully")

	return nil
}
