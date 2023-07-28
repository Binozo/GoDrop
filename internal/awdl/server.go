package awdl

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"howett.net/plist"
)

var DiscoverResponse []byte
var AskAcceptResponse []byte

// in case the user forgets:)
func init() {
	Init("GoDrop")
}

// Init fixed responses for the webserver
// Nothing wild here
func Init(hostname string) {
	type discoverData struct {
		ReceiverComputerName      string
		ReceiverModelName         string
		ReceiverMediaCapabilities []byte
		ReceiverRecordData        []byte
	}

	var discoverPlist bytes.Buffer
	plistEncoder := plist.NewBinaryEncoder(&discoverPlist)

	err := plistEncoder.Encode(discoverData{
		ReceiverComputerName:      hostname,
		ReceiverModelName:         "MacbookPro5.1",
		ReceiverMediaCapabilities: []byte("{\"Version\":1}"),
		ReceiverRecordData:        nil, // we don't need that lol
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't encode Discovery Plist")
	}

	DiscoverResponse = discoverPlist.Bytes()

	type askResponse struct {
		ReceiverModelName    string
		ReceiverComputerName string
	}

	var askResponsePlist bytes.Buffer
	askResponseEncoder := plist.NewBinaryEncoder(&askResponsePlist)

	if err := askResponseEncoder.Encode(askResponse{
		ReceiverModelName:    hostname,
		ReceiverComputerName: "MacbookPro5.1",
	}); err != nil {
		log.Fatal().Err(err).Msg("Couldn't encode Ask Response Plist")
	}

	AskAcceptResponse = askResponsePlist.Bytes()
}

func getRandomServiceID() string {
	rBytes := make([]byte, 6)
	if _, err := rand.Read(rBytes); err != nil {
		return ""
	}
	return hex.EncodeToString(rBytes)
}
