package modext

import (
	api "github.com/Autherain/saltyChat"
	"github.com/Autherain/saltyChat/pkg/store/models"
	"github.com/google/uuid"
)

func MapRoom(room *models.Room) *api.Room {
	roomID, _ := uuid.Parse(room.ID)

	return &api.Room{
		ID:           roomID,
		CreatedAt:    room.CreatedAt,
		LastActivity: room.LastActivity,
		IsActive:     room.IsActive,
	}
}
