package lru

import (
	"sync"
	"time"
)

type node[K comparable, V any] struct {
	next *node[K, V]
	prev *node[K, V]

	added time.Time
	key   K
	value V
}

type LRUCache[K comparable, V any] struct {
	cache   map[K]*node[K, V]
	mu      sync.RWMutex
	maxSize int
	len     int
	ttl     time.Duration

	head *node[K, V]
	tail *node[K, V]
}

func NewLRUCache[K comparable, V any](maxSize int, ttl time.Duration) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		maxSize: maxSize,
		ttl:     ttl,
		cache:   make(map[K]*node[K, V], maxSize),
	}
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	n, ok := c.cache[key]
	if !ok || c.ttl < time.Since(n.added) {
		var ret V

		return ret, false
	}

	return n.value, ok
}

func (c *LRUCache[K, V]) Put(key K, value V) (removed []K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.cache[key]; ok {
		node.value = value
		node.added = time.Now()

		return
	}

	if c.len == c.maxSize {
		for c.ttl < time.Since(c.head.added) {
			removed = append(removed, c.head.key)
			c.remove(c.head)
		}
	}

	newNode := &node[K, V]{
		value: value,
		added: time.Now(),
		prev:  c.tail,
		key:   key,
	}

	if c.tail == nil {
		c.head = newNode
		c.tail = newNode
	} else {
		c.tail.next = newNode
	}

	c.tail = newNode

	c.len++
	c.cache[key] = newNode

	return
}

// no concurrent safe
func (c *LRUCache[K, V]) remove(n *node[K, V]) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}

	c.len--
	delete(c.cache, n.key)
}
