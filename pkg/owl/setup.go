package owl

import (
	"github.com/rs/zerolog/log"
	"os/exec"
)

func setupWifi(wifiInterface string) {
	// We need to remove any programs that block our wifi interface
	cmd := exec.Command("ip", "link", "set", "dev", wifiInterface, "down")
	err := cmd.Start()

	if err != nil {
		// Setting up our Wifi Interface was not successful
		log.Fatal().Err(err).Msgf("Unable to setup Wifi Interface %s", wifiInterface)
	}
}
