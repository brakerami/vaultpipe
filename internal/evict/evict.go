// Package evict provides a least-recently-used (LRU) eviction policy
// for bounded secret caches. When the cache reaches capacity the oldest
// entry is removed to make room for the incoming one.
package evict

import (
	"container/list"
	"sync"
)

// entry is stored in both the map and the linked list.
type entry struct {
	key   string
	value string
}

// LRU is a thread-safe, fixed-capacity LRU cache.
type LRU struct {
	mu       sync.Mutex
	cap      int
	items    map[string]*list.Element
	order    *list.List
	Evicted  func(key, value string) // optional hook called on eviction
}

// New returns an LRU with the given capacity. Panics if cap < 1.
func New(cap int) *LRU {
	if cap < 1 {
		panic("evict: capacity must be at least 1")
	}
	return &LRU{
		cap:   cap,
		items: make(map[string]*list.Element, cap),
		order: list.New(),
	}
}

// Set inserts or updates key. If the cache is at capacity the least-recently
// used entry is evicted first.
func (l *LRU) Set(key, value string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if el, ok := l.items[key]; ok {
		el.Value.(*entry).value = value
		l.order.MoveToFront(el)
		return
	}

	if l.order.Len() >= l.cap {
		l.evictOldest()
	}

	e := &entry{key: key, value: value}
	el := l.order.PushFront(e)
	l.items[key] = el
}

// Get returns the value for key and marks it as recently used.
// The second return value is false when the key is absent.
func (l *LRU) Get(key string) (string, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	el, ok := l.items[key]
	if !ok {
		return "", false
	}
	l.order.MoveToFront(el)
	return el.Value.(*entry).value, true
}

// Remove deletes key from the cache. It is a no-op if the key is absent.
func (l *LRU) Remove(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if el, ok := l.items[key]; ok {
		l.removeElement(el)
	}
}

// Len returns the number of entries currently held.
func (l *LRU) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.order.Len()
}

func (l *LRU) evictOldest() {
	el := l.order.Back()
	if el != nil {
		l.removeElement(el)
	}
}

func (l *LRU) removeElement(el *list.Element) {
	e := el.Value.(*entry)
	l.order.Remove(el)
	delete(l.items, e.key)
	if l.Evicted != nil {
		l.Evicted(e.key, e.value)
	}
}
