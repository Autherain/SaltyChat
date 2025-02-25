package models

import (
	"time"

	"github.com/Autherain/saltyChat"
	timeutil "github.com/Autherain/saltyChat/internal/utils/timeutil"
	"github.com/google/uuid"
)

type RoomParams struct {
	CreatedAt    time.Time `json:"createdAt"`
	LastActivity time.Time `json:"lastActivity"`
	IsActive     bool      `json:"isActive"`
}

func (p *RoomParams) Map() *saltyChat.Room {
	return &saltyChat.Room{
		ID:           uuid.New(),
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
	}
}

// Room represents the JSON response structure
type Room struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"createdAt"`
	LastActivity string `json:"lastActivity"`
	IsActive     bool   `json:"isActive"`
}

// MapRoom converts api.Room to the response structure
func MapRoom(r *saltyChat.Room) *Room {
	return &Room{
		ID:           r.ID.String(),
		CreatedAt:    timeutil.FormatTime(r.CreatedAt),
		LastActivity: timeutil.FormatTime(r.LastActivity),
		IsActive:     r.IsActive,
	}
}
