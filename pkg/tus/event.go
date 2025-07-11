package tus

import (
	"context"
	"strings"
	"sync"

	"github.com/busyster996/dagflow/pkg/tus/types"
)

type HandleFn func(event types.HookEvent) error

type subscriber struct {
	ctx    context.Context
	cancel context.CancelFunc
	ch     chan types.HookEvent
}

type topic struct {
	name   string
	logger ILogger
	subs   []*subscriber
	mu     sync.RWMutex
}

func newTopic(name string, logger ILogger) *topic {
	return &topic{
		name:   name,
		logger: logger,
	}
}

func (t *topic) publish(event types.HookEvent) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, sub := range t.subs {
		select {
		case <-sub.ctx.Done():
			continue
		case sub.ch <- event:
		default:
			t.logger.Warnf("channel full, dropping message")
		}
	}
}

func (t *topic) subscribe(ctx context.Context, handler HandleFn) {
	subCtx, cancel := context.WithCancel(ctx)
	sub := &subscriber{
		ctx:    subCtx,
		cancel: cancel,
		ch:     make(chan types.HookEvent, 65535),
	}

	t.mu.Lock()
	t.subs = append(t.subs, sub)
	t.mu.Unlock()

	go func() {
		defer cancel()
		for {
			select {
			case <-sub.ctx.Done():
				t.logger.Infof("closed for topic %s", t.name)
				t.removeSubscriber(sub)
				return
			case event := <-sub.ch:
				if err := handler(event); err != nil {
					t.logger.Errorf("handling event %s", err)
				}
			}
		}
	}()
}

func (t *topic) removeSubscriber(target *subscriber) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, sub := range t.subs {
		if sub == target {
			t.subs = append(t.subs[:i], t.subs[i+1:]...)
			break
		}
	}
}

func (t *topic) close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, sub := range t.subs {
		sub.cancel()
	}
	t.subs = nil
}

type sMemoryBroker struct {
	logger ILogger
	topics sync.Map
}

func newMemoryBroker(logger ILogger) *sMemoryBroker {
	return &sMemoryBroker{logger: logger}
}

func (b *sMemoryBroker) PublishEvent(prefix string, event types.HookEvent) {
	b.topics.Range(func(key, value any) bool {
		if strings.HasPrefix(key.(string), prefix) {
			value.(*topic).publish(event)
		}
		return true
	})
}

func (b *sMemoryBroker) SubscribeEvent(ctx context.Context, prefix string, handler HandleFn) {
	t, _ := b.topics.LoadOrStore(prefix, newTopic(prefix, b.logger))
	t.(*topic).subscribe(ctx, handler)
}

func (b *sMemoryBroker) Shutdown(ctx context.Context) {
	var wg sync.WaitGroup
	b.topics.Range(func(_, value any) bool {
		wg.Add(1)
		go func(t *topic) {
			defer wg.Done()
			t.close()
		}(value.(*topic))
		return true
	})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		b.logger.Warnf("shutdown timed out")
	case <-done:
		b.logger.Infof("shutdown complete")
	}
}
