package interaction

import (
	"github.com/Binozo/GoDrop/pkg/air"
	"github.com/korylprince/go-cpio-odc"
	"github.com/rs/zerolog/log"
)

var confirmationCallback func(request air.Request) bool
var filesCallback func(request air.Request, files []*cpio.File)
var cachedRequests = make([]air.Request, 0)

func AddAskCallback(callback func(request air.Request) bool) {
	confirmationCallback = callback
}

func AskUserAboutAcceptingIncomingFiles(request air.Request) bool {
	if confirmationCallback != nil {
		requestConfirmed := confirmationCallback(request)
		if requestConfirmed {
			// Now we cache the request in order to
			// pass it with the sent files in the OnFiles callback

			cachedRequests = append(cachedRequests, request)
		}
		return requestConfirmed
	}
	return false
}

func AddFilesCallback(callback func(request air.Request, files []*cpio.File)) {
	filesCallback = callback
}

func OnFiles(files []*cpio.File, senderIP string) {
	if filesCallback != nil {
		// Now we get the cached air.Request
		var request air.Request
		for _, cachedRequest := range cachedRequests {
			if cachedRequest.SenderIP == senderIP {
				// Found the cached request!
				request = cachedRequest
				break
			}
		}

		if request.SenderID == "" {
			// no cached request found. Sender somehow sent without asking
			log.Warn().Msg("Declining uploaded files because no /Ask request has been sent")
		} else {
			filesCallback(request, files)
		}
	} else {
		log.Warn().Msg("AirDrop files received but no OnFiles handler has been registered")
	}
}
