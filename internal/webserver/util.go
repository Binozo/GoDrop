package webserver

import (
	"fmt"
	"github.com/Binozo/GoDrop/pkg/owl"
	"net/http"
	"net/netip"
	"strings"
)

func parseIPv6(r *http.Request) string {
	// RemoteAddr could look like this: [fe80::8c67:8eff:fe7b:def8%awdl0]:50965
	// Now we remove %awdl0
	remoteAddrRaw := strings.ReplaceAll(r.RemoteAddr, fmt.Sprintf("%%%s", owl.OwlInterface), "")
	addr, err := netip.ParseAddrPort(remoteAddrRaw)
	if err != nil {
		return r.RemoteAddr // Give up
	} else {
		return addr.Addr().String() // fe80::8c67:8eff:fe7b:def
	}
}
