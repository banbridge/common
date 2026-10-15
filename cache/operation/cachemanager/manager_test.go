package cachemanager

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/banbridge/common/cache"
	"github.com/banbridge/common/cache/testhelper"
)

func TestNewCacheManager(t *testing.T) {
	ctx := context.Background()
	redisCli, localCli, cleanup := testhelper.SetupTestCacheClient(t)
	defer cleanup()

	manager, err := NewCacheManager(redisCli, localCli, cache.Config{
		RedisPrefix: "cn_phone:",
		LocalPrefix: "cn_phone:local:",
	})
	require.NoError(t, err)
	require.NotNil(t, manager.GetRedis(ctx))
	require.NotNil(t, manager.GetLocal(ctx))
	require.NotNil(t, manager.GetLocalOnRedis(ctx))
	require.NotNil(t, manager.GetNoCache(ctx))

	value := testPayload{Name: "hit", Count: 1}
	require.NoError(t, manager.GetRedis(ctx).Save(ctx, "key", []byte(`{"name":"hit","count":1}`), time.Minute))
	got, err := manager.GetRedis(ctx).Get(ctx, "key")
	require.NoError(t, err)
	require.JSONEq(t, string(got), `{"name":"hit","count":1}`)
	_ = value
}

func TestNewCacheManagerRejectsNilInputs(t *testing.T) {
	redisCli, localCli, cleanup := testhelper.SetupTestCacheClient(t)
	defer cleanup()

	_, err := NewCacheManager(nil, localCli, cache.DefaultConfig())
	require.Error(t, err)

	_, err = NewCacheManager(redisCli, nil, cache.DefaultConfig())
	require.Error(t, err)
}
