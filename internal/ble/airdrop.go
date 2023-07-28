package ble

import "time"

// AirDropBeacon contains information about
// the BLE advertising from a sender.
// It shows that some is open for receiving requests.
type AirDropBeacon struct {
	DeviceMac         string
	RSSI              int16
	ReceivedTimestamp time.Time
}
