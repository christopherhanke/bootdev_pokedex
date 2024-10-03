package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	// fmt.Println("New cache was called")
	newCache := Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}
	go newCache.reapLoop(interval)
	return newCache
}

func (c Cache) Add(key string, val []byte) {
	// fmt.Printf("cache Add was called with key: %v, value %v\n", key, val)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	// fmt.Printf("cache Get was called with: %v\n", key)
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.cache[key]
	if !ok {
		return []byte{}, ok
	}
	return item.val, ok
}

func (c Cache) reapLoop(interval time.Duration) {
	// fmt.Println("reapLoop was called")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		// fmt.Println("reapLoop For begins")
		_, ok := <-ticker.C
		if !ok {
			break
		}
		// fmt.Println("reapLoop after Channel")
		func() {
			c.mu.Lock()
			defer c.mu.Unlock()
			for key, val := range c.cache {
				if time.Since(val.createdAt) > interval {
					fmt.Println("Deleting entry in cache.")
					delete(c.cache, key)
				}
			}
		}()
	}
}
