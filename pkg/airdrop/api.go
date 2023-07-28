package airdrop

import (
	"github.com/Binozo/GoDrop/internal/awdl"
	"github.com/Binozo/GoDrop/internal/ble"
	"github.com/Binozo/GoDrop/internal/certificate"
	"github.com/Binozo/GoDrop/internal/interaction"
	"github.com/Binozo/GoDrop/internal/utils"
	"github.com/Binozo/GoDrop/internal/webserver"
	"github.com/Binozo/GoDrop/pkg/account"
	"github.com/Binozo/GoDrop/pkg/air"
	"github.com/Binozo/GoDrop/pkg/owl"
	"github.com/korylprince/go-cpio-odc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const DefaultWifiInterface = "wlan1"

var hostname = "GoDrop"
var wifiInterface = DefaultWifiInterface
var debug = false
var appleAccount = account.Default

func init() {
	osHostname, err := os.Hostname()
	if err == nil {
		hostname = osHostname
	}
}

// Config the current GoDrop setup
func Config(config air.Config) {
	if config.HostName != "" {
		hostname = config.HostName
	}
	debug = config.DebugMode
	emptyAppleAccount := account.AppleAccount{}
	if config.Account != emptyAppleAccount {
		appleAccount = config.Account
	}
	if config.WifiInterface != "" {
		wifiInterface = config.WifiInterface
	}
}

// Init everything required.
// BLE, OWL, WebServer etc.
func Init() error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// We need to check if we are root
	if !utils.IsRoot() {
		if debug {
			log.Error().Msg("Please run this program as root")
		}
		return notRootError
	}

	if debug {
		log.Info().Msgf("Using hostname %s with interface %s", hostname, wifiInterface)
	}

	// Init default responses
	awdl.Init(hostname)

	// Checking for TLS certificate
	certificateExists := certificate.CertificatesPresent()
	if !certificateExists {
		if debug {
			log.Info().Msg("Generating certificate...")
		}
		err := certificate.GenerateCertificates(hostname)
		if err != nil {
			if debug {
				log.Error().Err(err).Msg("Failed to generate certificate")
			}
			return failedCertificateGenerationError
		}
	}

	// Init Bluetooth
	if err := ble.Init(); err != nil {
		if debug {
			log.Error().Err(err).Msg("Couldn't init Bluetooth stack")
		}
		return bleInitFailedError
	}

	// wait until owl is setup and ready
	err := owl.Setup(wifiInterface)
	if err != nil {
		if debug {
			log.Error().Err(err).Msg("Failed to setup OWL")
		}
		return failedOWLSetupError
	}

	if debug {
		log.Info().Msg("GoDrop setup successful")
	}
	go listenForShutdown()
	return nil
}

// StartReceiving starts the listening process for incoming files.
// OnAsk and OnFiles need to be called before StartReceiving
func StartReceiving() error {
	airDropBeaconReceiver := make(chan ble.AirDropBeacon)

	go func() {
		for {
			receivedBeacon := <-airDropBeaconReceiver
			if debug {
				log.Debug().Msgf("Received AirDrop Beacon from %s with RSSI %d", receivedBeacon.DeviceMac, receivedBeacon.RSSI)
			}
			err := awdl.RegisterService(hostname)
			if err != nil {
				// Failed to register service somehow
				if debug {
					log.Error().Err(err).Msg("Couldn't register AirDrop Service")
				}
			}
		}
	}()

	// Starting BLE scan
	go func() {
		for {
			err := ble.StartScan(airDropBeaconReceiver)
			if err != nil {
				restartTimeout := time.Second * 2
				if debug {
					log.Error().Err(err).Msgf("BLE scan failed. Restarting in %d seconds...", restartTimeout.Seconds())
				}
				time.Sleep(restartTimeout)
			}
		}
	}()

	// Now we start Bluetooth Advertisement
	err := ble.SendAirDropBeacon(appleAccount)
	if err != nil {
		if debug {
			log.Error().Err(err).Msg("Failed to start BLE advertising")
		}
		return failedStartingBLEAdvertising
	}

	err = webserver.Start()
	if err != nil {
		if debug {
			log.Error().Err(err).Msg("Failed to start webserver")
		}
		return failedStartingWebServer
	}

	return nil
}

// OnAsk gives you the option to accept or decline requests.
// If you don't define a callback, GoDrop will always decline requests
func OnAsk(callback func(request air.Request) bool) {
	interaction.AddAskCallback(callback)
}

// OnFiles will be called if GoDrop received ("AirDropped") files.
// air.Request shows you information about the sender.
// []*cpio.File contains the actual files.
func OnFiles(callback func(request air.Request, files []*cpio.File)) {
	interaction.AddFilesCallback(callback)
}

// Shutdown gracefully
func listenForShutdown() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM) // all OS signals
	sig := <-ch
	log.Info().Msgf("[GoDrop] Shutting down... (%v)", sig)
	if err := owl.Kill(); err != nil {
		log.Error().Err(err).Msg("Couldn't kill Owl")
	}

	if err := ble.Shutdown(); err != nil {
		log.Error().Err(err).Msg("Couldn't BLE scan")
	}

	awdl.ShutdownService()

	//os.Exit(-1)
}
