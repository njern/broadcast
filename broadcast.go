// Package broadcast provides a simple and efficient way to implement a
// publish/subscribe pattern in Go, allowing one sender to broadcast
// messages to multiple receivers. It supports dynamic subscription
// and unsubscription, ensuring flexibility and control over message
// delivery and resource management.

package broadcast

import (
	"fmt"
	"sync"
	"time"
)

// ErrBroadcasterClosed is returned when trying to subscribe to a closed Broadcaster.
var ErrBroadcasterClosed = fmt.Errorf("broadcaster is closed")

// A Broadcaster broadcasts values to multiple subscribers.
type Broadcaster[T any] struct {
	m           sync.RWMutex // Protects the subscribers slice
	subscribers map[chan<- T]struct{}
	valCh       chan T
	closeCh     chan struct{}
	timeout     time.Duration
}

// New creates a new Broadcaster with a buffer of size `n`
// and a timeout for each subscriber of `timeout`.
func New[T any](n int, timeout time.Duration) *Broadcaster[T] {
	b := &Broadcaster[T]{
		subscribers: make(map[chan<- T]struct{}),
		valCh:       make(chan T, n),
		closeCh:     make(chan struct{}),
		timeout:     timeout,
	}

	go b.run()
	return b
}

// run starts the broadcasting process, listening for new values and subscribers.
func (b *Broadcaster[T]) run() {
	for {
		select {
		case v := <-b.valCh:
			b.broadcast(v)
		case <-b.closeCh:
			return
		}
	}
}

// broadcast the value to all subscribers.
func (b *Broadcaster[T]) broadcast(v T) {
	b.m.RLock()
	defer b.m.RUnlock()

	for ch := range b.subscribers {
		if b.timeout <= 0 {
			select {
			case ch <- v:
			case <-b.closeCh:
				// NOTE(njern): Handle an edge case where the
				// Broadcaster is closed while broadcasting.
				return
			}
			continue
		}

		select {
		case ch <- v:
		case <-time.After(b.timeout):
			// NOTE(njern): The subscriber did not read from the
			// channel within the timeout, keep going.
		case <-b.closeCh:
			// NOTE(njern): Handle an edge case where the
			// Broadcaster is closed while broadcasting.
			return
		}
	}
}

// Subscribe adds a new subscriber to the broadcaster and returns a channel to listen on.
func (b *Broadcaster[T]) Subscribe(chSize int) (chan T, error) {
	b.m.Lock()
	defer b.m.Unlock()

	if b.isClosed() {
		return nil, ErrBroadcasterClosed
	}

	ch := make(chan T, chSize)
	b.subscribers[ch] = struct{}{}
	return ch, nil
}

// Unsubscribe removes a subscriber from the broadcaster.
func (b *Broadcaster[T]) Unsubscribe(ch chan<- T) {
	b.m.Lock()
	defer b.m.Unlock()

	delete(b.subscribers, ch)

	defer func() {
		// NOTE(njern): This may happen if the Broadcaster is closed, which also closes
		// subCh - but the subscriber subsequently calls Unsubscribe.
		_ = recover()
	}()

	close(ch)
}

// Close the broadcaster and all subscriber channels.
func (b *Broadcaster[T]) Close() {
	close(b.closeCh)

	b.m.Lock()
	defer b.m.Unlock()
	for ch := range b.subscribers {
		close(ch)
	}

	b.subscribers = nil
}

// Chan returns the input channel for the broadcaster.
func (b *Broadcaster[T]) Chan() chan<- T {
	return b.valCh
}

// isClosed checks if the broadcaster has been closed.
func (b *Broadcaster[T]) isClosed() bool {
	select {
	case <-b.closeCh:
		return true
	default:
		return false
	}
}
