package storage

import (
	"fmt"
	"log/slog"
	"sync"
)

type LocalUser struct {
	Username string
	UserID   int32
	IsOnline bool
}

// RWMutex in it
type Storage struct {
	sync.RWMutex
	Users  map[string]*LocalUser
	logger *slog.Logger
}

func InitStorage(logger *slog.Logger) *Storage {
	return &Storage{
		Users:  make(map[string]*LocalUser),
		logger: logger,
	}
}

func (s *Storage) InsertUser(userID int32, username string) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.Users[username]
	if ok {
		return fmt.Errorf("user already found in the storage")
	}

	s.Users[username] = &LocalUser{
		Username: username,
		UserID:   userID,
		IsOnline: false,
	}

	s.logger.Info("user inserted into storage memory", "name", username)
	return nil
}

func (s *Storage) Logging() {
	s.RLock()
	s.logger.Info("total users in storage", "count", len(s.Users))
	s.RUnlock()
}

func (s *Storage) IsUserOnline(username string) bool {
	s.RLock()
	_, ok := s.Users[username]
	s.RUnlock()
	if ok {
		return true
	}

	return false
}
