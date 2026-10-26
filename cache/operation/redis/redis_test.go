package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"

	"github.com/banbridge/common/cache/testhelper"
)

func TestRedisCache(t *testing.T) {
	ctx := context.Background()
	redisCli, _, cleanup := testhelper.SetupTestCacheClient(t)
	defer cleanup()

	cli := NewRedisCache(redisCli, "cn_phone:")

	_, err := cli.Get(ctx, "test")
	require.ErrorIs(t, err, redis.Nil)

	require.NoError(t, cli.Save(ctx, "test", []byte("value"), time.Minute))
	got, err := cli.Get(ctx, "test")
	require.NoError(t, err)
	require.Equal(t, "value", string(got))

	value, err := cli.IncrByExpire(ctx, "counter", 10, time.Minute)
	require.NoError(t, err)
	require.Equal(t, int64(10), value)

	require.NoError(t, cli.Delete(ctx, "test"))
	_, err = cli.Get(ctx, "test")
	require.ErrorIs(t, err, redis.Nil)
}
