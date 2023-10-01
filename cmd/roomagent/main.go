package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/cockroachdb/redact"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"techytechster.com/monkeyparty/internal/room"
	"techytechster.com/monkeyparty/internal/utils"
	"techytechster.com/monkeyparty/pkg/rooms_grpc"
)

const SECRET_LENGTH = 4128

const DEFAULT_ROOM_PORT uint = 8054

var flagParse = flag.Parse
var flagUint = flag.Uint
var netListen = net.Listen
var fatalLog = log.Fatal()
var newGRPCServer = func() *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	return grpcServer
}

/*
Roomagent handles joining, creating and performing room operations for Monkey party.

This is done through a variety of operations available through roomagent via socket.

Flags:
-port {Port_Number} uint: The port to run the grpc service on (Default - DEFAULT_ROOM_PORT=8054)

Usage:
-h: display help
*/
func main() {
	port := flagUint("port", DEFAULT_ROOM_PORT, "set the port to run grpc service on")
	flagParse()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// Start The Server
	lis, err := netListen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		fatalLog.Err(err).Uint("port", *port).Msg("failed to listen on port") // Return is only needed for unit tests
		return
	}
	// TODO: Probably swap for a certificate or two part key (API holds one key RoomAgent holds second)
	secretKey := utils.RandomBase64String(SECRET_LENGTH)
	secretKeyRedacted := redact.Sprintf("generated secretkey: %s", secretKey)
	secretKeyRedacted.StripMarkers()
	log.Info().Msg(string(secretKeyRedacted.Redact()))
	grpcServer := newGRPCServer()
	rooms_grpc.RegisterRoomsServer(grpcServer, room.NewGRPC())
	log.Info().Uint("port", *port).Msg("Started GRPC Service")
	reflection.Register(grpcServer)
	if err = grpcServer.Serve(lis); err != nil {
		fatalLog.Err(err).Uint("port", *port).Msg("failed to serve grpc service") // Return is only needed for unit tests
		return
	}
}
