package store

import (
	"context"
	"database/sql"

	api "github.com/Autherain/saltyChat"
	"github.com/Autherain/saltyChat/internal/utils/errors"
	"github.com/Autherain/saltyChat/internal/utils/pagination"
	"github.com/Autherain/saltyChat/pkg/store/models"
	"github.com/Autherain/saltyChat/pkg/store/modext"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type roomStore struct{ baseStore *Store }

func (s *roomStore) Read(selector *api.RoomSelector) (*api.Room, error) {
	model, err := read(s.baseStore.db, selector)
	if err != nil {
		return nil, err
	}

	return modext.MapRoom(model), err
}

func (s *roomStore) Create(room *api.Room) error {
	if err := room.Validate(); err != nil {
		return err
	}
	return create(s.baseStore.db, room)
}

func (s *roomStore) Delete(selector *api.RoomSelector) error {
	return delete(s.baseStore.db, selector)
}

func read(db *sql.DB, selector *api.RoomSelector) (*models.Room, error) {
	result, err := models.Rooms(models.RoomWhere.ID.EQ(selector.RoomID.String())).One(context.TODO(), db)

	return result, errors.MapSQLError(err)
}

func create(db *sql.DB, room *api.Room) error {
	model := &models.Room{
		ID:           room.ID.String(),
		CreatedAt:    room.CreatedAt,
		LastActivity: room.LastActivity,
		IsActive:     room.IsActive,
	}

	if err := errors.MapSQLError(model.Insert(context.TODO(), db, boil.Infer())); err != nil {
		return err
	}

	return nil
}

func delete(db *sql.DB, selector *api.RoomSelector) error {
	room := &models.Room{ID: selector.RoomID.String()}

	_, err := room.Delete(context.TODO(), db)

	return errors.MapSQLError(err)
}

func (s *roomStore) ReadAll(selector *api.RoomsSelector) ([]*api.Room, uuid.UUID, error) {
	result := []qm.QueryMod{
		qm.OrderBy(models.RoomColumns.ID),
	}

	const limit = 100
	if selector.LastKey != uuid.Nil {
		models.RoomWhere.ID.GT(selector.LastKey.String())
	}

	result = append(result, qm.Limit(int(pagination.NewLimit(selector.Size).Bound(limit))))

	modelsJustRead, err := models.Rooms(result...).All(context.Background(), s.baseStore.db)
	if err != nil {
		return nil, uuid.Nil, err
	}

	lastRoomID, rooms := uuid.Nil, []*api.Room{}
	for _, model := range modelsJustRead {
		roomID, err := uuid.Parse(model.ID)
		if err != nil {
			return nil, uuid.Nil, err
		}
		rooms = append(rooms, &api.Room{
			ID:           roomID,
			CreatedAt:    model.CreatedAt,
			LastActivity: model.LastActivity,
			IsActive:     model.IsActive,
		})
		lastRoomID = roomID
	}

	return rooms, lastRoomID, nil
}
