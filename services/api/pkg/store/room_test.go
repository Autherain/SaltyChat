package store

import (
	"testing"
	"time"

	api "github.com/Autherain/saltyChat"
	"github.com/Autherain/saltyChat/environment"
	"github.com/Autherain/saltyChat/internal/utils/errors"
	"github.com/Autherain/saltyChat/internal/utils/pagination"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

//nolint:revive //Test don't care
func TestRoomStore(t *testing.T) {
	t.Parallel()
	store := NewStore(WithDB(environment.MustInitPGSQLDB(environment.Parse()))).Rooms

	now := time.Now().UTC()
	room := &api.Room{
		ID:           uuid.New(),
		CreatedAt:    now,
		LastActivity: now,
		IsActive:     true,
	}

	//nolint:paralleltest // Need series of test
	t.Run("create room", func(t *testing.T) {
		require.NoError(t, store.Create(room))
	})

	//nolint:paralleltest // Need series
	t.Run("read room", func(t *testing.T) {
		selector := &api.RoomSelector{RoomID: room.ID}
		fetchedRoom, err := store.Read(selector)
		require.NoError(t, err)
		require.Equal(t, room.ID, fetchedRoom.ID)
		require.Equal(t, room.IsActive, fetchedRoom.IsActive)
		require.WithinDuration(t, room.CreatedAt, fetchedRoom.CreatedAt, time.Second)
		require.WithinDuration(t, room.LastActivity, fetchedRoom.LastActivity, time.Second)
	})

	//nolint:paralleltest // Need series
	t.Run("read all rooms", func(t *testing.T) {
		// Create a few more rooms for testing ReadAll
		additionalRooms := []*api.Room{
			{
				ID:           uuid.New(),
				CreatedAt:    now.Add(time.Minute),
				LastActivity: now.Add(time.Minute),
				IsActive:     true,
			},
			{
				ID:           uuid.New(),
				CreatedAt:    now.Add(2 * time.Minute),
				LastActivity: now.Add(2 * time.Minute),
				IsActive:     false,
			},
		}

		for _, r := range additionalRooms {
			require.NoError(t, store.Create(r))
		}

		// Test reading all rooms with basic selector
		selector := &api.RoomsSelector{
			KeysetSelector: &pagination.KeysetSelector[uuid.UUID]{
				Size: 100, // Large enough to get all rooms
			},
		}

		allRooms, _, err := store.ReadAll(selector)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(allRooms), 3) // Should have at least our 3 rooms

		// Verify we can find all our created rooms in the results
		createdIDs := map[string]bool{
			room.ID.String():               false,
			additionalRooms[0].ID.String(): false,
			additionalRooms[1].ID.String(): false,
		}

		for _, fetchedRoom := range allRooms {
			if _, exists := createdIDs[fetchedRoom.ID.String()]; exists {
				createdIDs[fetchedRoom.ID.String()] = true
			}
		}

		// Verify all rooms were found
		for id, found := range createdIDs {
			require.True(t, found, "Room with ID %s was not found in ReadAll results", id)
		}
	})

	t.Run("room not found", func(t *testing.T) {
		t.Parallel()
		nonExistentID := uuid.New()
		selector := &api.RoomSelector{RoomID: nonExistentID}
		_, err := store.Read(selector)
		require.Error(t, err)
		require.Equal(t, errors.CodeNotFound, errors.ErrorCode(err))
	})

	t.Run("validate room", func(t *testing.T) {
		t.Parallel()
		invalidRoom := &api.Room{
			// Missing ID
			CreatedAt:    now,
			LastActivity: now,
			IsActive:     true,
		}
		err := store.Create(invalidRoom)
		require.Error(t, err)
		require.Equal(t, errors.CodeInvalid, errors.ErrorCode(err))
	})

	t.Run("delete room", func(t *testing.T) {
		t.Parallel()
		selector := &api.RoomSelector{RoomID: room.ID}
		require.NoError(t, store.Delete(selector))

		// Verify room is deleted
		_, err := store.Read(selector)
		require.Error(t, err)
		require.Equal(t, errors.CodeNotFound, errors.ErrorCode(err))
	})
}
