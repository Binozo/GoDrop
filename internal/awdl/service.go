package awdl

import (
	"github.com/Binozo/GoDrop/pkg/owl"
	"github.com/grandcat/zeroconf"
	"net"
)

const AirDropService = "_airdrop._tcp"
const AirDropPort = 8772

var server *zeroconf.Server

// RegisterService registers the required mDNS service
// through the network provided by OWL
// Without this no AirDrop sender would be able to "see" you
func RegisterService(hostname string) error {
	if server != nil {
		// here we are avoiding the duplicate entries of this device
		server.TTL(120)
		server.SetText([]string{"flags=136"})
		return nil
	}
	addr := owl.GetOwlInterfaceAddr()

	serviceID := getRandomServiceID()
	var err error
	server, err = zeroconf.RegisterProxy(serviceID, AirDropService, "local.", AirDropPort, hostname, []string{addr}, []string{"flags=136"}, []net.Interface{owl.GetOwlInterface()})
	if err != nil {
		return err
	}
	server.TTL(120)
	return nil
}

func ShutdownService() {
	if server != nil {
		server.Shutdown()
	}
}
