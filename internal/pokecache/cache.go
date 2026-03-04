package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	value     []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		mu:      sync.Mutex{},
		entries: make(map[string]cacheEntry),
	}
	go cache.readLoop(interval)
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		value:     val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.entries[key]; !ok {
		return nil, false
	} else {
		return v.value, true
	}
}

func (c *Cache) readLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, val := range c.entries {
			if time.Since(val.createdAt) > interval {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}

}
