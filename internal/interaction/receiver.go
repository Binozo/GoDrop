package interaction

import "github.com/Binozo/GoDrop/pkg/air"

var receiverFoundCallback func(receiver air.Receiver)

func OnReceiverFound(receiver air.Receiver) {
	if receiverFoundCallback != nil {
		receiverFoundCallback(receiver)
	}
}

func AddReceiverFoundCallback(callback func(receiver air.Receiver)) {
	receiverFoundCallback = callback
}
