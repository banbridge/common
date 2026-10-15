package cachemanager

import (
	"context"
	"fmt"

	"github.com/banbridge/common/cache"
	"github.com/banbridge/common/cache/operation/local"
	localonredis "github.com/banbridge/common/cache/operation/local_on_redis"
	nocache "github.com/banbridge/common/cache/operation/no_cache"
	"github.com/redis/go-redis/v9"

	commonredis "github.com/banbridge/common/cache/operation/redis"
)

type cacheManager struct {
	redis        cache.RedisBaseCache
	local        cache.BaseCache
	localOnRedis cache.BaseCache
	noCache      cache.BaseCache
}

// NewCacheManager creates a CacheManager backed by Redis, bigcache local cache,
// a local-on-Redis two-level cache, and a no-cache implementation. The caller
// keeps ownership of the Redis client and local cache.
func NewCacheManager(redisCli *redis.Client, localCli cache.LocalStore, cfg cache.Config) (cache.CacheManager, error) {
	if redisCli == nil {
		return nil, fmt.Errorf("cache: redis client is nil")
	}
	if localCli == nil {
		return nil, fmt.Errorf("cache: local cache is nil")
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("cache: invalid config: %w", err)
	}
	return &cacheManager{
		redis:        commonredis.NewRedisCache(redisCli, cfg.RedisPrefix),
		local:        local.NewLocalCache(localCli, cfg.LocalPrefix),
		localOnRedis: localonredis.NewLocalOnRedis(redisCli, localCli, cfg.RedisPrefix),
		noCache:      nocache.NewNoCache(),
	}, nil
}

// GetRedis returns the Redis-backed base cache with increment support.
func (c *cacheManager) GetRedis(_ context.Context) cache.RedisBaseCache {
	return c.redis
}

// GetLocal returns the bigcache-backed local base cache.
func (c *cacheManager) GetLocal(_ context.Context) cache.BaseCache {
	return c.local
}

// GetLocalOnRedis returns the two-level base cache that reads local first and
// falls back to Redis.
func (c *cacheManager) GetLocalOnRedis(_ context.Context) cache.BaseCache {
	return c.localOnRedis
}

// GetNoCache returns the no-cache base cache implementation.
func (c *cacheManager) GetNoCache(_ context.Context) cache.BaseCache {
	return c.noCache
}

var _ cache.CacheManager = (*cacheManager)(nil)
