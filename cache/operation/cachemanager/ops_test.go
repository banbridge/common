package cachemanager

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"github.com/banbridge/common/cache"
)

type testPayload struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type fakeBaseCache struct {
	data       map[string][]byte
	getErr     error
	saveErr    error
	deleteErr  error
	lastKey    string
	lastValue  []byte
	lastExpire time.Duration
}

func newFakeBaseCache() *fakeBaseCache {
	return &fakeBaseCache{data: map[string][]byte{}}
}

func (f *fakeBaseCache) Get(_ context.Context, key string) ([]byte, error) {
	f.lastKey = key
	if f.getErr != nil {
		return nil, f.getErr
	}
	value, ok := f.data[key]
	if !ok {
		return nil, redis.Nil
	}
	return append([]byte(nil), value...), nil
}

func (f *fakeBaseCache) Save(_ context.Context, key string, value []byte, expire time.Duration) error {
	f.lastKey = key
	f.lastValue = append([]byte(nil), value...)
	f.lastExpire = expire
	if f.saveErr != nil {
		return f.saveErr
	}
	f.data[key] = append([]byte(nil), value...)
	return nil
}

func (f *fakeBaseCache) Delete(_ context.Context, key string) error {
	f.lastKey = key
	if f.deleteErr != nil {
		return f.deleteErr
	}
	delete(f.data, key)
	return nil
}

func (f *fakeBaseCache) IncrByExpire(_ context.Context, _ string, _ int64, _ time.Duration) (int64, error) {
	return 0, errors.New("not implemented")
}

type fakeManager struct {
	redis        cache.RedisBaseCache
	local        cache.BaseCache
	localOnRedis cache.BaseCache
	noCache      cache.BaseCache
}

func (f *fakeManager) GetRedis(context.Context) cache.RedisBaseCache   { return f.redis }
func (f *fakeManager) GetLocal(context.Context) cache.BaseCache        { return f.local }
func (f *fakeManager) GetLocalOnRedis(context.Context) cache.BaseCache { return f.localOnRedis }
func (f *fakeManager) GetNoCache(context.Context) cache.BaseCache      { return f.noCache }

func TestGetBaseCacheRoutesByOption(t *testing.T) {
	ctx := context.Background()
	redisCache := newFakeBaseCache()
	localCache := newFakeBaseCache()
	localOnRedis := newFakeBaseCache()
	noCache := newFakeBaseCache()
	customCache := newFakeBaseCache()
	manager := &fakeManager{
		redis:        redisCache,
		local:        localCache,
		localOnRedis: localOnRedis,
		noCache:      noCache,
	}
	op, err := NewCacheOperation[testPayload](manager)
	require.NoError(t, err)
	opImpl, ok := op.(*cacheOperation[testPayload])
	require.True(t, ok)

	base, err := opImpl.getBaseCache(ctx)
	require.NoError(t, err)
	require.Same(t, redisCache, base)

	base, err = opImpl.getBaseCache(ctx, cache.UseRedis())
	require.NoError(t, err)
	require.Same(t, redisCache, base)

	base, err = opImpl.getBaseCache(ctx, cache.UseLocal())
	require.NoError(t, err)
	require.Same(t, localCache, base)

	base, err = opImpl.getBaseCache(ctx, cache.UseLocalOnRedis())
	require.NoError(t, err)
	require.Same(t, localOnRedis, base)

	base, err = opImpl.getBaseCache(ctx, cache.UseNoCache())
	require.NoError(t, err)
	require.Same(t, noCache, base)

	base, err = opImpl.getBaseCache(ctx, cache.UseCustomCache(customCache))
	require.NoError(t, err)
	require.Same(t, customCache, base)

	_, err = opImpl.getBaseCache(ctx, func(opt *cache.CacheOptions) {
		opt.CacheType = cache.CacheType(999)
	})
	require.ErrorIs(t, err, ErrInvalidCacheType)

	_, err = opImpl.getBaseCache(ctx, cache.UseCustomCache(nil))
	require.ErrorIs(t, err, ErrInvalidCacheType)
}

func TestCacheOperationSaveGetDelete(t *testing.T) {
	ctx := context.Background()
	base := newFakeBaseCache()
	manager := &fakeManager{local: base}
	op, err := NewCacheOperation[testPayload](manager)
	require.NoError(t, err)
	value := testPayload{Name: "pod-a", Count: 3}

	require.NoError(t, op.Save(ctx, "cache-key", value, time.Minute, cache.UseLocal()))
	require.Equal(t, "cache-key", base.lastKey)
	require.Equal(t, time.Minute, base.lastExpire)

	got, err := op.Get(ctx, "cache-key", cache.UseLocal())
	require.NoError(t, err)
	require.Equal(t, value, got)

	require.NoError(t, op.Delete(ctx, "cache-key", cache.UseLocal()))
	_, err = op.Get(ctx, "cache-key", cache.UseLocal())
	require.ErrorIs(t, err, redis.Nil)
}

func TestCacheOperationErrors(t *testing.T) {
	ctx := context.Background()
	base := newFakeBaseCache()
	manager := &fakeManager{redis: base}
	op, err := NewCacheOperation[testPayload](manager)
	require.NoError(t, err)

	base.getErr = errors.New("cache get failed")
	_, err = op.Get(ctx, "key")
	require.ErrorContains(t, err, "cache get failed")

	base.getErr = nil
	base.data["bad-json"] = []byte("{")
	_, err = op.Get(ctx, "bad-json")
	require.Error(t, err)

	base.saveErr = errors.New("cache save failed")
	require.ErrorContains(t, op.Save(ctx, "key", testPayload{Name: "x"}, time.Second), "cache save failed")

	base.deleteErr = errors.New("cache delete failed")
	require.ErrorContains(t, op.Delete(ctx, "key"), "cache delete failed")
}

func TestNewCacheOperationRejectsNilManager(t *testing.T) {
	op, err := NewCacheOperation[testPayload](nil)
	require.Error(t, err)
	require.Nil(t, op)
}
