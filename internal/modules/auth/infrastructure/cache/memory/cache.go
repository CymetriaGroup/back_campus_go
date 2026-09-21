package memory

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type item struct {
	value   []byte
	expires time.Time
}
type Cache struct {
	mu    sync.RWMutex
	items map[string]item
}

func New() *Cache { return &Cache{items: make(map[string]item)} }
func (c *Cache) Get(_ context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	v, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(v.expires) {
		if ok {
			_ = c.Delete(context.Background(), key)
		}
		return nil, ErrCacheMiss
	}
	return append([]byte(nil), v.value...), nil
}
func (c *Cache) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = item{append([]byte(nil), value...), time.Now().Add(ttl)}
	return nil
}
func (c *Cache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}
