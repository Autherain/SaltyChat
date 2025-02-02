package store

import (
	"database/sql"

	api "github.com/Autherain/go_cyber"
)

type Store struct {
	db *sql.DB

	Rooms    api.RoomManager
	Messages api.MessageManager
}

type Option func(*Store)

func NewStore(options ...Option) *Store {
	blankStore := &Store{}

	blankStore.Rooms = &roomStore{baseStore: blankStore}
	blankStore.Messages = &messageStore{baseStore: blankStore}

	for _, option := range options {
		option(blankStore)
	}

	if blankStore.db == nil {
		panic("DB store is required")
	}

	return blankStore
}

func WithDB(db *sql.DB) Option {
	return func(s *Store) {
		s.db = db
	}
}
