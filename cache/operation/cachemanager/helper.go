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

// maxAttempts bounds how often each cache item helper tries the underlying
// call: one try plus a single immediate retry on a non-miss error.
const maxAttempts = 2

// SaveCacheItem saves value into cacheHandle, retrying once on a save error.
func SaveCacheItem[T any](ctx context.Context, cacheHandle cache.Cache[T], key string, value T, expiration time.Duration, optFns ...cache.CacheOption) (err error) {
	if cacheHandle == nil {
		return fmt.Errorf("cache handle is nil")
	}
	for attempt := 1; ; attempt++ {
		err = cacheHandle.Save(ctx, key, value, expiration, optFns...)
		if err == nil || attempt >= maxAttempts {
			return err
		}
	}
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

// GetCacheItem gets a typed value from cacheHandle, retrying once on errors
// that are not cache misses. Cache misses are not retried and are normalized
// to cache.ErrCacheNotFound.
func GetCacheItem[T any](ctx context.Context, cacheHandle cache.Cache[T], key string, optFns ...cache.CacheOption) (res T, err error) {
	if cacheHandle == nil {
		return res, fmt.Errorf("cache handle is nil")
	}

	for attempt := 1; ; attempt++ {
		res, err = cacheHandle.Get(ctx, key, optFns...)
		if err == nil || attempt >= maxAttempts || isCacheMiss(err) {
			break
		}
	}

	err = normalizeCacheMiss(err)
	return
}

// DeleteCacheItem deletes key from cacheHandle, retrying once on errors that
// are not cache misses. Cache misses are treated as successful deletes.
func DeleteCacheItem[T any](ctx context.Context, cacheHandle cache.Cache[T], key string, optFns ...cache.CacheOption) (err error) {
	if cacheHandle == nil {
		return fmt.Errorf("cache handle is nil")
	}

	for attempt := 1; ; attempt++ {
		err = cacheHandle.Delete(ctx, key, optFns...)
		if err == nil || attempt >= maxAttempts || isCacheMiss(err) {
			break
		}
	}

	if err != nil && !isCacheMiss(err) {
		return err
	}
	return nil
}
