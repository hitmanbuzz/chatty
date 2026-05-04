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
	store  *storage.Storage
}

func Init(logger *slog.Logger, store *storage.Storage) *Startup {
	return &Startup{
		logger: logger,
		store:  store,
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

	err = store.InsertUser(userId, DEFAULT_SERVER_OWNER)
	if err != nil {
		return err
	}

	users, err := u.LoadAllUsers(ctx)
	if err != nil {
		return err
	}

	for _, us := range users {
		err = s.store.InsertUser(us.Id, us.Username)
		if err != nil {
			s.logger.Error("failed to insert use to storage memory", "error", err)
			continue
		}
	}

	s.logger.Info("startup executed succesfully")
	s.store.Logging()

	return nil
}
