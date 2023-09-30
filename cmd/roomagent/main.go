package main

import (
	"crypto/rand"
	"encoding/base64"
	"math"

	"github.com/cockroachdb/redact"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"techytechster.com/monkeyparty/internal/room"
	"techytechster.com/monkeyparty/pkg"
)

const SECRET_LENGTH = 4128

func randomBase64String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l] // strip 1 extra character we get from odd length results
}

var rooms []pkg.Roomiface = []pkg.Roomiface{}

func addRoom(password *string) error {
	newRoom, err := room.New(password)
	if err != nil {
		log.Err(err).Msg("failed to create ")
	}
	rooms = append(rooms, newRoom)
	return nil
}

/*
Roomagent handles joining, creating and performing room operations for Monkey party.

This is done through a variety of operations available through roomagent via socket.
*/
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	secretKey := randomBase64String(SECRET_LENGTH)
	secretKeyRedacted := redact.Sprintf("started roomagent, secretkey: %s", secretKey)
	secretKeyRedacted.StripMarkers()
	log.Info().Msg(string(secretKeyRedacted.Redact()))
	var mockPsw = "thisIsAPassword"
	addRoom(&mockPsw)
}
