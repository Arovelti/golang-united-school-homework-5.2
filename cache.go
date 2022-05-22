package cache

import (
	"sync"
	"time"
)

type CacheProperties struct {
	value    string
	deadline time.Time
	elapse   bool
}

type Cache struct {
	mu       sync.Mutex
	cacheMap map[string]CacheProperties
}

func NewCache() Cache {
	cache := make(map[string]CacheProperties)
	return Cache{cacheMap: cache}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if k, ok := c.cacheMap[key]; ok {
		k.elapse = k.deadline.Before(time.Now())
		if !k.elapse || k.deadline.IsZero() {
			return k.value, true
		}
	}
	return "", false
}

func (c *Cache) Put(key, value string) {
	c.cacheMap[key] = CacheProperties{
		value:    value,
		deadline: time.Time{},
		elapse:   false,
	}
}

func (c *Cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]string, 0, len(c.cacheMap))
	for key, val := range c.cacheMap {
		val.elapse = val.deadline.Before(time.Now())
		if !val.elapse || val.deadline.IsZero() {
			keys = append(keys, key)
		}
	}
	return keys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheMap[key] = CacheProperties{
		value:    value,
		deadline: deadline,
		elapse:   false,
	}
}
