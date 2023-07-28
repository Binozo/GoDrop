package air

import "github.com/Binozo/GoDrop/pkg/account"

// Config helps you to adjust GoDrop
type Config struct {
	// Display helpful logs
	DebugMode bool

	// Define the name of your AirDrop receiver device.
	// Will show up on the receiver list of other devices.
	HostName string

	// Specify your target Wifi interface OWL should listen to.
	// Default is `wlan1`. This device MUST support Active Monitor Mode.
	// Take a look at https://github.com/Binozo/GoDrop for more help.
	WifiInterface string

	// You can define which information should be advertised
	// through BLE. You can ignore this. Doesn't change anything.
	// Because we don't use real Apple TLS certificates.
	Account account.AppleAccount
}
