// Package observe provides a lightweight event bus for broadcasting
// internal lifecycle events (secret fetched, lease renewed, cache hit, etc.)
// to registered subscribers without tight coupling between components.
package observe

import "sync"

// EventKind identifies the type of event.
type EventKind string

const (
	SecretFetched  EventKind = "secret.fetched"
	SecretRotated  EventKind = "secret.rotated"
	LeaseRenewed   EventKind = "lease.renewed"
	CacheHit       EventKind = "cache.hit"
	CacheMiss      EventKind = "cache.miss"
	ProcessStarted EventKind = "process.started"
	ProcessExited  EventKind = "process.exited"
)

// Event carries a kind and an optional string payload (e.g. secret path).
type Event struct {
	Kind    EventKind
	Payload string
}

// Handler is a callback invoked when a matching event is published.
type Handler func(Event)

// Bus is a thread-safe publish/subscribe event bus.
type Bus struct {
	mu       sync.RWMutex
	handlers map[EventKind][]Handler
}

// New returns an initialised, empty Bus.
func New() *Bus {
	return &Bus{
		handlers: make(map[EventKind][]Handler),
	}
}

// Subscribe registers h to be called whenever an event of kind k is published.
// The same handler may be registered multiple times.
func (b *Bus) Subscribe(k EventKind, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[k] = append(b.handlers[k], h)
}

// Publish delivers e to all handlers registered for e.Kind.
// Handlers are called synchronously in registration order.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	list := make([]Handler, len(b.handlers[e.Kind]))
	copy(list, b.handlers[e.Kind])
	b.mu.RUnlock()

	for _, h := range list {
		h(e)
	}
}

// Reset removes all registered handlers.
func (b *Bus) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = make(map[EventKind][]Handler)
}
