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
	room, err := New(payload.RoomName, payload.Passphrase)
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

func (r *RoomGRPCService) backgroundRoomChecker() {
	log.Info().Msg("doing a scan of all current rooms and checking if they are all valid")
	for key, value := range r.rooms {
		log.Info().Interface("room", value)
		if !(*value).IsStillValid() {
			log.Info().Str("room", key).Msg("room is no longer valid, deleting")
			delete(r.rooms, key)
		}
	}
}

func NewGRPC() rooms_grpc.RoomsServer {
	newGrpcServer := &RoomGRPCService{
		rooms: map[string]*pkg.Roomiface{},
	}
	go newGrpcServer.backgroundRoomChecker()
	return newGrpcServer
}

// Theres going to need to be background thread to check rooms are still active
