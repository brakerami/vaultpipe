package observe_test

import (
	"sync"
	"testing"

	"github.com/yourusername/vaultpipe/internal/observe"
)

func TestPublish_NoHandlers_DoesNotPanic(t *testing.T) {
	b := observe.New()
	// should not panic
	b.Publish(observe.Event{Kind: observe.SecretFetched, Payload: "secret/foo"})
}

func TestSubscribe_ReceivesMatchingEvent(t *testing.T) {
	b := observe.New()
	var got observe.Event
	b.Subscribe(observe.SecretFetched, func(e observe.Event) { got = e })

	want := observe.Event{Kind: observe.SecretFetched, Payload: "secret/db"}
	b.Publish(want)

	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestSubscribe_DoesNotReceiveOtherKind(t *testing.T) {
	b := observe.New()
	called := false
	b.Subscribe(observe.CacheHit, func(e observe.Event) { called = true })

	b.Publish(observe.Event{Kind: observe.CacheMiss})

	if called {
		t.Fatal("handler should not have been called for a different kind")
	}
}

func TestSubscribe_MultipleHandlers_AllCalled(t *testing.T) {
	b := observe.New()
	count := 0
	for i := 0; i < 3; i++ {
		b.Subscribe(observe.LeaseRenewed, func(e observe.Event) { count++ })
	}
	b.Publish(observe.Event{Kind: observe.LeaseRenewed})
	if count != 3 {
		t.Fatalf("expected 3 handler calls, got %d", count)
	}
}

func TestReset_RemovesAllHandlers(t *testing.T) {
	b := observe.New()
	called := false
	b.Subscribe(observe.SecretRotated, func(e observe.Event) { called = true })
	b.Reset()
	b.Publish(observe.Event{Kind: observe.SecretRotated})
	if called {
		t.Fatal("handler should not be called after Reset")
	}
}

func TestPublish_ConcurrentSafe(t *testing.T) {
	b := observe.New()
	var mu sync.Mutex
	var events []observe.Event
	b.Subscribe(observe.ProcessExited, func(e observe.Event) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	})

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Publish(observe.Event{Kind: observe.ProcessExited})
		}()
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if len(events) != 50 {
		t.Fatalf("expected 50 events, got %d", len(events))
	}
}
