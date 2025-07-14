package amqp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/busyster996/dagflow/internal/pubsub/common"
	"github.com/busyster996/dagflow/internal/pubsub/queue"
	"github.com/busyster996/dagflow/internal/utils"
	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/busyster996/dagflow/pkg/rabbitmq"
)

type sAmqp struct {
	conn    *rabbitmq.Conn
	channel *rabbitmq.Channel
	// 防止相同生产者重复创建
	publisherMap map[string]*rabbitmq.Publisher
	// 防止相同消费者重复创建
	consumerMap map[string]*rabbitmq.Consumer
	topics      sync.Map
	directs     sync.Map
	mu          sync.Mutex
}

func (a *sAmqp) directExchangeName() string {
	return utils.ServiceName
}

func (a *sAmqp) topicExchangeName() string {
	return utils.ServiceName + ".topic"
}

func (a *sAmqp) PublishTask(node string, data string) error {
	rkey := fmt.Sprintf("%s.%s", queue.TaskRoutingKey(), node)
	return a.publishDirect(rkey, data)
}

func (a *sAmqp) PublishTaskDelayed(node string, data string, delay time.Duration) error {
	rkey := fmt.Sprintf("%s.%s", queue.TaskRoutingKey(), node)
	delayedQueue := rkey + ".delayed"
	_, err := a.channel.QueueDeclareSafe(
		delayedQueue, true, false, false, false,
		amqp091.Table{
			"x-dead-letter-exchange":    a.directExchangeName(),
			"x-dead-letter-routing-key": rkey,
		})
	if err != nil {
		logx.Errorln(err)
		return err
	}
	err = a.channel.QueueBindSafe(delayedQueue, delayedQueue, a.directExchangeName(), false, amqp091.Table{})
	if err != nil {
		logx.Errorln(err)
		return err
	}
	return a.publish(a.directExchangeName(), delayedQueue, data, fmt.Sprintf("%.f", delay.Seconds()))
}

func (a *sAmqp) SubscribeTask(ctx context.Context, node string, handler common.HandleFn) error {
	qname := fmt.Sprintf("%s.%s", queue.TaskRoutingKey(), node)
	d, _ := a.directs.LoadOrStore(qname, common.NewMemDirect(qname))
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.consumerMap[qname]; !ok {
		var err error
		a.consumerMap[qname], err = a.newDirectConsumer(qname, func(data string) {
			d.(*common.SMemDirect).Publish(data)
		})
		if err != nil {
			return err
		}
	}
	d.(*common.SMemDirect).Subscribe(ctx, handler)
	return nil
}

func (a *sAmqp) PublishEvent(data string) error {
	rkey := fmt.Sprintf("%s.*", queue.EventRoutingKey())
	return a.publishTopic(rkey, data)
}

func (a *sAmqp) SubscribeEvent(ctx context.Context, handler common.HandleFn) error {
	rkey := fmt.Sprintf("%s.*", queue.EventRoutingKey())
	t, _ := a.topics.LoadOrStore(rkey, common.NewMemTopic(rkey))
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.consumerMap[rkey]; !ok {
		qname := queue.EventRoutingKey()
		var err error
		a.consumerMap[rkey], err = a.newTopicConsumer(rkey, qname, func(data string) {
			t.(*common.SMemTopic).Publish(data)
		})
		if err != nil {
			return err
		}
	}
	t.(*common.SMemTopic).Subscribe(ctx, handler)
	return nil
}

func (a *sAmqp) PublishManager(node string, data string) error {
	routingKey := fmt.Sprintf("%s.%s", queue.ManagerRoutingKey(), node)
	return a.publishTopic(routingKey, data)
}

