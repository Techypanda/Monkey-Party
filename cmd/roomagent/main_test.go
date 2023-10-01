package main

import (
	"errors"
	"net"
	"testing"

	"github.com/rs/zerolog/log"
)

var mockNetListen = func(network, address string) (net.Listener, error) {
	return &net.TCPListener{}, nil
}

type MockListener struct {
}

func (m MockListener) Accept() (net.Conn, error) {
	return nil, errors.New("mock accept")
}
func (m MockListener) Close() error {
	return errors.New("mock close")
}
func (m MockListener) Addr() net.Addr {
	return &net.IPAddr{}
}

func setup() {
	flagParse = func() {}
	flagUint = func(name string, value uint, usage string) *uint {
		var mockFlagPort = DEFAULT_ROOM_PORT
		return &mockFlagPort
	}
	netListen = mockNetListen
	fatalLog = log.Info() // TODO find a better way to test: Fatal will cause a os.exit call which is hard to unit test
}

func TestRoomAgentTestNetListenerFailure(t *testing.T) {
	setup()
	netListen = func(network, address string) (net.Listener, error) {
		return nil, errors.New("mock error")
	}
	// TODO: should spy on stdout
	main()
}

func TestServeError(t *testing.T) {
	setup()
	netListen = func(network, address string) (net.Listener, error) {
		return MockListener{}, nil
	}
	main()
}
