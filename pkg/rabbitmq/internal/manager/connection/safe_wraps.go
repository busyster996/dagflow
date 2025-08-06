package connection

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// NotifyBlockedSafe safely wraps the (*amqp.Connection).NotifyBlocked method
func (m *Manager) NotifyBlockedSafe(
	receiver chan amqp.Blocking,
) chan amqp.Blocking {
	m.connectionMu.Lock()
	defer m.connectionMu.Unlock()

	// add receiver to connection manager.
	m.publisherNotifyBlockingReceiversMu.Lock()
	m.publisherNotifyBlockingReceivers = append(m.publisherNotifyBlockingReceivers, receiver)
	m.publisherNotifyBlockingReceiversMu.Unlock()

	if !m.universalNotifyBlockingReceiverUsed {
		m.connection.NotifyBlocked(
			m.universalNotifyBlockingReceiver,
		)
		m.universalNotifyBlockingReceiverUsed = true
	}

	return receiver
}

// readUniversalBlockReceiver reads on universal blocking receiver and broadcasts event to all blocking receivers of
// connection manager.
func (m *Manager) readUniversalBlockReceiver() {
	for b := range m.universalNotifyBlockingReceiver {
		m.publisherNotifyBlockingReceiversMu.RLock()
		for _, br := range m.publisherNotifyBlockingReceivers {
			br <- b
		}
		m.publisherNotifyBlockingReceiversMu.RUnlock()
	}
}

func (m *Manager) RemovePublisherBlockingReceiver(receiver chan amqp.Blocking) {
	m.publisherNotifyBlockingReceiversMu.Lock()
	for i, br := range m.publisherNotifyBlockingReceivers {
		if br == receiver {
			m.publisherNotifyBlockingReceivers = append(m.publisherNotifyBlockingReceivers[:i], m.publisherNotifyBlockingReceivers[i+1:]...)
		}
	}
	m.publisherNotifyBlockingReceiversMu.Unlock()
	close(receiver)
}
