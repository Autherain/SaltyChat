package server

import (
	"encoding/json"
	"testing"

	"github.com/Autherain/saltyChat"
	"github.com/Autherain/saltyChat/pkg/server/models"
	"github.com/Autherain/saltyChat/pkg/store"
	"github.com/google/uuid"
	"github.com/jirenius/go-res"
	"github.com/jirenius/go-res/restest"
	"github.com/stretchr/testify/require"
)

func TestHandleCreateRoom(t *testing.T) {
	roomID := uuid.Nil

	params := &models.RoomParams{}
	encodedParams, err := json.Marshal(params)
	require.NoError(t, err)

	server := New(
		WithService(res.NewService("test")),
		WithStore(&store.Store{
			Rooms: &roomsMock{
				createFunc: func(room *saltyChat.Room) error {
					roomID = room.ID
					return nil
				},
			},
		}),
	)

	session := newTestSession(t, server.Service)
	defer session.Close()

	session.Call("test.rooms", "new", &restest.Request{Params: encodedParams})
	session.GetMsg().AssertResource("test.rooms." + roomID.String())
	session.GetMsg().AssertQueryEvent("test.rooms", nil)
}

type roomsMock struct {
	createFunc  func(room *saltyChat.Room) error
	readAllFunc func(selector *saltyChat.RoomsSelector) ([]*saltyChat.Room, uuid.UUID, error)

	// ReadAll(selector *RoomsSelector) ([]*Room, uuid.UUID, error)
}

func (m *roomsMock) Create(selector *saltyChat.Room) error {
	return m.createFunc(selector)
}

func (m *roomsMock) ReadAll(selector *saltyChat.RoomsSelector) ([]*saltyChat.Room, uuid.UUID, error) {
	return m.readAllFunc(selector)
}

func (m *roomsMock) Delete(selector *saltyChat.RoomSelector) error {
	return nil
}

func (m *roomsMock) Read(selector *saltyChat.RoomSelector) (*saltyChat.Room, error) {
	return nil, nil
}
