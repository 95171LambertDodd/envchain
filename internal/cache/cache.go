// Package cache provides a simple in-memory key-value cache with optional
// TTL and size-bounded eviction for resolved environment entries.
package cache

import (
	"errors"
	"sync"
	"time"
)

// ErrCacheMiss is returned when a key is not found in the cache.
var ErrCacheMiss = errors.New("cache: key not found")

// ErrInvalidCapacity is returned when capacity is less than 1.
var ErrInvalidCapacity = errors.New("cache: capacity must be >= 1")

// entry holds a cached value and its expiry time.
type entry struct {
	value   string
	expiry  time.Time
	hasExpiry bool
}

// Cache is a thread-safe in-memory store for string key-value pairs.
type Cache struct {
	mu       sync.RWMutex
	items    map[string]entry
	capacity int
	clock    func() time.Time
}

// New creates a Cache with the given maximum capacity.
func New(capacity int) (*Cache, error) {
	if capacity < 1 {
		return nil, ErrInvalidCapacity
	}
	return &Cache{
		items:    make(map[string]entry, capacity),
		capacity: capacity,
		clock:    time.Now,
	}, nil
}

// Set stores a value under key with no expiry.
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictIfFull(key)
	c.items[key] = entry{value: value}
}

// SetWithTTL stores a value under key that expires after ttl duration.
func (c *Cache) SetWithTTL(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictIfFull(key)
	c.items[key] = entry{value: value, expiry: c.clock().Add(ttl), hasExpiry: true}
}

// Get retrieves a value by key. Returns ErrCacheMiss if absent or expired.
func (c *Cache) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[key]
	if !ok {
		return "", ErrCacheMiss
	}
	if e.hasExpiry && c.clock().After(e.expiry) {
		return "", ErrCacheMiss
	}
	return e.value, nil
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Len returns the number of items currently in the cache (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// evictIfFull removes an arbitrary entry when at capacity and key is new.
// Must be called with the write lock held.
func (c *Cache) evictIfFull(newKey string) {
	if _, exists := c.items[newKey]; exists {
		return
	}
	if len(c.items) >= c.capacity {
		for k := range c.items {
			delete(c.items, k)
			break
		}
	}
}
