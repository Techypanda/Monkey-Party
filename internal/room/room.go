package room

import (
	"bytes"
	cryptoRand "crypto/rand"
	"errors"
	"fmt"
	"math/rand"

	"golang.org/x/crypto/argon2"
	"techytechster.com/monkeyparty/pkg"
)

type Room struct {
	Roomname      string
	passphrase    []byte
	encryptionKey []byte
}

func (r Room) Join(password *string) error {
	if r.passphrase != nil {
		return checkPasswordIsCorrect(password, r.passphrase, r.encryptionKey)
	}
	return nil
}

var POSSIBLE_NAMES = []string{"Red", "Green", "Blue", "Yellow", "Rabbit", "Pig", "Tortise", "Yahoo", "Google", "Bing", "Meta", "Amazon"}

func randName() string {
	return POSSIBLE_NAMES[rand.Intn(len(POSSIBLE_NAMES))]
}
func funName() string {
	// TODO: actually generate a 'fun' name like AMONG US
	return fmt.Sprintf("%s-%s-%s-%d", randName(), randName(), randName(), rand.Intn(1024))
}

var cryptoRead = cryptoRand.Read
var errFailedToGenerateBytes = errors.New("failed to generate random bytes via cryptoRead")

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := cryptoRead(b)
	if err != nil {
		return nil, errors.Join(errFailedToGenerateBytes, err)
	}
	return b, nil
}

var errNoPasswordProvided = errors.New("expected a password, none was provided")
var errIncorrectPassword = errors.New("the room password provided is incorrect")
var encryptionMemoryKIB uint32 = 65536
var threads uint8 = 2
var keyLen uint32 = 32
var time uint32 = 3

func checkPasswordIsCorrect(password *string, targetHash []byte, salt []byte) error {
	if password == nil {
		return errNoPasswordProvided
	}
	hash := argon2.IDKey([]byte(*password), salt, time, encryptionMemoryKIB, threads, keyLen)
	if !bytes.Equal(hash, targetHash) {
		return errIncorrectPassword
	}
	return nil
}

var errFailedToGenerateSalt = errors.New("failed to generate salt")

func generateFromPassword(password string) (b64hash []byte, b64salt []byte, err error) {
	salt, err := generateRandomBytes(16)
	if err != nil {
		return nil, nil, errors.Join(errFailedToGenerateSalt, err)
	}
	hash := argon2.IDKey([]byte(password), salt, time, encryptionMemoryKIB, threads, keyLen)
	return hash, salt, nil
}

var errFailedToCreateRoom = errors.New("failed to generate new room")

func New(passphrase *string) (pkg.Roomiface, error) {
	if passphrase != nil {
		hash, salt, err := generateFromPassword(*passphrase)
		if err != nil {
			return nil, errors.Join(errFailedToCreateRoom, err)
		}
		return &Room{
			Roomname:      funName(),
			passphrase:    hash,
			encryptionKey: salt,
		}, nil
	}
	return &Room{
		Roomname: funName(),
	}, nil
}
