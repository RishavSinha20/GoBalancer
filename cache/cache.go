package cache

import(
	"sync"
	"time"
)

type CacheItem struct{
	Data	[]byte
	ExpiresAt	time.Time
}

type Cache struct{
	store	map[string]CacheItem
	mu 	sync.RWMutex
	ttl	time.Duration
}

func NewCache(ttl time.Duration) *Cache{
	return &Cache{
		store : make(map[string]CacheItem),
		ttl : ttl,
	}
}

func (c *Cache) Get(key string)([]byte, bool){
	c.mu.RLock()
	item, exists := c.store[key]
	c.mu.RUnlock()

	if !exists || time.Now().After(item.ExpiresAt){
		return nil, false
	}

	return item.Data, true
}

func (c *Cache) Set(key string, data[] byte){
	c.mu.Lock()
	c.store[key] = CacheItem{
		Data : data, 
		ExpiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
}
