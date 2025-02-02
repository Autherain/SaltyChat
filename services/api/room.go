package api

import (
	"time"

	"github.com/gofrs/uuid"
)

type Room struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	LastActivity time.Time
	is_active    bool
}

type RoomSelector struct {
	RoomID uuid.UUID
}

type RoomManager interface {
	ReadRoom(selector *RoomSelector) (*Room, error)
	CreateRoom(selector *RoomSelector) error
	DeleteRoom(selector *RoomSelector) error
}
