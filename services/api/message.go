package api

import (
	"time"

	"github.com/Autherain/go_cyber/internal/pagination"
	"github.com/gofrs/uuid"
)

type Message struct {
	ID                uuid.UUID
	RoomID            uuid.UUID
	encrypted_content byte
	nonce             byte
	timestamp         time.Time
}

type MessagesSelector struct {
	*pagination.KeysetSelector[uuid.UUID]

	RoomID uuid.UUID
}

type MessageManager interface {
	ReadMessages(selector *MessagesSelector) (*[]Message, error)
}
