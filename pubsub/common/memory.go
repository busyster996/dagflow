package common

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/segmentio/ksuid"
)

type HandleFn func(data string)

const defaultQueueSize = 2 << 16

type SMemDirect struct {
	name    string
	ch      chan string
	subs    []*sSub
	unacked int32
	closed  atomic.Bool
	mu      sync.RWMutex
	wg      sync.WaitGroup
}

func NewMemDirect(name string) *SMemDirect {
	return &SMemDirect{
		name: name,
		ch:   make(chan string, defaultQueueSize),
		subs: make([]*sSub, 0),
	}
}

type sSub struct {
	ctx    context.Context
	cancel context.CancelFunc
	cname  string

	// topic only
	ch chan string
}

// Publish messages to all subscribers in a non-blocking manner.
func (d *SMemDirect) Publish(data string) {
	if d.closed.Load() {
		return
	}
	select {
	case d.ch <- data:
	default:
		// queue full, dropping message
	}
}

func (d *SMemDirect) Size() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.ch)
}

func (d *SMemDirect) Close() {
	if !d.closed.CompareAndSwap(false, true) {
		return
	}
	close(d.ch)

	d.mu.Lock()
	for _, sub := range d.subs {
		sub.cancel()
	}
	d.mu.Unlock()

	d.wg.Wait()
}

func (d *SMemDirect) Subscribe(ctx context.Context, handle HandleFn) {
	subCtx, cancel := context.WithCancel(ctx)
	sub := &sSub{
		cname:  ksuid.New().String(),
		ctx:    subCtx,
		cancel: cancel,
		ch:     make(chan string, 100),
	}
	// Add subscription safely
	d.mu.Lock()
	d.subs = append(d.subs, sub)
	d.mu.Unlock()

	d.wg.Add(1)

	// Handle subscription in a separate goroutine
	go func() {
		defer d.wg.Done()
		defer cancel()

		for {
			select {
			case <-sub.ctx.Done():
				d.removeSubscriber(sub.cname)
				return
			case msg, ok := <-d.ch:
				if !ok {
					return
				}
				atomic.AddInt32(&d.unacked, 1)
				if handle != nil {
					handle(msg)
				}
				atomic.AddInt32(&d.unacked, -1)
			}
		}
	}()
}

func (d *SMemDirect) removeSubscriber(cname string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, sub := range d.subs {
		if sub.cname == cname {
			d.subs = append(d.subs[:i], d.subs[i+1:]...)
			break
		}
	}
}

type SMemTopic struct {
	name      string
	subs      []*sSub
	terminate chan struct{}
	mu        sync.RWMutex
}

func NewMemTopic(name string) *SMemTopic {
	return &SMemTopic{
		name:      name,
		terminate: make(chan struct{}),
	}
}

func (t *SMemTopic) Publish(event string) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, sub := range t.subs {
		select {
		case <-sub.ctx.Done():
			continue
		case sub.ch <- event:
		default:
			// queue full, dropping message
		}
	}
}

func (t *SMemTopic) Subscribe(ctx context.Context, handler HandleFn) {
	subCtx, cancel := context.WithCancel(ctx)
	sub := &sSub{
		cname:  ksuid.New().String(),
		ctx:    subCtx,
		cancel: cancel,
		ch:     make(chan string, 100),
	}
	t.mu.Lock()
	t.subs = append(t.subs, sub)
	t.mu.Unlock()

	// Launch subscriber handling in a separate goroutine
	go func() {
		defer cancel()

		for {
			select {
			case <-sub.ctx.Done():
				t.removeSubscriber(sub.cname)
				return
			case m := <-sub.ch:
				if handler != nil {
					handler(m)
				}
			}
		}
	}()
}

func (t *SMemTopic) removeSubscriber(cname string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, sub := range t.subs {
		if sub.cname == cname {
			t.subs = append(t.subs[:i], t.subs[i+1:]...)
			break
		}
	}
}

func (t *SMemTopic) Close() {
	close(t.terminate)

	t.mu.Lock()
	defer t.mu.Unlock()

	for _, sub := range t.subs {
		sub.cancel()
	}
}
