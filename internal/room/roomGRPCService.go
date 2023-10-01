package room

import (
	"context"

	"github.com/rs/zerolog/log"
	"techytechster.com/monkeyparty/pkg/rooms_grpc"
)

type RoomGRPCService struct {
	rooms_grpc.UnimplementedRoomsServer
}

func (r *RoomGRPCService) JoinRoom(c context.Context, payload *rooms_grpc.JoinRoomRequest) (*rooms_grpc.JoinRoomResponse, error) {
	log.Warn().Interface("payload", payload).Msg("JoinRoom is not currently implemented")
	return &rooms_grpc.JoinRoomResponse{}, nil
}

func NewGRPC() rooms_grpc.RoomsServer {
	return &RoomGRPCService{}
}
