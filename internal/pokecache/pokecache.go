package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval int) {
	fmt.Println("New cache was called")

}

func (c cache) Add(key string, val []byte) {
	fmt.Println("cache Add was called")
}

func (c cache) Get(key string) ([]byte, bool) {
	fmt.Printf("cache Get was called with: %v\n", key)
	return []byte{}, true
}

func (c cache) reapLoop() {
	fmt.Println("reapLoop was called")
}
