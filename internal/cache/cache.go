package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// TODO(ppacher): add a redis implementation for the cache.

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Cache interface {
	PutKey(ctx context.Context, key string, value any) error
	PutKeyTTL(ctx context.Context, key string, value any, ttl time.Duration) error
	GetKey(ctx context.Context, key string, receiver any) error
	GetAndDeleteKey(ctx context.Context, key string, receiver any) error
	DeleteKey(ctx context.Context, key string) error
}

func NewInMemoryCache() Cache {
	return new(inMemoryCache)
}

type cacheEntry struct {
	value   []byte
	expires time.Time
}

type inMemoryCache struct {
	l sync.RWMutex

	keys map[string]cacheEntry
}

func (c *inMemoryCache) PutKey(_ context.Context, key string, value any) error {
	c.l.Lock()
	defer c.l.Unlock()

	blob, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.keys[key] = cacheEntry{value: blob}

	return nil
}

func (c *inMemoryCache) PutKeyTTL(_ context.Context, key string, value any, ttl time.Duration) error {
	c.l.Lock()
	defer c.l.Unlock()

	blob, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.keys[key] = cacheEntry{
		value:   blob,
		expires: time.Now().Add(ttl),
	}

	return nil
}

func (c *inMemoryCache) GetKey(_ context.Context, key string, receiver any) error {
	c.l.RLock()
	defer c.l.RUnlock()

	val, ok := c.keys[key]
	if !ok {
		return ErrKeyNotFound
	}

	if !val.expires.IsZero() && val.expires.Before(time.Now()) {
		return ErrKeyNotFound
	}

	if receiver == nil {
		return nil
	}

	return json.Unmarshal(val.value, receiver)
}

func (c *inMemoryCache) GetAndDeleteKey(_ context.Context, key string, receiver any) error {
	c.l.Lock()
	defer c.l.Unlock()

	val, ok := c.keys[key]
	if !ok {
		return ErrKeyNotFound
	}

	delete(c.keys, key)

	if !val.expires.IsZero() && val.expires.Before(time.Now()) {
		return ErrKeyNotFound
	}

	if receiver == nil {
		return nil
	}

	return json.Unmarshal(val.value, receiver)
}

func (c *inMemoryCache) DeleteKey(ctx context.Context, key string) error {
	return c.GetAndDeleteKey(ctx, key, nil)
}
