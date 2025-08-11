//go:build !windows

package listeners

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"

	"github.com/busyster996/dagflow/pkg/sockets"
)

// New creates new listeners for the server.
func New(ctx context.Context, proto, addr string, tlsConfig *tls.Config) (net.Listener, error) {
	switch proto {
	case "tcp":
		ln, err := sockets.NewTCPSocket(ctx, addr, tlsConfig)
		if err != nil {
			return nil, err
		}
		return &Keepalive{Listener: ln}, err
	case "unix":
		return sockets.NewUnixSocket(addr, os.Getegid())
	default:
		return nil, fmt.Errorf("invalid protocol format: %q", proto)
	}
}
