package owl

import (
	"bufio"
	"io"
	"time"
)

var outLogs []string = make([]string, 0)
var logListeners = make([]chan string, 0)

// SubscribeToLogs allows you read access to OWL logs
func SubscribeToLogs() *chan string {
	receiver := make(chan string)
	logListeners = append(logListeners, receiver)
	return &receiver
}

func UnsubscribeFromLogs(receiver *chan string) {
	newReceiverList := make([]chan string, 0)
	for _, curReceiver := range logListeners {
		if curReceiver != *receiver {
			newReceiverList = append(newReceiverList, curReceiver)
		}
	}
	logListeners = newReceiverList
}

func broadcastLog(logMsg string) {
	for _, receiver := range logListeners {
		select {
		case receiver <- logMsg:
			// Message sent
			break
		case <-time.After(time.Second):
			// Message not sent
			break
		default:
			// Message not sent
			break
		}
	}
}

func listenOnOutput(pipe io.ReadCloser) {
	scanner := bufio.NewScanner(pipe)

	for scanner.Scan() {
		outLogs = append(outLogs, scanner.Text())
		//fmt.Println("Owl Output:", scanner.Text())
		go broadcastLog(scanner.Text())
	}
}
