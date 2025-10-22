package amqp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"

	"github.com/busyster996/dagflow/internal/pubsub/queue"
	"github.com/busyster996/dagflow/internal/utility"
	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/busyster996/dagflow/pkg/rabbitmq"
)

type sAmqp struct {
	conn *rabbitmq.Conn
	// 防止相同生产者重复创建
	publisherMap map[string]*rabbitmq.Publisher
	// 防止相同消费者重复创建
	consumerMap map[string]*rabbitmq.Consumer
	topics      sync.Map
	directs     sync.Map
	mu          sync.Mutex
}

func (a *sAmqp) directExchangeName() string {
	return utility.ServiceName
}

func (a *sAmqp) topicExchangeName() string {
	return utility.ServiceName + ".topic"
}

func (a *sAmqp) PublishTask(node string, data string) error {
	rkey := fmt.Sprintf("%s.%s", queue.TaskRoutingKey(), node)
	return a.publish(amqp091.ExchangeDirect, rkey, data, "")
}

func (a *sAmqp) PublishTaskDelayed(node string, data string, delay time.Duration) error {
	rkey := fmt.Sprintf("%s.%s", queue.TaskRoutingKey(), node)
	delayedQueue := rkey + ".delayed"
	if a.publisherMap[delayedQueue] == nil {
		channel, err := rabbitmq.NewChannel(a.conn, rabbitmq.WithChannelOptionsLogger(logx.GetSubLogger()))
		if err != nil {
			logx.Errorln(err)
			return err
		}
		defer channel.Close()
		if _, err = channel.QueueDeclareSafe(
			delayedQueue, true, false, false, false,
			amqp091.Table{
				"x-queue-type":              "quorum",
				"x-dead-letter-exchange":    a.directExchangeName(),
				"x-dead-letter-routing-key": rkey,
			}); err != nil {
			logx.Errorln(err)
			return err
		}
		err = channel.QueueBindSafe(delayedQueue, delayedQueue, a.directExchangeName(), false, amqp091.Table{})
		if err != nil {
			logx.Errorln(err)
			return err
		}
	}

	return a.publish(amqp091.ExchangeDirect, delayedQueue, data, fmt.Sprintf("%.f", delay.Seconds()))
}

func (a *sAmqp) SubscribeTask(ctx context.Context, node string, handler utility.QueueHandleFn) error {
	qname := fmt.Sprintf("%s.%s", queue.TaskRoutingKey(), node)
	d, _ := a.directs.LoadOrStore(qname, utility.NewMemDirectQueue(qname))
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.consumerMap[qname]; !ok {
		var err error
		a.consumerMap[qname], err = a.newConsumer(amqp091.ExchangeDirect, qname, qname, false, func(data string) {
			d.(utility.IQueue).Publish(data)
		})
		if err != nil {
			return err
		}
	}
	d.(utility.IQueue).Subscribe(ctx, handler)
	return nil
}

func (a *sAmqp) PublishEvent(data string) error {
	rkey := fmt.Sprintf("%s.*", queue.EventRoutingKey())
	return a.publish(amqp091.ExchangeTopic, rkey, data, "")
}

func (a *sAmqp) SubscribeEvent(ctx context.Context, handler utility.QueueHandleFn) error {
	rkey := fmt.Sprintf("%s.*", queue.EventRoutingKey())
	t, _ := a.topics.LoadOrStore(rkey, utility.NewMemTopicQueue(rkey))
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.consumerMap[rkey]; !ok {
		qname := fmt.Sprintf("%s.%s", queue.EventRoutingKey(), ksuid.New().String())
		var err error
		a.consumerMap[rkey], err = a.newConsumer(amqp091.ExchangeTopic, rkey, qname, true, func(data string) {
			t.(utility.IQueue).Publish(data)
		})
		if err != nil {
			return err
		}
	}
	t.(utility.IQueue).Subscribe(ctx, handler)
	return nil
}

func (a *sAmqp) PublishManager(node string, data string) error {
	routingKey := fmt.Sprintf("%s.%s", queue.ManagerRoutingKey(), node)
	return a.publish(amqp091.ExchangeTopic, routingKey, data, "")
}

func (a *sAmqp) SubscribeManager(ctx context.Context, node string, handler utility.QueueHandleFn) error {
	rkey := fmt.Sprintf("%s.%s", queue.ManagerRoutingKey(), node)
	t, _ := a.topics.LoadOrStore(rkey, utility.NewMemTopicQueue(rkey))
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.consumerMap[rkey]; !ok {
		qname := fmt.Sprintf("%s.%s", queue.ManagerRoutingKey(), node)
		var err error
		a.consumerMap[rkey], err = a.newConsumer(amqp091.ExchangeTopic, rkey, qname, false, func(data string) {
			t.(utility.IQueue).Publish(data)
		})
		if err != nil {
			return err
		}
	}
	t.(utility.IQueue).Subscribe(ctx, handler)
	return nil
}

