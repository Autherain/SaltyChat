package store

import api "github.com/Autherain/go_cyber"

type messageStore struct{ baseStore *Store }

var _ api.MessageManager = (*messageStore)(nil)

func (s *messageStore) ReadMessages(selector *api.MessagesSelector) (*[]api.Message, error) {
	return nil, nil
}
