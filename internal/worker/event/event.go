package event

import (
	"errors"
	"sync"
	"sync/atomic"
)

// Package event provides a simple generic event emitter that reads events from
// a source channel and fans them out to multiple subscribers. It is safe for
// concurrent use: subscribers may subscribe and unsubscribe concurrently while
// events are being processed.
//
// Naming choices follow Go conventions:
// - Exported types and methods are capitalized (Emitter, Subscribe, Unsubscribe).
// - Implementation types are unexported (emitter, subscriber).
// - Interface name uses -er suffix (Emitter) rather than an I-prefix.
// - Avoid stuttering: package event + type Emitter -> event.Emitter.

// EventStream is a receive-only channel representing an event stream of T.
type Stream[T any] <-chan T

// Emitter is the public interface for the generic event emitter. It allows
// callers to subscribe to events and to unsubscribe by id.
type Emitter[T any] interface {
	// Subscribe registers a new subscriber and returns a stream that will
	// deliver events to that subscriber, a unique subscriber id, and an error
	// if subscription is not possible (for example, if the emitter has been closed).
	Subscribe() (Stream[T], int64, error)

	// Unsubscribe removes the subscriber with the given id and closes the
	// subscriber's channel so the receiver can detect completion.
	Unsubscribe(id int64)
}

const (
	// defaultBufferSize is the default buffer size for subscriber channels.
	// Extracted as a constant to avoid magic numbers and to make it easier to
	// expose configurability in the future.
	defaultBufferSize = 200
)

// emitter is the concrete (unexported) implementation of Emitter[T].
// It listens for events on an input stream and forwards them to all
// registered subscribers.
type emitter[T any] struct {
	stream      Stream[T]
	subscribers sync.Map // map[int64]*subscriber[T]
	done        bool
	mux         sync.Mutex
	nextID      int64
}

// New constructs a new generic emitter that reads events from the provided
// EventStream. It starts a background goroutine to forward events to
// subscribers and returns the public Emitter[T] interface.
func New[T any](s Stream[T]) Emitter[T] {
	e := &emitter[T]{stream: s}
	go e.process()
	return e
}

// process continuously receives events from the source stream and forwards
// each event to all current subscribers. When the source stream is closed,
// process calls shutdown to close all subscriber channels.
func (e *emitter[T]) process() {
	for event := range e.stream {
		e.subscribers.Range(func(_, value any) bool {
			sub := value.(*subscriber[T])
			sub.send(event)
			return true
		})
	}
	e.shutdown()
}

// shutdown closes all subscribers and marks the emitter as done. It is safe
// to call multiple times; subsequent calls are no-ops.
func (e *emitter[T]) shutdown() {
	e.mux.Lock()
	defer e.mux.Unlock()

	if e.done {
		return
	}
	e.done = true

	e.subscribers.Range(func(_, value any) bool {
		sub := value.(*subscriber[T])
		sub.close()
		return true
	})
}

// Subscribe registers a new subscriber and returns its EventStream along with
// a unique id. If the emitter has already been closed, Subscribe returns an
// error. The returned stream will be closed when the subscriber is
// unsubscribed or when the emitter is closed.
func (e *emitter[T]) Subscribe() (Stream[T], int64, error) {
	e.mux.Lock()
	defer e.mux.Unlock()

	if e.done {
		return nil, 0, errors.New("emitter has been closed, cannot subscribe")
	}
	id := atomic.AddInt64(&e.nextID, 1)
	sub := newSubscriber[T](defaultBufferSize)
	e.subscribers.Store(id, sub)
	return sub.channel(), id, nil
}

// Unsubscribe removes the subscriber with the specified id if it exists and
// closes its channel so the subscriber goroutine can exit. If the id does
// not exist this is a no-op.
func (e *emitter[T]) Unsubscribe(id int64) {
	value, exist := e.subscribers.Load(id)
	if !exist {
		return
	}
	sub := value.(*subscriber[T])
	e.subscribers.Delete(id)
	sub.close()
}

// subscriber represents an individual subscriber. Each subscriber owns a
// buffered channel used to deliver events. The close operation is guarded by
// sync.Once so it can be called safely multiple times.
type subscriber[T any] struct {
	ch   chan T
	once sync.Once
}

// send attempts to deliver the event to the subscriber's buffer in a
// non-blocking manner. If the buffer is full the event is dropped. This
// avoids a slow subscriber from blocking the emitter.
func (s *subscriber[T]) send(event T) {
	select {
	case s.ch <- event:
	default:
		// drop when buffer is full
	}
}

// channel returns the subscriber's receive-only channel.
func (s *subscriber[T]) channel() Stream[T] {
	return s.ch
}

// close closes the subscriber's channel (only once).
func (s *subscriber[T]) close() {
	s.once.Do(func() {
		close(s.ch)
	})
}

// newSubscriber creates a subscriber with the given buffer size.
func newSubscriber[T any](bufferSize int) *subscriber[T] {
	return &subscriber[T]{
		ch: make(chan T, bufferSize),
	}
}
