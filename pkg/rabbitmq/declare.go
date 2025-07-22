package rabbitmq

import (
	"errors"

	"github.com/busyster996/dagflow/pkg/rabbitmq/internal/manager/channel"
)

type Declarator struct {
	chanManager *channel.Manager
}

func NewDeclarator(conn *Conn, optionFuncs ...func(*DeclareOptions)) (*Declarator, error) {
	defaultOptions := getDefaultDeclareOptions()
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

	result := &Declarator{
		chanManager: chanManager,
	}

	return result, nil
}

func (d *Declarator) Close() {
	_ = d.chanManager.Close()
}

func (d *Declarator) Exchange(options ExchangeOptions) error {
	return declareExchange(d.chanManager, options)
}

func (d *Declarator) Queue(options QueueOptions) error {
	return declareQueue(d.chanManager, options)
}

func (d *Declarator) BindExchanges(bindings []Binding) error {
	for _, binding := range bindings {
		err := d.chanManager.ExchangeBindSafe(
			binding.Destination,
			binding.RoutingKey,
			binding.Source,
			binding.NoWait,
			tableToAMQPTable(binding.Args),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Declarator) BindQueues(bindings []Binding) error {
	for _, binding := range bindings {
		err := d.chanManager.QueueBindSafe(
			binding.Destination,
			binding.RoutingKey,
			binding.Source,
			binding.NoWait,
			tableToAMQPTable(binding.Args),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func declareQueue(chanManager *channel.Manager, options QueueOptions) error {
	if !options.Declare {
		return nil
	}
	if options.Passive {
		_, err := chanManager.QueueDeclarePassiveSafe(
			options.Name,
			options.Durable,
			options.AutoDelete,
			options.Exclusive,
			options.NoWait,
			tableToAMQPTable(options.Args),
		)
		if err != nil {
			return err
		}
		return nil
	}
	_, err := chanManager.QueueDeclareSafe(
		options.Name,
		options.Durable,
		options.AutoDelete,
		options.Exclusive,
		options.NoWait,
		tableToAMQPTable(options.Args),
	)
	if err != nil {
		return err
	}
	return nil
}

func declareExchange(chanManager *channel.Manager, options ExchangeOptions) error {
	if !options.Declare {
		return nil
	}
	if options.Passive {
		err := chanManager.ExchangeDeclarePassiveSafe(
			options.Name,
			options.Kind,
			options.Durable,
			options.AutoDelete,
			options.Internal,
			options.NoWait,
			tableToAMQPTable(options.Args),
		)
		if err != nil {
			return err
		}
		return nil
	}
	err := chanManager.ExchangeDeclareSafe(
		options.Name,
		options.Kind,
		options.Durable,
		options.AutoDelete,
		options.Internal,
		options.NoWait,
		tableToAMQPTable(options.Args),
	)
	if err != nil {
		return err
	}
	return nil
}

func declareBindings(chanManager *channel.Manager, options ConsumerOptions) error {
	for _, exchangeOption := range options.ExchangeOptions {
		for _, binding := range exchangeOption.Bindings {
			if !binding.Declare {
				continue
			}
			err := chanManager.QueueBindSafe(
				options.QueueOptions.Name,
				binding.RoutingKey,
				exchangeOption.Name,
				binding.NoWait,
				tableToAMQPTable(binding.Args),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
