// Package cache defines deployment-independent cache contracts and the
// building blocks for typed, multi-backend cache operations.
package cache

import (
	"context"
	"time"
)

// Cache is a typed cache operation interface. Backend selection and expiration
// are controlled by CacheOption values.
type Cache[T any] interface {
	Get(ctx context.Context, key string, optFns ...CacheOption) (T, error)
	Save(ctx context.Context, key string, entity T, expire time.Duration, optFns ...CacheOption) error
	Delete(ctx context.Context, key string, optFns ...CacheOption) error
	GetCacheManager() CacheManager
}

// BaseCache is a byte-oriented cache backend operation interface.
type BaseCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Save(ctx context.Context, key string, value []byte, expire time.Duration) error
	Delete(ctx context.Context, key string) error
}

// Incrementer is implemented by backends that support atomic increment with
// expiration, currently the Redis backend.
type Incrementer interface {
	IncrByExpire(ctx context.Context, key string, value int64, expiration time.Duration) (int64, error)
}

// RedisBaseCache is the Redis backend surface: generic byte operations plus
// atomic increment with expiration.
type RedisBaseCache interface {
	BaseCache
	Incrementer
}

// CacheManager provides access to the available cache backends.
type CacheManager interface {
	GetRedis(ctx context.Context) RedisBaseCache
	GetLocal(ctx context.Context) BaseCache
	GetLocalOnRedis(ctx context.Context) BaseCache
	GetNoCache(ctx context.Context) BaseCache
}

// LocalStore is the minimal local cache surface used by the cache component.
// bigcache.BigCache satisfies it directly.
type LocalStore interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

// Config contains backend-independent cache configuration.
type Config struct {
	// RedisPrefix is prepended to every Redis cache key.
	RedisPrefix string
	// LocalPrefix is prepended to every standalone local cache key.
	LocalPrefix string
}

// DefaultConfig returns cache configuration with no key prefixes.
func DefaultConfig() Config {
	return Config{}
}

// Validate checks cache configuration.
func (c Config) Validate() error {
	return nil
}