func (a *sAmqp) Shutdown(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, publisher := range a.publisherMap {
		publisher.Close()
	}
	for _, consumer := range a.consumerMap {
		consumer.Close()
	}
	if a.conn != nil {
		_ = a.conn.Close()
	}
	var wg sync.WaitGroup
	a.directs.Range(func(_, value any) bool {
		wg.Add(1)
		go func(t utility.IQueue) {
			defer wg.Done()
			t.Close()
		}(value.(utility.IQueue))
		return true
	})
	a.topics.Range(func(_, value any) bool {
		wg.Add(1)
		go func(d utility.IQueue) {
			defer wg.Done()
			d.Close()
		}(value.(utility.IQueue))
		return true
	})

	doneChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case <-ctx.Done():
		logx.Infoln("shutting down broker")
	case <-doneChan:
		logx.Infoln("broker shutdown complete")
	}
}

func (a *sAmqp) subscribe(ctx context.Context, consumer *rabbitmq.Consumer, handler utility.QueueHandleFn) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		err := consumer.Run(func(d rabbitmq.Delivery) (action rabbitmq.Action) {
			if handler != nil {
				handler(string(d.Body))
			}
			return rabbitmq.Ack
		})
		if err != nil {
			logx.Errorln("unexpected error occurred while processing task", err)
		}
	}()
	go func() {
		<-ctx.Done()
		logx.Infof("subscribe closed")
		consumer.Close()
	}()
}

func (a *sAmqp) publish(kind, rkey, data, expiration string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	var ename = a.directExchangeName()
	if kind == amqp091.ExchangeTopic {
		ename = a.topicExchangeName()
	}
	publisher, ok := a.publisherMap[ename]
	if !ok {
		var err error
		publisher, err = a.newPublisher(kind, ename)
		if err != nil {
			logx.Errorln(err)
			return err
		}
		a.publisherMap[ename] = publisher
	}
	ops := []func(*rabbitmq.PublishOptions){
		rabbitmq.WithPublishOptionsExchange(ename),    // 交换机名称
		rabbitmq.WithPublishOptionsMandatory,          // 强制发布
		rabbitmq.WithPublishOptionsPersistentDelivery, // 立即发布
	}
	if expiration != "" {
		ops = append(ops, rabbitmq.WithPublishOptionsExpiration(expiration))
	}
	return publisher.Publish(
		[]byte(data), []string{rkey},
		ops...,
	)
}

func (a *sAmqp) newPublisher(kind, ename string) (*rabbitmq.Publisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		a.conn,
		rabbitmq.WithPublisherOptionsLogger(logx.GetSubLoggerWithOption(zap.AddCallerSkip(-1))), // 日志
		rabbitmq.WithPublisherOptionsExchangeName(ename),                                        // 交换机名称
		rabbitmq.WithPublisherOptionsExchangeKind(kind),                                         // 交换机类型
		rabbitmq.WithPublisherOptionsExchangeDeclare,                                            // 声明交换机
		rabbitmq.WithPublisherOptionsExchangeDurable,                                            // 交换机持久化
		rabbitmq.WithPublisherOptionsConfirm,                                                    // 启用确认模式
	)
	if err != nil {
		return nil, err
	}
	publisher.NotifyPublish(func(confirm rabbitmq.Confirmation) {
		logx.Infof("publish success: %v", confirm)
	})
	publisher.NotifyReturn(func(r rabbitmq.Return) {
		logx.Infoln("message returned from server", r.Exchange, r.RoutingKey, r.ReplyCode, r.ReplyText)
	})
	return publisher, nil
}

func (a *sAmqp) newConsumer(kind, rkey, qname string, autoDel bool, handle utility.QueueHandleFn) (*rabbitmq.Consumer, error) {
	var ename = a.directExchangeName()
	if kind == amqp091.ExchangeTopic {
		ename = a.topicExchangeName()
	}
	ops := []func(*rabbitmq.ConsumerOptions){
		rabbitmq.WithConsumerOptionsLogger(logx.GetSubLoggerWithOption(zap.AddCallerSkip(-1))), // 日志
		rabbitmq.WithConsumerOptionsExchangeName(ename),                                        // 交换机名称
		rabbitmq.WithConsumerOptionsRoutingKey(rkey),                                           // routing key
		rabbitmq.WithConsumerOptionsExchangeKind(kind),                                         // 交换机类型
		rabbitmq.WithConsumerOptionsExchangeDeclare,                                            // 声明交换机
		rabbitmq.WithConsumerOptionsExchangeDurable,                                            // 交换机持久化
		rabbitmq.WithConsumerOptionsQueueDurable,                                               // 队列持久化
		rabbitmq.WithConsumerOptionsQueueQuorum,                                                // 使用仲裁队列
	}
	if autoDel {
		ops = append(ops, rabbitmq.WithConsumerOptionsQueueExpires(60*time.Second))
	}
	consumer, err := rabbitmq.NewConsumer(a.conn, qname, ops...)
	if err != nil {
		return nil, err
	}
	a.subscribe(context.Background(), consumer, func(data string) {
		handle(data)
	})
	return consumer, nil
}
