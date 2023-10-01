package room

import (
	"context"
	"testing"

	"techytechster.com/monkeyparty/pkg/rooms_grpc"
)

func TestJoinRoom(t *testing.T) {
	mockRoomName := "MockRoomName"
	mockRoomPassword := "MockRoomPassword"
	NewGRPC().JoinRoom(context.TODO(), &rooms_grpc.JoinRoomRequest{
		RoomName:   mockRoomName,
		Passphrase: &mockRoomPassword,
	})
}
