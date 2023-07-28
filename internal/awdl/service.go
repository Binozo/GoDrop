package awdl

import (
	"context"
	"github.com/Binozo/GoDrop/internal/interaction"
	"github.com/Binozo/GoDrop/pkg/owl"
	"github.com/grandcat/zeroconf"
	"net"
)

const AirDropService = "_airdrop._tcp"
const AirDropPort = 8772

var server *zeroconf.Server
var discoverCancel context.CancelFunc

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

// Discover discovers other AirDrop receiver (mDNS service)
func StartDiscover() error {
	interfaces := []net.Interface{
		{
			Index: owl.GetOwlInterface().Index,
		},
	}
	options := []zeroconf.ClientOption{
		zeroconf.SelectIPTraffic(zeroconf.IPv6),
		zeroconf.SelectIfaces(interfaces),
	}
	resolver, err := zeroconf.NewResolver(options...)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	discoverCancel = cancel
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			if len(entry.AddrIPv6) != 1 || entry.AddrIPv6[0].To16().String() == owl.GetOwlInterfaceAddr() {
				// Not us pls
				continue
			}

			go func(entry *zeroconf.ServiceEntry) {
				ipAddr := entry.AddrIPv6[0]
				port := entry.Port

				receiver, err := sendDiscoverRequest(ipAddr, port)
				if err != nil {
					// device doesn't answer:(
					// This is really weird. I don't know why this happens.
					// The first request always works, all other won't
					// TODO
					// Note: I just found out if you use a real
					// certificate from apple it is more likely to work
					return
				}
				interaction.OnReceiverFound(receiver)
			}(entry)
		}
	}(entries)
	err = resolver.Browse(ctx, AirDropService, "local.", entries)
	if err != nil {
		cancel()
		return err
	}
	return nil
}

// StopDiscover stops discovering AirDrop receiver (mDNS service)
func StopDiscover() {
	if discoverCancel != nil {
		discoverCancel()
	}
}
