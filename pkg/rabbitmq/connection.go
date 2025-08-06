package rabbitmq

import (
	"math/rand"
	"slices"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/busyster996/dagflow/pkg/rabbitmq/internal/manager/connection"
)

// Conn manages the connection to a rabbit cluster
// it is intended to be shared across publishers and consumers
type Conn struct {
	connManager                *connection.Manager
	reconnectErrCh             <-chan error
	closeConnectionToManagerCh chan<- struct{}

	options ConnectionOptions
}

// Config wraps amqp.Config
// Config is used in DialConfig and Open to specify the desired tuning
// parameters used during a connection open handshake.  The negotiated tuning
// will be stored in the returned connection's Config field.
type Config amqp.Config

type Resolver = connection.Resolver

type StaticResolver struct {
	urls   []string
	shuffe bool
}

func (r *StaticResolver) Resolve() ([]string, error) {
	var urls = slices.Clone(r.urls)
	if r.shuffe {
		rand.Shuffle(len(urls), func(i, j int) {
			urls[i], urls[j] = urls[j], urls[i]
		})
	}
	return urls, nil
}

func NewStaticResolver(urls []string, shuffle bool) *StaticResolver {
	return &StaticResolver{urls: urls, shuffe: shuffle}
}

// NewConn creates a new connection manager
func NewConn(url string, opts ...func(*ConnectionOptions)) (*Conn, error) {
	return NewClusterConn(NewStaticResolver([]string{url}, false), opts...)
}

func NewClusterConn(resolver Resolver, opts ...func(*ConnectionOptions)) (*Conn, error) {
	defaultOptions := getDefaultConnectionOptions()
	options := &defaultOptions
	for _, optFn := range opts {
		optFn(options)
	}

	conn := &Conn{
		options: *options,
	}
	var err error
	conn.connManager, err = connection.New(resolver, amqp.Config(options.Config), options.Logger, options.ReconnectInterval)
	if err != nil {
		return nil, err
	}

	conn.reconnectErrCh, conn.closeConnectionToManagerCh = conn.connManager.NotifyReconnect()
	go conn.handleRestarts()
	return conn, nil
}

func (conn *Conn) handleRestarts() {
	for err := range conn.reconnectErrCh {
		conn.options.Logger.Infof("successful connection recovery from: %v", err)
	}
}

// Close closes the connection, it's not safe for re-use.
// You should also close any consumers and publishers before
// closing the connection
func (conn *Conn) Close() error {
	conn.closeConnectionToManagerCh <- struct{}{}
	return conn.connManager.Close()
}

// IsClosed returns whether the connection is closed or not
func (conn *Conn) IsClosed() bool {
	return conn.connManager.IsClosed()
}
