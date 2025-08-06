package channel

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumeSafe safely wraps the (*amqp.Channel).Consume method
func (m *Manager) ConsumeSafe(
	queue,
	consumer string,
	autoAck,
	exclusive,
	noLocal,
	noWait bool,
	args amqp.Table,
) (<-chan amqp.Delivery, error) {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.Consume(
		queue,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)
}

// QueueDeclarePassiveSafe safely wraps the (*amqp.Channel).QueueDeclarePassive method
func (m *Manager) QueueDeclarePassiveSafe(
	name string,
	durable bool,
	autoDelete bool,
	exclusive bool,
	noWait bool,
	args amqp.Table,
) (amqp.Queue, error) {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.QueueDeclarePassive(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
}

// QueueDeclareSafe safely wraps the (*amqp.Channel).QueueDeclare method
func (m *Manager) QueueDeclareSafe(
	name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table,
) (amqp.Queue, error) {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
}

// ExchangeDeclarePassiveSafe safely wraps the (*amqp.Channel).ExchangeDeclarePassive method
func (m *Manager) ExchangeDeclarePassiveSafe(
	name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table,
) error {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.ExchangeDeclarePassive(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args,
	)
}

// ExchangeDeclareSafe safely wraps the (*amqp.Channel).ExchangeDeclare method
func (m *Manager) ExchangeDeclareSafe(
	name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table,
) error {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args,
	)
}

// ExchangeBindSafe safely wraps the (*amqp.Channel).ExchangeBind method
func (m *Manager) ExchangeBindSafe(
	name string, key string, exchange string, noWait bool, args amqp.Table,
) error {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.ExchangeBind(
		name,
		key,
		exchange,
		noWait,
		args,
	)
}

// QueueBindSafe safely wraps the (*amqp.Channel).QueueBind method
func (m *Manager) QueueBindSafe(
	name string, key string, exchange string, noWait bool, args amqp.Table,
) error {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.QueueBind(
		name,
		key,
		exchange,
		noWait,
		args,
	)
}

// QosSafe safely wraps the (*amqp.Channel).Qos method
func (m *Manager) QosSafe(
	prefetchCount int, prefetchSize int, global bool,
) error {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.Qos(
		prefetchCount,
		prefetchSize,
		global,
	)
}

/*
PublishSafe safely wraps the (*amqp.Channel).Publish method.
*/
func (m *Manager) PublishSafe(
	exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing,
) error {
	return m.PublishWithContextSafe(
		context.Background(),
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
}

// PublishWithContextSafe safely wraps the (*amqp.Channel).PublishWithContext method.
func (m *Manager) PublishWithContextSafe(
	ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing,
) error {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()
	confirm, err := m.channel.PublishWithDeferredConfirmWithContext(
		ctx,
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
	if err != nil {
		return err
	}
	if confirm != nil {
		var ok bool
		ok, err = confirm.WaitContext(ctx)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("message publishing not confirmed")
		}
	}
	return nil
}

func (m *Manager) PublishWithDeferredConfirmWithContextSafe(
	ctx context.Context, exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing,
) (*amqp.DeferredConfirmation, error) {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.PublishWithDeferredConfirmWithContext(
		ctx,
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
}

// NotifyReturnSafe safely wraps the (*amqp.Channel).NotifyReturn method
func (m *Manager) NotifyReturnSafe(
	c chan amqp.Return,
) chan amqp.Return {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.NotifyReturn(
		c,
	)
}

// ConfirmSafe safely wraps the (*amqp.Channel).Confirm method
func (m *Manager) ConfirmSafe(
	noWait bool,
) error {
	m.channelMu.Lock()
	defer m.channelMu.Unlock()

	return m.channel.Confirm(
		noWait,
	)
}

// NotifyPublishSafe safely wraps the (*amqp.Channel).NotifyPublish method
func (m *Manager) NotifyPublishSafe(
	confirm chan amqp.Confirmation,
) chan amqp.Confirmation {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.NotifyPublish(
		confirm,
	)
}

// NotifyFlowSafe safely wraps the (*amqp.Channel).NotifyFlow method
func (m *Manager) NotifyFlowSafe(
	c chan bool,
) chan bool {
	m.channelMu.RLock()
	defer m.channelMu.RUnlock()

	return m.channel.NotifyFlow(
		c,
	)
}
