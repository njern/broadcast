package broadcast

import (
	"sync"
	"testing"
	"time"
)

func TestBroadcastBasic(t *testing.T) {
	b := New[int](10, 0)
	defer b.Close()

	received := make(chan int, 10)
	subCh, err := b.Subscribe(10)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	defer b.Unsubscribe(subCh)

	go func() {
		for val := range subCh {
			received <- val
		}
	}()

	b.Chan() <- 1
	b.Chan() <- 2

	// Allow some time for messages to be received
	time.Sleep(100 * time.Millisecond)

	select {
	case val := <-received:
		if val != 1 {
			t.Errorf("Expected to receive 1, got %d", val)
		}
	default:
		t.Errorf("Expected to receive a message but did not")
	}

	select {
	case val := <-received:
		if val != 2 {
			t.Errorf("Expected to receive 2, got %d", val)
		}
	default:
		t.Errorf("Expected to receive a second message but did not")
	}
}

func TestBroadcastChannelWithTimeout(t *testing.T) {
	b := New[int](10, 50*time.Millisecond)
	defer b.Close()

	received := make(chan int, 10)
	subCh, err := b.Subscribe(0)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	defer b.Unsubscribe(subCh)

	go func() {
		for range subCh {
			time.Sleep(100 * time.Millisecond) // Simulate slow consumer
			received <- 1
		}
	}()

	b.Chan() <- 1
	b.Chan() <- 2

	// Allow enough time for timeout and message processing
	time.Sleep(200 * time.Millisecond)

	if len(received) != 1 {
		t.Errorf("Expected only one message to be processed due to timeout, got %d", len(received))
	}
}

func TestBroadcastChannelClose(t *testing.T) {
	b := New[int](10, 0)
	subCh, err := b.Subscribe(10)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	defer b.Unsubscribe(subCh)

	b.Close()

	select {
	case _, ok := <-subCh:
		if ok {
			t.Errorf("Expected subscriber channel to be closed but it was still open")
		}
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Expected subscriber channel to be closed immediately")
	}
}

func TestBroadcastSubscriberClosesChannel(t *testing.T) {
	b := New[int](10, 0)
	subCh, err := b.Subscribe(10)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	close(subCh)

	select {
	case _, ok := <-subCh:
		if ok {
			t.Errorf("Expected subscriber channel to be closed but it was still open")
		}
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Expected subscriber channel to be closed immediately")
	}
}

func TestSubscribeSubscriberClosed(t *testing.T) {
	b := New[int](10, 0)
	b.Close()

	subCh, err := b.Subscribe(10)
	if err != ErrBroadcasterClosed {
		t.Fatalf("expected ErrBroadcasterClosed, got %v", err)
	}

	defer b.Unsubscribe(subCh)
}

func TestConcurrentSubscriptions(t *testing.T) {
	b := New[int](10, 0)
	defer b.Close()

	const subCount = 50
	var wg sync.WaitGroup
	wg.Add(subCount)

	var subs []chan int

	for i := 0; i < subCount; i++ {
		go func() {
			subCh, err := b.Subscribe(10)
			if err != nil {
				t.Errorf("Failed to subscribe: %v", err)
			}

			subs = append(subs, subCh)

			wg.Done()
		}()
	}

	wg.Wait()

	b.Chan() <- 1

	// Give subscribers time to receive the message.
	time.Sleep(50 * time.Millisecond)

	for _, ch := range subs {
		if len(ch) != 1 {
			t.Errorf("Expected 1 message in channel, got %d", len(ch))
		}
	}

	// This is a simplistic check. Ideally, you should verify all subscribers receive the message.
	if len(b.subscribers) != subCount {
		t.Errorf("Expected %d subscribers, got %d", subCount, len(b.subscribers))
	}
}
