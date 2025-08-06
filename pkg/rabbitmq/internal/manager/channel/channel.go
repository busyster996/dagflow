package channel

import (
	"errors"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/busyster996/dagflow/pkg/rabbitmq/internal/dispatcher"
	"github.com/busyster996/dagflow/pkg/rabbitmq/internal/logger"
	"github.com/busyster996/dagflow/pkg/rabbitmq/internal/manager/connection"
)

// Manager -
type Manager struct {
	logger              logger.Logger
	channel             *amqp.Channel
	connManager         *connection.Manager
	channelMu           *sync.RWMutex
	reconnectInterval   time.Duration
	reconnectionCount   uint
	reconnectionCountMu *sync.Mutex
	dispatcher          *dispatcher.Dispatcher
}

// New creates a new connection manager
func New(connManager *connection.Manager, log logger.Logger, reconnectInterval time.Duration) (*Manager, error) {
	chanManager := &Manager{
		logger:              log,
		connManager:         connManager,
		channelMu:           &sync.RWMutex{},
		reconnectInterval:   reconnectInterval,
		reconnectionCount:   0,
		reconnectionCountMu: &sync.Mutex{},
		dispatcher:          dispatcher.New(),
	}

	ch, err := chanManager.getNewChannel()
	if err != nil {
		return nil, err
	}

	chanManager.channel = ch
	go chanManager.startNotifyCancelOrClosed()

	return chanManager, nil
}

func (m *Manager) getNewChannel() (*amqp.Channel, error) {
	conn := m.connManager.CheckoutConnection()
	defer m.connManager.CheckinConnection()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

// startNotifyCancelOrClosed listens on the channel's cancelled and closed
// notifiers. When it detects a problem, it attempts to reconnect.
// Once reconnected, it sends an error back on the manager's notifyCancelOrClose
// channel
func (m *Manager) startNotifyCancelOrClosed() {
	notifyCloseChan := m.channel.NotifyClose(make(chan *amqp.Error, 1))
	notifyCancelChan := m.channel.NotifyCancel(make(chan string, 1))

	select {
	case err := <-notifyCloseChan:
		if err != nil {
			m.logger.Errorf("attempting to reconnect to amqp server after close with error: %v", err)
			m.reconnectLoop()
			m.logger.Warnf("successfully reconnected to amqp server")
			_ = m.dispatcher.Dispatch(err)
		}

		if err == nil {
			m.logger.Infof("amqp channel closed gracefully")
		}
	case err := <-notifyCancelChan:
		m.logger.Errorf("attempting to reconnect to amqp server after cancel with error: %s", err)
		m.reconnectLoop()
		m.logger.Warnf("successfully reconnected to amqp server after cancel")
		if _err := m.dispatcher.Dispatch(errors.New(err)); _err != nil {
			m.logger.Warnf("channel dispatch err: %v", err)
		}
	}
}

// GetReconnectionCount -
func (m *Manager) GetReconnectionCount() uint {
	m.reconnectionCountMu.Lock()
	defer m.reconnectionCountMu.Unlock()

	return m.reconnectionCount
}

func (m *Manager) incrementReconnectionCount() {
	m.reconnectionCountMu.Lock()
	defer m.reconnectionCountMu.Unlock()

	m.reconnectionCount++
}

// reconnectLoop continuously attempts to reconnect
func (m *Manager) reconnectLoop() {
	for {
		m.logger.Infof("waiting %s seconds to attempt to reconnect to amqp server", m.reconnectInterval)
		time.Sleep(m.reconnectInterval)
		err := m.reconnect()
		if err != nil {
			m.logger.Errorf("error reconnecting to amqp server: %v", err)
		} else {
			m.incrementReconnectionCount()
			go m.startNotifyCancelOrClosed()
			return
		}
	}
}

// reconnect safely closes the current channel and obtains a new one
func (m *Manager) reconnect() error {
	m.channelMu.Lock()
	defer m.channelMu.Unlock()

	newChannel, err := m.getNewChannel()
	if err != nil {
		return err
	}

	if m.channel != nil {
		if err = m.channel.Close(); err != nil {
			m.logger.Warnf("error closing channel while reconnecting: %v", err)
		}
	}

	m.channel = newChannel
	return nil
}

// Close safely closes the current channel and connection
func (m *Manager) Close() error {
	m.logger.Infof("closing channel manager...")
	m.channelMu.Lock()
	defer m.channelMu.Unlock()

	err := m.channel.Close()
	if err != nil {
		m.logger.Errorf("close err: %v", err)
		return err
	}

	return nil
}

// NotifyReconnect adds a new subscriber that will receive error messages whenever
// the connection manager has successfully reconnect to the server
func (m *Manager) NotifyReconnect() (<-chan error, chan<- struct{}) {
	return m.dispatcher.AddSubscriber()
}
