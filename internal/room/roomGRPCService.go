package room

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"techytechster.com/monkeyparty/pkg"
	"techytechster.com/monkeyparty/pkg/rooms_grpc"
)

type RoomGRPCService struct {
	rooms map[string]*pkg.Roomiface
	rooms_grpc.UnimplementedRoomsServer
}

var errRoomDoesNotExist = errors.New("provided room does not exist")

func (r *RoomGRPCService) JoinRoom(c context.Context, payload *rooms_grpc.JoinRoomRequest) (*rooms_grpc.JoinRoomResponse, error) {
	log.Info().Str("room", payload.RoomName).Msg("Attempting To Join Room")
	if _, ok := r.rooms[payload.RoomName]; !ok {
		log.Info().Msg("Room does not exist")
		return nil, errRoomDoesNotExist
	}
	// TODO: add the consumer to that room
	return &rooms_grpc.JoinRoomResponse{}, nil
}

var errFailedToGenerateRoom = errors.New("failed to generate a room")

func (r *RoomGRPCService) CreateRoom(context context.Context, payload *rooms_grpc.CreateRoomRequest) (*rooms_grpc.CreatedRoomResponse, error) {
	log.Info().Msg("Creating a new room")
	if payload.RoomName != nil { // TODO: rework room.go to allow name
		log.Warn().Msg("Currently not supported: naming of rooms manually - automatic naming done via funName")
	}
	room, err := New(payload.Passphrase)
	if err != nil {
		log.Err(err).Msg("failed to create a room")
		return nil, errFailedToGenerateRoom
	}
	log.Info().Str("room", room.Name()).Msg("Generated a new room")
	r.rooms[room.Name()] = &room
	return &rooms_grpc.CreatedRoomResponse{
		RoomName: room.Name(),
	}, nil
}

func NewGRPC() rooms_grpc.RoomsServer {
	return &RoomGRPCService{
		rooms: map[string]*pkg.Roomiface{},
	}
}

// Theres going to need to be background thread to check rooms are still active
