package cache

import (
	"context"
	"sync"
	"time"
)

// Opts - options to create new cache instance
// CleanInterval - uses if value > 0
type Opts struct {
	CleanInterval time.Duration
	Logger        Logger
}

// New - create new expirable cache instance with clean interval in ms
func New[K comparable, V any](ctx context.Context, opts Opts) *Cache[K, V] {
	c := &Cache[K, V]{
		cache:  make(map[K]value[V]),
		logger: opts.Logger,
	}
	if opts.CleanInterval > 0 {
		go c.runCleaner(ctx, opts.CleanInterval)
	}

	return c
}

// Cache - generic key-value cache with time expiration
type Cache[K comparable, V any] struct {
	cache  map[K]value[V]
	mu     sync.RWMutex
	logger Logger
}

// AddWithExp - add value by key with time expiration
func (c *Cache[K, V]) AddWithExp(k K, v V, exp time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[k] = value[V]{
		data: v,
		exp:  exp.UnixNano(),
	}
}

// Add - add value by key
func (c *Cache[K, V]) Add(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[k] = value[V]{
		data: v,
	}
}

// Delete - delete value by key
func (c *Cache[K, V]) Delete(k K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, k)
}

// Add - get actual value by key
func (c *Cache[K, V]) Get(k K) (v V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, ok := c.cache[k]
	if ok && value.actual() {
		v, ok = value.data, true
		return
	}

	return v, false
}

// Keys - get cache keys
func (c *Cache[K, V]) Keys() (keys []K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.cache {
		if v.actual() {
			keys = append(keys, k)
		}
	}

	return
}

// Keys - get cache values
func (c *Cache[K, V]) Values() (values []V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, v := range c.cache {
		if v.actual() {
			values = append(values, v.data)
		}
	}

	return
}

// ClearExpired - clear expired keys
func (c *Cache[K, V]) ClearExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.cache {
		if !v.actual() {
			delete(c.cache, k)
		}
	}
}

func (c *Cache[K, V]) runCleaner(ctx context.Context, interval time.Duration) {
	const op = "cache.runCleaner"

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug("context canceled", "op", op)
			return
		case <-ticker.C:
			c.logger.Debug("clean cache", "op", op)
			c.ClearExpired()
		}
	}
}

type value[V any] struct {
	data V
	exp  int64
}

func (v value[V]) actual() bool {
	return v.exp == 0 || time.Now().UnixNano() < v.exp
}
