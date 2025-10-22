// Package sockets provides helper functions to create and configure Unix or TCP sockets.
package sockets

import (
	"context"
	"crypto/tls"
	"net"
)

// NewTCPSocket creates a TCP socket listener with the specified address and
// the specified tls configuration. If TLSConfig is set, will encapsulate the
// TCP listener inside a TLS one.
func NewTCPSocket(ctx context.Context, addr string, tlsConfig *tls.Config) (net.Listener, error) {
	lc := &net.ListenConfig{}
	l, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		tlsConfig.NextProtos = []string{"http/1.1"}
		l = tls.NewListener(l, tlsConfig)
	}
	return l, nil
}
