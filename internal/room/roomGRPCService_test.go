package room

import (
	"context"
	"errors"
	"testing"

	"techytechster.com/monkeyparty/pkg"
	"techytechster.com/monkeyparty/pkg/rooms_grpc"
)

func TestCanCreateAndJoinRoom(t *testing.T) {
	initialize()
	service := NewGRPC()
	mockRoomName := "MockRoomName"
	mockRoomPassword := "MockRoomPassword"
	res, err := service.CreateRoom(context.TODO(), &rooms_grpc.CreateRoomRequest{
		Passphrase: &mockRoomPassword,
		RoomName:   &mockRoomName,
	})
	if err != nil {
		t.Fatalf("expected no err for creating room: %s", err.Error())
	}
	_, err = service.JoinRoom(context.TODO(), &rooms_grpc.JoinRoomRequest{
		RoomName:   res.RoomName,
		Passphrase: &mockRoomPassword,
	})
	if err != nil {
		t.Fatalf("expected no err for joining room: %s", err.Error())
	}
}

type MockRoom struct {
	pkg.Roomiface
}

func (m MockRoom) IsStillValid() bool {
	return false
}

func TestBackgroundRoomChecker(t *testing.T) {
	initialize()
	service := &RoomGRPCService{
		rooms: map[string]*pkg.Roomiface{},
	}
	_, err := service.CreateRoom(context.TODO(), &rooms_grpc.CreateRoomRequest{})
	if err != nil {
		t.Fatalf("failed to create room: %s", err.Error())
	}
	var v pkg.Roomiface = MockRoom{}
	service.rooms["asdf"] = &v
	service.backgroundRoomChecker()
	getUnixNow = func() int64 {
		return timeNow().AddDate(0, 0, 1).Unix()
	}
	service.backgroundRoomChecker()
}

func TestCannotJoinNonExistingRooms(t *testing.T) {
	initialize()
	_, err := NewGRPC().JoinRoom(context.TODO(), &rooms_grpc.JoinRoomRequest{
		RoomName: "Yahtzee",
	})
	if err != errRoomDoesNotExist {
		t.Fatalf("expected errRoomDoesNotExist: %s", err.Error())
	}
}

func TestRoomCreationFailure(t *testing.T) {
	initialize()
	cryptoRead = func(b []byte) (n int, err error) {
		return 0, errors.New("mockerr")
	}
	service := NewGRPC()
	mockRoomPassword := "MockRoomPassword"
	_, err := service.CreateRoom(context.TODO(), &rooms_grpc.CreateRoomRequest{
		Passphrase: &mockRoomPassword,
	})
	if err != errFailedToGenerateRoom {
		t.Fatalf("expected errFailedToGenerateRoom: %s", err.Error())
	}
}
