package store

import (
	"context"
	"time"
)

type Cache interface {
	// Set stores a value with an optional TTL (0 = no expiration)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	// Get retrieves a value by key.
	// Returns (value, true) if found and not expired, otherwise (nil, false)
	Get(ctx context.Context, key string) (any, bool, error)

	// Delete removes a key.
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists (and is not expired)
	Exists(ctx context.Context, key string) (bool, error)

	// Keys returns a list of all keys (use carefully; may be expensive)
	Keys(ctx context.Context) ([]string, error)

	// Clear removes all keys from the cache
	Clear(ctx context.Context) error

	// Close gracefully shuts down the cache (e.g., closes Redis connection)
	Close(ctx context.Context) error
}
