// Package cache wraps go-redis to provide a simple string key-value cache
// used in the service layer. All values are stored as strings — callers are
// responsible for serializing (e.g. JSON) before storing and deserializing
// after retrieval.
package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache holds a go-redis client and is the only exported type in this
// package. Callers interact through Get, Set, and Delete — the underlying
// client is never exposed directly, keeping the Redis dependency internal.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a Redis client pointing at the given address.
// No connection is made here — go-redis connects lazily on first command.
// addr should be in "host:port" format, e.g. "localhost:6379".
func NewRedisCache(addr string) *RedisCache {

	var redCache RedisCache
	newClient := redis.NewClient(&redis.Options{Addr: addr})

	redCache.client = newClient
	return &redCache
}

// Get returns the value stored under key. Returns an error if the key does
// not exist (redis.Nil) or if Redis is unreachable. Callers should treat any
// error as a cache miss and fall back to the source of truth (e.g. the DB).
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return val, err
	}

	return val, nil

}

// Set stores value under key with the given TTL. After the TTL expires, Redis
// automatically evicts the key — no manual cleanup is needed. A TTL of 0
// means the key never expires, which should be avoided for user data.
func (r *RedisCache) Set(ctx context.Context, key string, val string, ttl time.Duration) error {

	err := r.client.Set(ctx, key, val, ttl).Err()
	return err
}

// Delete removes key from the cache. Used for cache invalidation after writes
// so subsequent reads don't serve stale data. A failed delete is not fatal —
// the TTL will eventually evict the stale entry on its own.
func (r *RedisCache) Delete(ctx context.Context, key string) error {

	err := r.client.Del(ctx, key).Err()
	return err
}
