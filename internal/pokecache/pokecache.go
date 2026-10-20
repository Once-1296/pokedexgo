package pokecache

import (
	"sync"
	"time"
)
var mu sync.Mutex
type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	cacheMap map[string]cacheEntry
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{make(map[string]cacheEntry)}
	go cache.reapLoop(interval)
	return cache
}

func (cache Cache) Add(key string, val []byte) {
	mu.Lock()
	defer mu.Unlock()
	cache.cacheMap[key] = cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
}

func (cache Cache) Get(key string) ([]byte, bool) {
	mu.Lock()
	defer mu.Unlock()
	entry, ok := cache.cacheMap[key]
	if !ok {
		return nil, ok
	}
	return entry.val, ok
}

func (cache Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		mu.Lock()
		for k,v := range cache.cacheMap {
			if v.createdAt.Add(interval).Before(time.Now()) {
				delete(cache.cacheMap,k)
			}
		}
		mu.Unlock()
	}
}