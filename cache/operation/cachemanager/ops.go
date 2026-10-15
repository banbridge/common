// Package cachemanager provides typed cache operations and the concrete
// CacheManager implementation.
package cachemanager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"

	"github.com/banbridge/common/cache"
)

type cacheOperation[T any] struct {
	cacheManager cache.CacheManager
}

// ErrInvalidCacheType indicates that a cache operation was configured with an
// unsupported CacheType.
var ErrInvalidCacheType = errors.New("invalid cache type")

// NewCacheOperation creates a typed cache operation wrapper over manager.
func NewCacheOperation[T any](manager cache.CacheManager) (cache.Cache[T], error) {
	if manager == nil {
		return nil, fmt.Errorf("cache manager is nil")
	}
	return &cacheOperation[T]{cacheManager: manager}, nil
}

// Get reads key from the selected base cache and unmarshals the JSON payload
// into T.
func (c *cacheOperation[T]) Get(ctx context.Context, key string, optFns ...cache.CacheOption) (T, error) {
	var res T
	base, err := c.getBaseCache(ctx, optFns...)
	if err != nil {
		return res, err
	}
	data, err := base.Get(ctx, key)
	if err != nil {
		return res, err
	}
	err = sonic.Unmarshal(data, &res)
	return res, err
}

// Save marshals entity as JSON and writes it to the selected base cache.
func (c *cacheOperation[T]) Save(ctx context.Context, key string, entity T, expire time.Duration, optFns ...cache.CacheOption) error {
	data, err := sonic.Marshal(entity)
	if err != nil {
		return err
	}
	base, err := c.getBaseCache(ctx, optFns...)
	if err != nil {
		return err
	}
	return base.Save(ctx, key, data, expire)
}

// Delete removes key from the selected base cache.
func (c *cacheOperation[T]) Delete(ctx context.Context, key string, optFns ...cache.CacheOption) error {
	base, err := c.getBaseCache(ctx, optFns...)
	if err != nil {
		return err
	}
	return base.Delete(ctx, key)
}

// GetCacheManager returns the CacheManager used by this operation wrapper.
func (c *cacheOperation[T]) GetCacheManager() cache.CacheManager {
	return c.cacheManager
}

func (c *cacheOperation[T]) getBaseCache(ctx context.Context, optFns ...cache.CacheOption) (cache.BaseCache, error) {
	option := cache.MergeCacheOptions(optFns...)

	switch option.CacheType {
	case cache.RedisCache:
		return c.cacheManager.GetRedis(ctx), nil
	case cache.LocalCache:
		return c.cacheManager.GetLocal(ctx), nil
	case cache.LocalOnRedisLevel:
		return c.cacheManager.GetLocalOnRedis(ctx), nil
	case cache.CustomCache:
		if option.CustomBaseCache == nil {
			return nil, fmt.Errorf("%w: custom base cache is nil", ErrInvalidCacheType)
		}
		return option.CustomBaseCache, nil
	case cache.NoCache:
		return c.cacheManager.GetNoCache(ctx), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrInvalidCacheType, option.CacheType)
	}
}
