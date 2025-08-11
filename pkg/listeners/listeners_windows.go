//go:build windows

package listeners

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"

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
	case "npipe":
		// allow Administrators and SYSTEM, plus whatever additional users or groups were specified
		sddl := "D:P(A;;GA;;;BA)(A;;GA;;;SY)"
		c := winio.PipeConfig{
			SecurityDescriptor: sddl,
			MessageMode:        true,  // Use message mode so that CloseWrite() is supported
			InputBufferSize:    65536, // Use 64KB buffers to improve performance
			OutputBufferSize:   65536,
		}
		return winio.ListenPipe(addr, &c)
	default:
		return nil, fmt.Errorf("invalid protocol format: windows only supports tcp and npipe")
	}
}
