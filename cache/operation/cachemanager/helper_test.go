package cachemanager

import (
	"context"
	"errors"
	"testing"
	"time"

	bigcache "github.com/allegro/bigcache/v3"
	"github.com/stretchr/testify/require"

	"github.com/banbridge/common/cache"
)

type fakeGenericCache[T any] struct {
	manager cache.CacheManager

	getFn    func(context.Context, string, ...cache.CacheOption) (T, error)
	saveFn   func(context.Context, string, T, time.Duration, ...cache.CacheOption) error
	deleteFn func(context.Context, string, ...cache.CacheOption) error
}

func (f *fakeGenericCache[T]) Get(ctx context.Context, key string, optFns ...cache.CacheOption) (T, error) {
	if f.getFn != nil {
		return f.getFn(ctx, key, optFns...)
	}
	var zero T
	return zero, nil
}

func (f *fakeGenericCache[T]) Save(ctx context.Context, key string, entity T, expire time.Duration, optFns ...cache.CacheOption) error {
	if f.saveFn != nil {
		return f.saveFn(ctx, key, entity, expire, optFns...)
	}
	return nil
}

func (f *fakeGenericCache[T]) Delete(ctx context.Context, key string, optFns ...cache.CacheOption) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, key, optFns...)
	}
	return nil
}

func (f *fakeGenericCache[T]) GetCacheManager() cache.CacheManager {
	return f.manager
}

func TestSaveCacheItemRetries(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("temporary save error")
	attempts := 0
	handle := &fakeGenericCache[testPayload]{
		saveFn: func(_ context.Context, key string, value testPayload, expiration time.Duration, optFns ...cache.CacheOption) error {
			attempts++
			require.Equal(t, "key", key)
			require.Equal(t, time.Minute, expiration)
			if attempts == 1 {
				return expectedErr
			}
			return nil
		},
	}

	require.NoError(t, SaveCacheItem(ctx, handle, "key", testPayload{Name: "value"}, time.Minute, cache.UseRedis()))
	require.Equal(t, 2, attempts)
}

func TestGetCacheItemNormalizesBigcacheMiss(t *testing.T) {
	ctx := context.Background()
	handle := &fakeGenericCache[testPayload]{
		getFn: func(context.Context, string, ...cache.CacheOption) (testPayload, error) {
			return testPayload{}, bigcache.ErrEntryNotFound
		},
	}

	_, err := GetCacheItem(ctx, handle, "missing", cache.UseLocal())
	require.ErrorIs(t, err, cache.ErrCacheNotFound)
	require.ErrorIs(t, err, bigcache.ErrEntryNotFound)
}

func TestGetCacheItemNoCacheIsMiss(t *testing.T) {
	ctx := context.Background()
	handle := &fakeGenericCache[testPayload]{
		getFn: func(context.Context, string, ...cache.CacheOption) (testPayload, error) {
			return testPayload{}, cache.ErrNoCache
		},
	}

	_, err := GetCacheItem(ctx, handle, "key", cache.UseNoCache())
	require.ErrorIs(t, err, cache.ErrCacheNotFound)
	require.ErrorIs(t, err, cache.ErrNoCache)
}

func TestCacheItemHelpersRejectNilHandle(t *testing.T) {
	ctx := context.Background()
	require.Error(t, SaveCacheItem(ctx, nil, "key", testPayload{}, time.Minute))
	_, err := GetCacheItem[testPayload](ctx, nil, "key")
	require.Error(t, err)
	require.Error(t, DeleteCacheItem[testPayload](ctx, nil, "key"))
}
