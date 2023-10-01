package room

import (
	"bytes"
	cryptoRand "crypto/rand"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/argon2"
	"techytechster.com/monkeyparty/pkg"
)

type Room struct {
	Roomname      string
	nextCheck     int64
	passphrase    []byte
	encryptionKey []byte
}

var timeNow = time.Now
var getUnixNow = func() int64 {
	return timeNow().Unix()
}

func (r *Room) roomContainsActivePlayers() bool {
	r.nextCheck = timeNow().Add(time.Minute * time.Duration(timeToCheckMinutesPeriod)).Unix()
	return true
}
func (r *Room) IsStillValid() bool {
	now := getUnixNow()
	if r.nextCheck < now {
		return r.roomContainsActivePlayers()
	}
	return true
}
func (r *Room) Join(password *string) error {
	if r.passphrase != nil {
		return checkPasswordIsCorrect(password, r.passphrase, r.encryptionKey)
	}
	return nil
}
func (r *Room) Name() string {
	return r.Roomname
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

const threads uint8 = 2
const keyLen uint32 = 32
const timeConstant uint32 = 3

func checkPasswordIsCorrect(password *string, targetHash []byte, salt []byte) error {
	if password == nil {
		return errNoPasswordProvided
	}
	hash := argon2.IDKey([]byte(*password), salt, timeConstant, encryptionMemoryKIB, threads, keyLen)
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
	hash := argon2.IDKey([]byte(password), salt, timeConstant, encryptionMemoryKIB, threads, keyLen)
	return hash, salt, nil
}

var errFailedToCreateRoom = errors.New("failed to generate new room")

const timeToCheckMinutesPeriod = 2

func New(roomName *string, passphrase *string) (pkg.Roomiface, error) {
	name := roomName
	if roomName == nil {
		genName := funName()
		name = &genName
	}
	if passphrase != nil {
		hash, salt, err := generateFromPassword(*passphrase)
		if err != nil {
			return nil, errors.Join(errFailedToCreateRoom, err)
		}
		return &Room{
			Roomname:      *name,
			passphrase:    hash,
			encryptionKey: salt,
			nextCheck:     timeNow().Add(time.Minute * time.Duration(timeToCheckMinutesPeriod)).Unix(),
		}, nil
	}
	return &Room{
		Roomname:  *name,
		nextCheck: timeNow().Add(time.Minute * time.Duration(timeToCheckMinutesPeriod)).Unix(),
	}, nil
}
