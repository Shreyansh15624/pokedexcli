package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (val []byte, ok bool) {
	c.mu.Lock()
	entry, ok := c.cache[key]
	c.mu.Unlock()
	if ok {
		return entry.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.cache {
			if time.Since(entry.createdAt) > interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{}
	c.cache = make(map[string]cacheEntry)

	go c.reapLoop(interval)

	return &c
}
