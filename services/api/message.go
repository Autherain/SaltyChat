package saltyChat

import (
	"time"

	"github.com/Autherain/saltyChat/internal/utils/pagination"
	"github.com/gofrs/uuid"
)

type Message struct {
	ID               uuid.UUID
	RoomID           uuid.UUID
	EncryptedContent byte
	Nonce            byte
	Timestamp        time.Time
}

type MessagesSelector struct {
	*pagination.KeysetSelector[uuid.UUID]

	RoomID uuid.UUID
}

type MessageManager interface {
	ReadMessages(selector *MessagesSelector) (*[]Message, error)
}
