package ble

import "github.com/Binozo/tinygo-bluetooth"

var adapter = bluetooth.DefaultAdapter

const device = "hci0"

// Init the BLE stack using the default hci0 Bluetooth interface
func Init() error {
	return adapter.Enable()
}
