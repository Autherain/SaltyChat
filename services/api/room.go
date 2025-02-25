package saltyChat

import (
	"time"

	"github.com/Autherain/saltyChat/internal/utils/errors"
	"github.com/Autherain/saltyChat/internal/utils/pagination"
	vld "github.com/tiendc/go-validator"

	"github.com/google/uuid"
)

type Room struct {
	ID           uuid.UUID
	CreatedAt    time.Time
	LastActivity time.Time
	IsActive     bool
}

func (r *Room) Validate() error {
	errs := vld.Validate(
		// Validation de l'ID
		vld.Required(&r.ID).OnError(
			vld.SetField("id", nil),
			vld.SetCustomKey("ERR_ROOM_ID_REQUIRED"),
		),
	)

	if len(errs) > 0 {
		detail, _ := errs[0].BuildDetail()
		return &errors.Error{
			Code:      errors.CodeInvalid,
			Message:   detail,
			Operation: "Room.Validate",
		}
	}

	return nil
}

type RoomSelector struct {
	RoomID uuid.UUID
}

type RoomsSelector struct {
	*pagination.KeysetSelector[uuid.UUID]
}

type RoomManager interface {
	ReadAll(selector *RoomsSelector) ([]*Room, uuid.UUID, error)
	Create(selector *Room) error
	Delete(selector *RoomSelector) error
	Read(selector *RoomSelector) (*Room, error)
}
