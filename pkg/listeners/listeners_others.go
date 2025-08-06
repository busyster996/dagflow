//go:build !windows

package listeners

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"

	"github.com/busyster996/dagflow/pkg/sockets"
)

// Init creates new listeners for the server.
func Init(proto, addr string, tlsConfig *tls.Config) (net.Listener, error) {
	switch proto {
	case "tcp":
		return sockets.NewTCPSocket(addr, tlsConfig)
	case "unix":
		return sockets.NewUnixSocket(addr, os.Getegid())
	default:
		return nil, fmt.Errorf("invalid protocol format: %q", proto)
	}
}
