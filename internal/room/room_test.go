package room

import (
	cryptoRand "crypto/rand"
	"errors"
	"strings"
	"testing"
)

func initialize() {
	cryptoRead = cryptoRand.Read
}

func TestCanJoinRoom(t *testing.T) {
	initialize()
	r, err := New(nil)
	if err != nil {
		t.Fatalf("expected no err for creating room: %s", err.Error())
	}
	inputPassword := "something"
	if err = r.Join(&inputPassword); err != nil {
		t.Fatalf("expected no err for joining room: %s", err.Error())
	}
	mockPassword := "ThisIsASecurePassword"
	r, err = New(&mockPassword)
	if err != nil {
		t.Fatalf("expected no err for creating room with password: %s", err.Error())
	}
	if err := r.Join(nil); err != errNoPasswordProvided {
		t.Fatalf("expected a error for joining room with no password: %s", err.Error())
	}
	if err = r.Join(&inputPassword); err != errIncorrectPassword {
		t.Fatalf("expected a error if i provide incorrect password: %s", err.Error())
	}
	inputPassword = strings.ToUpper(mockPassword)
	if err = r.Join(&inputPassword); err != errIncorrectPassword {
		t.Fatalf("expected a error if i provide wrong casing: %s", err.Error())
	}
	inputPassword = "ThisIsASecurePassword"
	if err = r.Join(&inputPassword); err != nil {
		t.Fatalf("expected no error as i provided correct password: %s", err.Error())
	}
}

func TestCryptoFailures(t *testing.T) {
	initialize()
	cryptoRead = func(b []byte) (n int, err error) {
		return 0, errors.New("mock error")
	}
	mockPassword := "ThisIsASecurePassword"
	_, err := New(&mockPassword)
	if !errors.Is(err, errFailedToCreateRoom) {
		t.Fatalf("expected err to be errFailedToCreateRoom: %s", err.Error())
	}
}
