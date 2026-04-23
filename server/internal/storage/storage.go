package storage

import (
	"fmt"
	"log/slog"
)

type LocaUser struct {
	Username  string
	Groupname string

	UserID  int32
	GroupID int32
}

type Storage struct {
	Users  map[string]LocaUser
	logger *slog.Logger
}

func InitStorage(logger *slog.Logger) *Storage {
	return &Storage{
		Users:  make(map[string]LocaUser),
		logger: logger,
	}
}

func (s *Storage) InsertUser(userID int32, username string, groupID int32, groupname string) error {
	_, ok := s.Users[username]
	if ok {
		return fmt.Errorf("user already linked to another session")
	}

	s.Users[username] = LocaUser{
		Username:  username,
		Groupname: groupname,
		UserID:    userID,
		GroupID:   groupID,
	}

	s.logger.Info("user inserted into storage memory", "name", username)

	return nil
}

func (s *Storage) Logging() {
	s.logger.Info("total users in storage", "count", len(s.Users))
}

func (s *Storage) IsUserOnline(username string) bool {
	_, ok := s.Users[username]
	if ok {
		return true
	}

	return false
}
