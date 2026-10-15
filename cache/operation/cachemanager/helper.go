package cachemanager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/banbridge/common/cache"
	"github.com/redis/go-redis/v9"
)

// SaveCacheItem saves value into cacheHandle with retry on any save error.
func SaveCacheItem[T any](ctx context.Context, cacheHandle cache.Cache[T], key string, value T, expiration time.Duration, optFns ...cache.CacheOption) (err error) {
	if cacheHandle == nil {
		return fmt.Errorf("cache handle is nil")
	}
	return cacheHandle.Save(ctx, key, value, expiration, optFns...)
}

func isCacheMiss(err error) bool {
	return errors.Is(err, cache.ErrCacheNotFound) ||
		errors.Is(err, cache.ErrNoCache) ||
		errors.Is(err, redis.Nil) ||
		errors.Is(err, bigcache.ErrEntryNotFound)
}

func normalizeCacheMiss(err error) error {
	if err == nil || !isCacheMiss(err) {
		return err
	}
	if errors.Is(err, cache.ErrCacheNotFound) {
		return err
	}
	return fmt.Errorf("%w: %w", cache.ErrCacheNotFound, err)
}

// GetCacheItem gets a typed value from cacheHandle with retry on transient
// errors. Cache misses are not retried and are normalized to
// cache.ErrCacheNotFound.
func GetCacheItem[T any](ctx context.Context, cacheHandle cache.Cache[T], key string, optFns ...cache.CacheOption) (res T, err error) {
	if cacheHandle == nil {
		return res, fmt.Errorf("cache handle is nil")
	}

	res, err = cacheHandle.Get(ctx, key, optFns...)

	err = normalizeCacheMiss(err)
	return
}

// DeleteCacheItem deletes key from cacheHandle with retry on transient errors.
// Cache misses are treated as successful deletes.
func DeleteCacheItem[T any](ctx context.Context, cacheHandle cache.Cache[T], key string, optFns ...cache.CacheOption) (err error) {
	if cacheHandle == nil {
		return fmt.Errorf("cache handle is nil")
	}

	err = cacheHandle.Delete(ctx, key, optFns...)
	
	if err != nil && !isCacheMiss(err) {
		return err
	}
	return nil
}
