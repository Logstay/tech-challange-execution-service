package cache

import (
	"sync"
	"time"

	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/cache"
)

type CacheItem struct {
	Data      any
	ExpiresAt time.Time
}

type MemoryCache struct {
	items map[string]CacheItem
	mutex sync.RWMutex
}

func NewMemoryCache() ports.Cache {
	cache := &MemoryCache{
		items: make(map[string]CacheItem),
	}

	go cache.cleanupExpired()
	return cache
}

func (c *MemoryCache) Set(key string, value any, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(duration),
	}
}

func (c *MemoryCache) Get(key string) (any, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists || time.Now().After(item.ExpiresAt) {
		return nil, false
	}

	return item.Data, true
}

func (c *MemoryCache) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items = make(map[string]CacheItem)
}

func (c *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}

var _ ports.Cache = (*MemoryCache)(nil)
