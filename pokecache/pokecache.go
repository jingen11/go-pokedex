package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mu       *sync.RWMutex
}

func NewCache(duration time.Duration) Cache {
	cache := Cache{
		interval: duration,
		entries:  map[string]cacheEntry{},
		mu:       &sync.RWMutex{},
	}
	ticker := time.NewTicker(cache.interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				cache.reapLoop()
			}
		}
	}()

	return cache
}

func (c Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.entries[key]

	if ok {
		return val.val, ok
	}
	return nil, ok
}

func (c Cache) reapLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.entries {
		life := time.Now().Sub(v.createdAt)
		if life > c.interval {
			delete(c.entries, k)
		}
	}
}
