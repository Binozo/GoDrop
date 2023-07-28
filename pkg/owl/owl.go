package owl

import (
	"github.com/rs/zerolog/log"
	"net"
	"os/exec"
	"strings"
)

const OwlInterface = "awdl0"

var owlCmd *exec.Cmd

// Setup the Wifi interface and start OWL
func Setup(wifiInterface string) error {
	setupWifi(wifiInterface)
	return startOwl(wifiInterface)
}

func startOwl(wifiInterface string) error {
	owlCmd = exec.Command("owl", "-i", wifiInterface, "-c", "44") // -c 44 works best

	// Owl is somehow only using stderr as main log output
	//stdout, _ := owlCmd.StdoutPipe()
	stderr, _ := owlCmd.StderrPipe()
	//go listenOnOutput(stdout)
	go listenOnOutput(stderr)

	err := owlCmd.Start()
	if err != nil {
		return err
	}

	// Now we wait until owl is ready
	WaitUntilOwlReady()
	//log.Info().Msgf("Owl is up and running at %s", GetOwlInterfaceAddr())

	go func() {
		err = owlCmd.Wait()
		if err != nil {
			// Exited with error
			owlTerminated(outLogs)
		}
	}()
	return nil
}

func GetOwlInterface() net.Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't fetch interfaces")
	}
	for _, curInterface := range interfaces {
		if curInterface.Name == OwlInterface {
			return curInterface
		}
	}
	// TODO: add failure recovery -> restarting owl
	log.Fatal().Msg("Couldn't find Owl interface")
	return net.Interface{}
}

func GetOwlInterfaceAddr() string {
	owlInterface := GetOwlInterface()
	addrs, err := owlInterface.Addrs()
	if err != nil {
		log.Fatal().Msg("Couldn't find Owl interface address")
	}
	for _, curAddrs := range addrs {
		// https://gist.github.com/schwarzeni/f25031a3123f895ff3785970921e962c
		return curAddrs.(*net.IPNet).IP.String() // this is a better solution

		rawAddr := curAddrs.String() // could look like this: fe80::20f:c9ff:fe13:1106/64

		// we need to remove the "/64"
		return strings.Split(rawAddr, "/")[0]
	}

	log.Fatal().Msg("Couldn't find Owl interface")
	return ""
}

// Kill Owl
func Kill() error {
	return owlCmd.Process.Kill()
}
