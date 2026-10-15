package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeCacheOptions(t *testing.T) {
	require.Equal(t, RedisCache, MergeCacheOptions().CacheType)
	require.Equal(t, LocalCache, MergeCacheOptions(UseLocal()).CacheType)
	require.Equal(t, LocalOnRedisLevel, MergeCacheOptions(UseLocalOnRedis()).CacheType)
	require.Equal(t, NoCache, MergeCacheOptions(UseNoCache()).CacheType)
	require.Equal(t, RedisCache, MergeCacheOptions(UseRedis()).CacheType)
}

func TestUseCustomCache(t *testing.T) {
	bc := &fakeBaseCache{}
	opt := MergeCacheOptions(UseCustomCache(bc))
	require.Equal(t, CustomCache, opt.CacheType)
	require.Same(t, bc, opt.CustomBaseCache)
}
