package server

import (
	"github.com/Autherain/saltyChat"
	"github.com/Autherain/saltyChat/internal/utils/errors"
	"github.com/Autherain/saltyChat/internal/utils/pagination"
	"github.com/Autherain/saltyChat/internal/utils/resutil"
	"github.com/Autherain/saltyChat/pkg/server/models"
	"github.com/google/uuid"
	"github.com/jirenius/go-res"
)

const (
	RessourceRooms = "rooms"
	RoomIDPram     = "$roomID"
)

func (s *Server) registerRoomRoutes() {
	s.Service.Handle(RessourceRooms,
		res.Access(res.AccessGranted),
		s.handleCreateRoom(),
	)
}

func (s *Server) handleCreateRoom() res.Option {
	return res.Call("new", func(request res.CallRequest) {
		params := &models.RoomParams{}
		request.ParseParams(params)

		room := params.Map()
		if err := s.store.Rooms.Create(room); err != nil {
			errors.LogAndWriteRESError(s.log, request, err)
		}

		request.Resource(
			s.getRoomSelectorRID(s.Service, &saltyChat.RoomSelector{
				RoomID: room.ID,
			}),
		)

		if err := s.sendRoomsQueryEvent(); err != nil {
			errors.LogAndWriteRESError(s.log, request, err)
		}
	})
}

func (s *Server) handleReadRoom() res.Option {
	return res.GetModel(func(request res.ModelRequest) {
	})
}

func (s *Server) getRoomSelectorRID(service *res.Service, selector *saltyChat.RoomSelector) string {
	return resutil.JoinResourcePath(service.FullPath(), RessourceRooms, selector.RoomID.String())
}

func (s *Server) sendRoomsQueryEvent() error {
	return resutil.HandleCollectionQueryRequest(
		s.Service,
		resutil.JoinResourcePath(s.Service.FullPath(), RessourceRooms),
		func(request res.QueryRequest) ([]res.Ref, error) {
			return s.getRoomReferencesFromRessource(&roomsResource{request})
		},
	)
}

func (s *Server) getRoomReferencesFromRessource(resource *roomsResource) ([]res.Ref, error) {
	selector, err := resource.parseSelector()
	if err != nil {
		return nil, err
	}

	rooms, _, err := s.store.Rooms.ReadAll(selector)
	if err != nil {
		return nil, err
	}

	result := []res.Ref{}
	for _, room := range rooms {
		result = append(result, res.Ref(
			s.getRoomSelectorRID(s.Service, &saltyChat.RoomSelector{
				RoomID: room.ID,
			}),
		))
	}

	return result, nil
}

type roomsResource struct{ res.Resource }

func (r *roomsResource) parseSelector() (*saltyChat.RoomsSelector, error) {
	query := r.ParseQuery()

	result := &saltyChat.RoomsSelector{}

	keysetSelector, err := pagination.ParseKeysetSelector(query, func(key string) (uuid.UUID, error) {
		parsedKey, err := uuid.Parse(key)
		if err != nil {
			return uuid.Nil, err
		}
		return parsedKey, nil
	})
	if err != nil {
		return nil, err
	}
	result.KeysetSelector = keysetSelector

	return result, nil
}
