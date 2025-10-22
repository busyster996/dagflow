package listeners

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Keepalive struct {
	Listener net.Listener
}

func (kln *Keepalive) Accept() (net.Conn, error) {
	c, err := kln.Listener.Accept()
	if err != nil {
		return nil, err
	}

	kac, err := kln.createKeepaliveConn(c)
	if err != nil {
		return nil, fmt.Errorf("create keepalive connection failed, %w", err)
	}
	// detection time: tcp_keepalive_time + tcp_keepalive_probes + tcp_keepalive_intvl
	// default on linux:  30 + 8 * 30
	// default on osx:    30 + 8 * 75
	if err = kac.SetKeepAlive(true); err != nil {
		return nil, fmt.Errorf("SetKeepAlive failed, %w", err)
	}
	if err = kac.SetKeepAlivePeriod(10 * time.Second); err != nil {
		return nil, fmt.Errorf("SetKeepAlivePeriod failed, %w", err)
	}
	if err = kac.SetReadDeadline(time.Now().Add(120 * time.Second)); err != nil {
		return nil, fmt.Errorf("SetReadDeadline failed, %w", err)
	}
	return kac, nil
}

func (kln *Keepalive) Close() error {
	if kln.Listener == nil {
		return nil
	}
	return kln.Listener.Close()
}

func (kln *Keepalive) Addr() net.Addr {
	if kln.Listener == nil {
		return nil
	}
	return kln.Listener.Addr()
}

func (kln *Keepalive) createKeepaliveConn(c net.Conn) (*keepAliveConn, error) {
	tcpc, ok := c.(*net.TCPConn)
	if !ok {
		return nil, errors.New("only tcp connections have keepalive")
	}
	return &keepAliveConn{tcpc}, nil
}

type keepAliveConn struct {
	*net.TCPConn
}

// SetKeepAlive sets keepalive
func (l *keepAliveConn) SetKeepAlive(doKeepAlive bool) error {
	return l.TCPConn.SetKeepAlive(doKeepAlive)
}

func (l *keepAliveConn) SetKeepAlivePeriod(d time.Duration) error {
	return l.TCPConn.SetKeepAlivePeriod(d)
}
