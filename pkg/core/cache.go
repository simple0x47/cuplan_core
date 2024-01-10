package core

import (
	"sync"
	"time"
)

// CacheItem provides a storage facility for a value which is set to expire.
type CacheItem struct {
	value    any
	expireAt time.Time
}

// Cache provides a storage facility for caching data for a specific duration.
type Cache struct {
	cache       *sync.Map
	stopChannel chan struct{}
}

// NewCacheItem creates an instance of CacheItem which is set to expire after the specified duration.
func NewCacheItem(value any, duration time.Duration) *CacheItem {
	item := new(CacheItem)

	item.value = value
	item.expireAt = time.Now().Add(duration)

	return item
}

// NewCache creates an instance of Cache which cleans expired cache items, every time, after the specified duration.
func NewCache(cleanAfter time.Duration) *Cache {
	cache := new(Cache)

	cache.cache = new(sync.Map)
	cache.stopChannel = make(chan struct{})

	go cache.cleanEveryTime(cleanAfter, cache.stopChannel)

	return cache
}

// Close stops the cleaning go routine.
func (c *Cache) Close() {
	close(c.stopChannel)
}

// Set stores the specified value at the specified key for the specified duration.
func (c *Cache) Set(key any, value any, duration time.Duration) {
	item := NewCacheItem(value, duration)

	c.cache.Store(key, item)
}

// Get returns the cached value or None if there is not a valid cache hit.
func (c *Cache) Get(key any) Option[any] {
	obj, ok := c.cache.Load(key)

	if !ok {
		return None
	}

	item, ok := obj.(*CacheItem)

	if !ok {
		return None
	}

	return Some(item.value)
}

// HasKey returns whether the specific key has been cached already.
func (c *Cache) HasKey(key any) bool {
	_, ok := c.cache.Load(key)

	return ok
}

// Remove eliminates the value for a key.
func (c *Cache) Remove(key any) {
	c.cache.Delete(key)
}

// Clear clears the cache completely.
func (c *Cache) Clear() {
	c.cache = new(sync.Map)
}

func (c *Cache) cleanEveryTime(cleanAfter time.Duration, stopChannel chan struct{}) {
	for {
		select {
		case <-stopChannel:
			return
		default:
			time.Sleep(cleanAfter)
			keysToBeRemoved := make([]any, 0, 10)

			c.cache.Range(func(key, value any) bool {
				item, ok := value.(*CacheItem)

				if !ok {
					return true
				}

				if time.Now().After(item.expireAt) {
					keysToBeRemoved = append(keysToBeRemoved, key)
				}

				return true
			})

			for _, key := range keysToBeRemoved {
				c.cache.Delete(key)
			}
		}
	}
}
