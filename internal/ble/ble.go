package ble

import "github.com/Binozo/tinygo-bluetooth"

var adapter = bluetooth.DefaultAdapter
var adv *bluetooth.Advertisement

const device = "hci0"

// Init the BLE stack using the default hci0 Bluetooth interface
func Init() error {
	return adapter.Enable()
}
