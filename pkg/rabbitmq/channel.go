package rabbitmq

import (
	"errors"
	"fmt"

	"github.com/busyster996/dagflow/pkg/rabbitmq/internal/manager/channel"
)

type Channel struct {
	*channel.Manager
}

// NewChannel returns a new channel to the cluster.
func NewChannel(conn *Conn, optionFuncs ...func(*ChannelOptions)) (*Channel, error) {
	defaultOptions := getDefaultChannelOptions()
	options := &defaultOptions
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	if conn.connManager == nil {
		return nil, errors.New("connection manager can't be nil")
	}
	chanManager, err := channel.New(conn.connManager, options.Logger, conn.connManager.ReconnectInterval)
	if err != nil {
		return nil, err
	}
	err = chanManager.QosSafe(options.QOSPrefetch, 0, options.QOSGlobal)
	if err != nil {
		_ = chanManager.Close()
		return nil, fmt.Errorf("declare qos failed: %w", err)
	}
	return &Channel{chanManager}, nil
}
