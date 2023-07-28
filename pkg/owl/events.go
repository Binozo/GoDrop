package owl

import (
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

var errorListeners = make([]chan OwlError, 0)

// SubscribeToErrors Subscribe to any error events like Owl exited with non-zero code
func SubscribeToErrors() chan OwlError {
	receiver := make(chan OwlError)
	errorListeners = append(errorListeners, receiver)
	return receiver
}

// WaitUntilOwlReady Waits until Owl is ready to access the network
func WaitUntilOwlReady() {
	var receiver *chan string = SubscribeToLogs()
	for true {
		logMsg := <-*receiver
		if strings.Contains(logMsg, "Host device: awdl") {
			// Owl is ready
			// We have to wait a bit until the interface is ready for us
			time.Sleep(time.Second * 2) // 2 seconds is ideal
			break
		}
	}
	UnsubscribeFromLogs(receiver)
}

// Owl unexpectedly terminated with an error code non-zero
// We need to broadcast this information
func owlTerminated(errorLogs []string) {

	// we need to investigate which kind of error occurred
	errorType := getExitReason(errorLogs)

	owlError := OwlError{
		ErrorType: errorType,
		Logs:      errorLogs,
	}
	for _, receiver := range errorListeners {
		receiver <- owlError
	}

	// In case nobody is listening we need to print to the console
	if len(errorListeners) == 0 {
		log.Error().Msg("Owl sadly exited with non-zero code. Can't continue.")
		if errorType == Unknown {
			log.Error().Msgf("I can't tell you what exactly went wrong. Maybe you should take a look: %s", errorLogs)
		} else {
			log.Error().Msgf("The following Owl error occurred: %s.", errorType)
		}
	}
}

func getExitReason(errorLogs []string) OwlErrorType {
	for _, curLine := range errorLogs {

		if strings.Contains(curLine, "No such interface exists") {
			return UnknownInterface
		}

	}

	// couldn't determine error type
	log.Error().Msgf("Couldn't determine error type from logs: %s", errorLogs)
	return Unknown
}
