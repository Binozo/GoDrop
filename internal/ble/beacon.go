package ble

import (
	"github.com/Binozo/GoDrop/pkg/account"
	"github.com/Binozo/tinygo-bluetooth"
	"time"
)

const appleBit = 0x004C
const airDropBit = 0x05

// SendAirDropBeacon makes your device discoverable for
// AirDrop sending requests using BLE
func SendAirDropBeacon(appleAccount account.AppleAccount) error {
	manufacturerData := appleAccount.BuildManufacturerData()

	adv := adapter.DefaultAdvertisement()
	adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    device,
		ServiceUUIDs: nil,
		Interval:     bluetooth.NewDuration(time.Millisecond * 200),
		ManufacturerData: map[uint16]interface{}{
			appleBit: manufacturerData,
		},
	})
	return adv.Start()
}

// StartScan searches for AirDrop beacons from other AirDrop senders
func StartScan(onBeaconReceive chan AirDropBeacon) error {
	return adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		// Now we need to check if our result's ManufacturerData contains the holy AppleBit
		for k, v := range result.ManufacturerData() {
			if k == appleBit {
				// Now we know that this beacon has been sent from an apple device

				// Now we need to check if the AirDropBit is present
				if v[0] == airDropBit {
					// we received an AirDrop beacon:)
					beacon := AirDropBeacon{
						DeviceMac:         result.Address.MAC.String(),
						RSSI:              result.RSSI,
						ReceivedTimestamp: time.Now(),
					}

					onBeaconReceive <- beacon
				}
			}
		}
	})
}
