package queue

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/busyster996/dagflow/internal/pubsub/common"
	"github.com/busyster996/dagflow/internal/utils"
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
	d, _ := m.directs.LoadOrStore(routingKey, common.NewMemDirect(routingKey))
	d.(*common.SMemDirect).Publish(data)
	return nil
}

func (m *sMemoryBroker) PublishTaskDelayed(node string, data string, delay time.Duration) error {
	key := utils.MD5(fmt.Sprintf("%s_%s_%s", TaskRoutingKey(), node, data))
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

func (m *sMemoryBroker) SubscribeTask(ctx context.Context, node string, handler common.HandleFn) error {
	routingKey := fmt.Sprintf("%s_%s", TaskRoutingKey(), node)
	d, _ := m.directs.LoadOrStore(routingKey, common.NewMemDirect(routingKey))
	d.(*common.SMemDirect).Subscribe(ctx, handler)
	return nil
}

func (m *sMemoryBroker) PublishEvent(data string) error {
	routingKey := fmt.Sprintf("%s.*", EventRoutingKey())
	m.topics.Range(func(key, value any) bool {
		if wildcard.Match(routingKey, key.(string)) {
			value.(*common.SMemTopic).Publish(data)
		}
		return true
	})
	return nil
}

func (m *sMemoryBroker) SubscribeEvent(ctx context.Context, handler common.HandleFn) error {
	routingKey := fmt.Sprintf("%s.%s", EventRoutingKey(), ksuid.New().String())
	t, _ := m.topics.LoadOrStore(routingKey, common.NewMemTopic(routingKey))
	t.(*common.SMemTopic).Subscribe(ctx, handler)
	return nil
}

func (m *sMemoryBroker) PublishManager(node string, data string) error {
	m.topics.Range(func(key, value any) bool {
		if wildcard.Match(fmt.Sprintf("%s.%s", ManagerRoutingKey(), node), key.(string)) {
			value.(*common.SMemTopic).Publish(data)
		}
		return true
	})
	return nil
}

func (m *sMemoryBroker) SubscribeManager(ctx context.Context, node string, handler common.HandleFn) error {
	qname := fmt.Sprintf("%s.%s", ManagerRoutingKey(), node)
	t, _ := m.topics.LoadOrStore(qname, common.NewMemTopic(qname))
	t.(*common.SMemTopic).Subscribe(ctx, handler)
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
		go func(d *common.SMemDirect) {
			defer wg.Done()
			d.Close()
		}(value.(*common.SMemDirect))
		return true
	})
	m.topics.Range(func(_, value any) bool {
		wg.Add(1)
		go func(t *common.SMemTopic) {
			defer wg.Done()
			t.Close()
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