func (a *sAmqp) SubscribeManager(ctx context.Context, node string, handler common.HandleFn) error {
	rkey := fmt.Sprintf("%s.%s", queue.ManagerRoutingKey(), node)
	t, _ := a.topics.LoadOrStore(rkey, common.NewMemTopic(rkey))
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.consumerMap[rkey]; !ok {
		qname := fmt.Sprintf("%s.%s", queue.ManagerRoutingKey(), node)
		var err error
		a.consumerMap[rkey], err = a.newTopicConsumer(rkey, qname, func(data string) {
			t.(*common.SMemTopic).Publish(data)
		})
		if err != nil {
			return err
		}
	}
	t.(*common.SMemTopic).Subscribe(ctx, handler)
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
		go func(t *common.SMemDirect) {
			defer wg.Done()
			t.Close()
		}(value.(*common.SMemDirect))
		return true
	})
	a.topics.Range(func(_, value any) bool {
		wg.Add(1)
		go func(d *common.SMemTopic) {
			defer wg.Done()
			d.Close()
		}(value.(*common.SMemTopic))
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

func (a *sAmqp) subscribe(ctx context.Context, consumer *rabbitmq.Consumer, handler common.HandleFn) {
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

func (a *sAmqp) publishDirect(rkey, data string) error {
	return a.publish(a.directExchangeName(), rkey, data, "")
}

func (a *sAmqp) publishTopic(rkey, data string) error {
	return a.publish(a.topicExchangeName(), rkey, data, "")
}

func (a *sAmqp) publish(ename, rkey, data, expiration string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	publisher, ok := a.publisherMap[ename]
	if !ok {
		return fmt.Errorf("exchange %s publisher not found", ename)
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

func (a *sAmqp) newDirectPublisher() (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.publisherMap[a.directExchangeName()], err = a.newPublisher(amqp091.ExchangeDirect, a.directExchangeName())
	return
}

func (a *sAmqp) newTopicPublisher() (err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.publisherMap[a.topicExchangeName()], err = a.newPublisher(amqp091.ExchangeTopic, a.topicExchangeName())
	return
}

func (a *sAmqp) newPublisher(kind, ename string) (*rabbitmq.Publisher, error) {
	return rabbitmq.NewPublisher(
		a.conn,
		rabbitmq.WithPublisherOptionsLogger(logx.GetSubLoggerWithOption(zap.AddCallerSkip(-1))), // 日志
		rabbitmq.WithPublisherOptionsExchangeName(ename),                                        // 交换机名称
		rabbitmq.WithPublisherOptionsExchangeKind(kind),                                         // 交换机类型
		rabbitmq.WithPublisherOptionsExchangeDeclare,                                            // 声明交换机
		rabbitmq.WithPublisherOptionsExchangeDurable,                                            // 交换机持久化
	)
}

func (a *sAmqp) newDirectConsumer(qname string, handle common.HandleFn) (*rabbitmq.Consumer, error) {
	return a.newConsumer(amqp091.ExchangeDirect, a.directExchangeName(), qname, qname, handle)
}

func (a *sAmqp) newTopicConsumer(rkey, qname string, handle common.HandleFn) (*rabbitmq.Consumer, error) {
	return a.newConsumer(amqp091.ExchangeTopic, a.topicExchangeName(), rkey, qname, handle)
}

func (a *sAmqp) newConsumer(kind, ename, rkey, qname string, handle common.HandleFn) (*rabbitmq.Consumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		a.conn, qname,
		rabbitmq.WithConsumerOptionsLogger(logx.GetSubLoggerWithOption(zap.AddCallerSkip(-1))), // 日志
		rabbitmq.WithConsumerOptionsExchangeName(ename),                                        // 交换机名称
		rabbitmq.WithConsumerOptionsRoutingKey(rkey),                                           // routing key
		rabbitmq.WithConsumerOptionsExchangeKind(kind),                                         // 交换机类型
		rabbitmq.WithConsumerOptionsExchangeDeclare,                                            // 声明交换机
		rabbitmq.WithConsumerOptionsExchangeDurable,                                            // 交换机持久化
		rabbitmq.WithConsumerOptionsQueueDurable,                                               // 队列持久化
		rabbitmq.WithConsumerOptionsQueueQuorum,                                                // 使用仲裁队列
	)
	if err != nil {
		return nil, err
	}
	a.subscribe(context.Background(), consumer, func(data string) {
		handle(data)
	})
	return consumer, nil
}
