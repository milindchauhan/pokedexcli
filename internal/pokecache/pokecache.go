package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type PokeCache struct {
	cache    map[string]cacheEntry
	interval time.Duration
	mtx      sync.RWMutex
}

func NewCache(intervalInSeconds int) *PokeCache {
	cacheToBeReturned := PokeCache{
		cache:    map[string]cacheEntry{},
		interval: time.Second * time.Duration(intervalInSeconds),
		mtx:      sync.RWMutex{},
	}

	go cacheToBeReturned.reapLoop()

	return &cacheToBeReturned
}

func (p *PokeCache) Add(key string, val []byte) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (p *PokeCache) Get(key string) ([]byte, bool) {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	res, ok := p.cache[key]

	// leaving an intentional bug here for now. the res value can be nil and can cause the program to panic
	return res.val, ok
}

func (p *PokeCache) reapLoop() {
	for {
		p.mtx.Lock()
		for k, v := range p.cache {
			if time.Since(v.createdAt) > p.interval {
				delete(p.cache, k)
			}
		}
		p.mtx.Unlock()
		time.Sleep(p.interval)
	}
}
