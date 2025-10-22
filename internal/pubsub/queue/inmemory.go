package queue

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/busyster996/dagflow/internal/utility"
	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/busyster996/dagflow/pkg/wildcard"
)

type delayed struct {
	node  string
	data  string
	timer *time.Timer
}

type sMemoryBroker struct {
	directs sync.Map
	topics  sync.Map

	terminate atomic.Bool

	delayed sync.Map
}

func (m *sMemoryBroker) PublishTask(node string, data string) error {
	routingKey := fmt.Sprintf("%s_%s", TaskRoutingKey(), node)
	d, _ := m.directs.LoadOrStore(routingKey, utility.NewMemDirectQueue(routingKey))
	d.(utility.IQueue).Publish(data)
	return nil
}

func (m *sMemoryBroker) PublishTaskDelayed(node string, data string, delay time.Duration) error {
	key := utility.MD5(fmt.Sprintf("%s_%s_%s", TaskRoutingKey(), node, data))
	t := time.AfterFunc(delay, func() {
		err := m.PublishTask(node, data)
		if err != nil {
			logx.Errorln("error publishing delayed task:", err)
		}

		value, loaded := m.delayed.LoadAndDelete(key)
		if !loaded {
			return
		}
		if t, ok := value.(*delayed); ok {
			t.timer.Stop()
		}
	})
	m.delayed.Store(key, &delayed{
		node:  node,
		data:  data,
		timer: t,
	})
	return nil
}

func (m *sMemoryBroker) SubscribeTask(ctx context.Context, node string, handler utility.QueueHandleFn) error {
	routingKey := fmt.Sprintf("%s_%s", TaskRoutingKey(), node)
	d, _ := m.directs.LoadOrStore(routingKey, utility.NewMemDirectQueue(routingKey))
	d.(utility.IQueue).Subscribe(ctx, handler)
	return nil
}

func (m *sMemoryBroker) PublishEvent(data string) error {
	routingKey := fmt.Sprintf("%s.*", EventRoutingKey())
	m.topics.Range(func(key, value any) bool {
		if wildcard.Match(routingKey, key.(string)) {
			value.(utility.IQueue).Publish(data)
		}
		return true
	})
	return nil
}

func (m *sMemoryBroker) SubscribeEvent(ctx context.Context, handler utility.QueueHandleFn) error {
	routingKey := fmt.Sprintf("%s.%s", EventRoutingKey(), ksuid.New().String())
	t, _ := m.topics.LoadOrStore(routingKey, utility.NewMemTopicQueue(routingKey))
	t.(utility.IQueue).Subscribe(ctx, handler)
	return nil
}

func (m *sMemoryBroker) PublishManager(node string, data string) error {
	m.topics.Range(func(key, value any) bool {
		if wildcard.Match(fmt.Sprintf("%s.%s", ManagerRoutingKey(), node), key.(string)) {
			value.(utility.IQueue).Publish(data)
		}
		return true
	})
	return nil
}

func (m *sMemoryBroker) SubscribeManager(ctx context.Context, node string, handler utility.QueueHandleFn) error {
	qname := fmt.Sprintf("%s.%s", ManagerRoutingKey(), node)
	t, _ := m.topics.LoadOrStore(qname, utility.NewMemTopicQueue(qname))
	t.(utility.IQueue).Subscribe(ctx, handler)
	return nil
}

func (m *sMemoryBroker) Shutdown(ctx context.Context) {
	if !m.terminate.CompareAndSwap(false, true) {
		return
	}

	m.delayed.Range(func(key, value any) bool {
		if t, ok := value.(*delayed); ok {
			t.timer.Stop()
		}
		return true
	})

	var wg sync.WaitGroup
	m.directs.Range(func(_, value any) bool {
		wg.Add(1)
		go func(d utility.IQueue) {
			defer wg.Done()
			d.Close()
		}(value.(utility.IQueue))
		return true
	})
	m.topics.Range(func(_, value any) bool {
		wg.Add(1)
		go func(t utility.IQueue) {
			defer wg.Done()
			t.Close()
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
