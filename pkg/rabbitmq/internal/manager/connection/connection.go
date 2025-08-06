package connection

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/busyster996/dagflow/pkg/rabbitmq/internal/dispatcher"
	"github.com/busyster996/dagflow/pkg/rabbitmq/internal/logger"
)

// Manager -
type Manager struct {
	logger              logger.Logger
	resolver            Resolver
	connection          *amqp.Connection
	amqpConfig          amqp.Config
	connectionMu        *sync.RWMutex
	ReconnectInterval   time.Duration
	reconnectionCount   uint
	reconnectionCountMu *sync.Mutex
	dispatcher          *dispatcher.Dispatcher

	// universalNotifyBlockingReceiver receives block signal from underlying
	// connection which are broadcasted to all publisherNotifyBlockingReceivers
	universalNotifyBlockingReceiver     chan amqp.Blocking
	universalNotifyBlockingReceiverUsed bool
	publisherNotifyBlockingReceiversMu  *sync.RWMutex
	publisherNotifyBlockingReceivers    []chan amqp.Blocking
}

type Resolver interface {
	Resolve() ([]string, error)
}

// dial will attempt to connect to the a list of urls in the order they are
// given.
func (m *Manager) dial() (*amqp.Connection, error) {
	urls, err := m.resolver.Resolve()
	if err != nil {
		return nil, fmt.Errorf("error resolving amqp server urls: %w", err)
	}

	var errs []error
	for _, _url := range urls {
		conn, err := amqp.DialConfig(_url, m.amqpConfig)
		if err == nil {
			return conn, err
		}
		m.logger.Warnf("failed to connect to amqp server %s: %v", m.maskPassword(_url), err)
		errs = append(errs, err)
	}
	return nil, errors.Join(errs...)
}

// New creates a new connection manager
func New(resolver Resolver, conf amqp.Config, log logger.Logger, reconnectInterval time.Duration) (*Manager, error) {
	connManager := Manager{
		logger:                             log,
		resolver:                           resolver,
		amqpConfig:                         conf,
		connectionMu:                       &sync.RWMutex{},
		ReconnectInterval:                  reconnectInterval,
		reconnectionCount:                  0,
		reconnectionCountMu:                &sync.Mutex{},
		dispatcher:                         dispatcher.New(),
		universalNotifyBlockingReceiver:    make(chan amqp.Blocking),
		publisherNotifyBlockingReceiversMu: &sync.RWMutex{},
	}
	conn, err := connManager.dial()
	if err != nil {
		return nil, err
	}
	connManager.connection = conn
	go connManager.startNotifyClose()
	go connManager.readUniversalBlockReceiver()
	return &connManager, nil
}

// Close safely closes the current channel and connection
func (m *Manager) Close() error {
	m.logger.Infof("closing connection manager...")
	m.connectionMu.Lock()
	defer m.connectionMu.Unlock()

	err := m.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

// NotifyReconnect adds a new subscriber that will receive error messages whenever
// the connection manager has successfully reconnected to the server
func (m *Manager) NotifyReconnect() (<-chan error, chan<- struct{}) {
	return m.dispatcher.AddSubscriber()
}

// CheckoutConnection -
func (m *Manager) CheckoutConnection() *amqp.Connection {
	m.connectionMu.RLock()
	return m.connection
}

// CheckinConnection -
func (m *Manager) CheckinConnection() {
	m.connectionMu.RUnlock()
}

// startNotifyCancelOrClosed listens on the channel's cancelled and closed
// notifiers. When it detects a problem, it attempts to reconnect.
// Once reconnected, it sends an error back on the manager's notifyCancelOrClose
// channel
func (m *Manager) startNotifyClose() {
	notifyCloseChan := m.connection.NotifyClose(make(chan *amqp.Error, 1))

	err := <-notifyCloseChan
	if err != nil {
		m.logger.Errorf("attempting to reconnect to amqp server after connection close with error: %v", err)
		m.reconnectLoop()
		m.logger.Warnf("successfully reconnected to amqp server")
		_ = m.dispatcher.Dispatch(err)
	}
	if err == nil {
		m.logger.Infof("amqp connection closed gracefully")
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
		m.logger.Infof("waiting %s seconds to attempt to reconnect to amqp server", m.ReconnectInterval)
		time.Sleep(m.ReconnectInterval)
		err := m.reconnect()
		if err != nil {
			m.logger.Errorf("error reconnecting to amqp server: %v", err)
		} else {
			m.incrementReconnectionCount()
			go m.startNotifyClose()
			return
		}
	}
}

// reconnect safely closes the current channel and obtains a new one
func (m *Manager) reconnect() error {
	m.connectionMu.Lock()
	defer m.connectionMu.Unlock()

	conn, err := m.dial()
	if err != nil {
		return err
	}

	if err = m.connection.Close(); err != nil {
		m.logger.Warnf("error closing connection while reconnecting: %v", err)
	}

	m.connection = conn
	return nil
}

// IsClosed checks if the connection is closed
func (m *Manager) IsClosed() bool {
	m.connectionMu.Lock()
	defer m.connectionMu.Unlock()

	return m.connection.IsClosed()
}

func (m *Manager) maskPassword(urlToMask string) string {
	parsedUrl, _ := url.Parse(urlToMask)
	return parsedUrl.Redacted()
}
