package memory

import (
	"context"
	"sync"
	"time"

	"github.com/ave1995/grpc-chat/domain/model"
	"github.com/ave1995/grpc-chat/domain/store"
)

var _ store.Cache = (*Cache)(nil)

type Cache struct {
	mu       sync.RWMutex
	items    map[string]model.CacheItem
	stopChan chan struct{}
}

func NewCache(cleanupInterval time.Duration) *Cache {
	mc := &Cache{
		items:    make(map[string]model.CacheItem),
		stopChan: make(chan struct{}),
	}

	if cleanupInterval > 0 {
		go mc.cleanupLoop(cleanupInterval)
	}

	return mc
}

// Set stores a value with an optional TTL (0 = no expiration)
func (mc *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	mc.mu.Lock()
	mc.items[key] = model.CacheItem{Value: value, Expiration: exp}
	mc.mu.Unlock()
	return nil
}

// Get retrieves a value if present and not expired
func (mc *Cache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	select {
	case <-ctx.Done():
		return nil, false, ctx.Err()
	default:
	}

	mc.mu.RLock()
	item, found := mc.items[key]
	mc.mu.RUnlock()

	if !found {
		return nil, false, nil
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		_ = mc.Delete(ctx, key)
		return nil, false, nil
	}

	return item.Value, true, nil
}

// Delete removes a key
func (mc *Cache) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	mc.mu.Lock()
	delete(mc.items, key)
	mc.mu.Unlock()
	return nil
}

// Exists checks if a key exists (and is not expired)
func (mc *Cache) Exists(ctx context.Context, key string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}
	_, ok, _ := mc.Get(ctx, key)
	return ok, nil
}

// Keys returns all current keys
func (mc *Cache) Keys(ctx context.Context) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	mc.mu.RLock()
	keys := make([]string, 0, len(mc.items))
	for k := range mc.items {
		keys = append(keys, k)
	}
	mc.mu.RUnlock()
	return keys, nil
}

// Clear removes all items
func (mc *Cache) Clear(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	mc.mu.Lock()
	mc.items = make(map[string]model.CacheItem)
	mc.mu.Unlock()
	return nil
}

// Close stops background cleanup
func (mc *Cache) Close() error {
	close(mc.stopChan)
	return nil
}

// Internal cleanup loop
func (mc *Cache) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.deleteExpired()
		case <-mc.stopChan:
			return
		}
	}
}

// deleteExpired removes all expired keys
func (mc *Cache) deleteExpired() {
	now := time.Now().UnixNano()
	mc.mu.Lock()
	for k, v := range mc.items {
		if v.Expiration > 0 && now > v.Expiration {
			delete(mc.items, k)
		}
	}
	mc.mu.Unlock()
}
