package redis

import (
	"context"
	"errors"
	"time"

	"github.com/ave1995/grpc-chat/domain/store"
	"github.com/redis/go-redis/v9"
)

var _ store.Cache = (*Cache)(nil)

type Cache struct {
	client *redis.Client
}

// NewCache creates a new Redis-backed cache
func NewCache(addr string, password string, db int) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Cache{client: rdb}
}

// Set stores a value with optional TTL
func (c *Cache) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value by key
func (c *Cache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return data, true, nil
}

// Delete removes a key
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Keys returns all keys (⚠️ be careful in production)
func (c *Cache) Keys(ctx context.Context) ([]string, error) {
	keys, err := c.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// Clear removes all keys from the cache
func (c *Cache) Clear(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}
