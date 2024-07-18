package tenor

import (
	"context"
	"sync"
	"time"

	"github.com/samber/lo"
)

// todo: он по хорошему должен быть многоуровневым
// но он пока одну гифку хранит так что пох
type gifCache struct {
	lock          *sync.Mutex
	m             map[GifId]gifCacheEntry
	cleanInterval time.Duration
}

type gifCacheEntry struct {
	createdAt time.Time
	ttm       *time.Duration
	content   []byte
}

func (e gifCacheEntry) isExpired() bool {
	if e.ttm == nil {
		return false
	}
	return time.Since(e.createdAt) > *e.ttm
}

func newGifCache(ctx context.Context, cleanInterval time.Duration) *gifCache {
	c := &gifCache{
		lock:          &sync.Mutex{},
		m:             map[GifId]gifCacheEntry{},
		cleanInterval: cleanInterval,
	}
	go c.cleanerWorker(ctx)
	return c
}

func (c *gifCache) cleanerWorker(ctx context.Context) {
	ticker := time.NewTicker(c.cleanInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		c.lock.Lock()
		for gifId, entry := range c.m {
			if entry.isExpired() {
				delete(c.m, gifId)
			}
		}
		c.lock.Unlock()
	}
}

func (c *gifCache) get(id GifId) ([]byte, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	result, ok := c.m[id]
	if !ok {
		return nil, false
	}
	if result.isExpired() {
		delete(c.m, id)
		return nil, false
	}
	return result.content, true
}

func (c *gifCache) set(id GifId, content []byte, ttm ...time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var ttmPtr *time.Duration
	if ttm, ok := lo.First(ttm); ok {
		ttmPtr = &ttm
	}
	c.m[id] = gifCacheEntry{
		createdAt: time.Now(),
		ttm:       ttmPtr,
		content:   content,
	}
}
