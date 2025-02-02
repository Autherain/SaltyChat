package store

import api "github.com/Autherain/go_cyber"

type roomStore struct{ baseStore *Store }

var _ api.RoomManager = (*roomStore)(nil)

func (s *roomStore) ReadRoom(selector *api.RoomSelector) (*api.Room, error) {
	return nil, nil
}

func (s *roomStore) CreateRoom(selector *api.RoomSelector) error {
	return nil
}

func (s *roomStore) DeleteRoom(selector *api.RoomSelector) error {
	return nil
}
