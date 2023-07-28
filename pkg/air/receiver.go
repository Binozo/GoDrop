package air

import (
	"net"
	"time"
)

// Receiver defines a receiver to which you can "airdrop"
type Receiver struct {
	ReceiverName string
	IP           net.IP
	Port         int

	// Defines the time of finding this receiver
	// It defines how long a device is "visible"
	Timestamp time.Time
}
