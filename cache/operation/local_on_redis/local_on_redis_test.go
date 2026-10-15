package localonredis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/banbridge/common/cache/testhelper"
)

func TestLocalOnRedisBackfillsLocal(t *testing.T) {
	ctx := context.Background()
	redisCli, localCli, cleanup := testhelper.SetupTestCacheClient(t)
	defer cleanup()

	ca := NewLocalOnRedis(redisCli, localCli, "cn_phone:")
	require.NoError(t, redisCli.Set(ctx, "cn_phone:remote", "remote-value", time.Minute).Err())

	got, err := ca.Get(ctx, "remote")
	require.NoError(t, err)
	require.Equal(t, "remote-value", string(got))

	localValue, err := localCli.Get("cn_phone:remote")
	require.NoError(t, err)
	require.Equal(t, "remote-value", string(localValue))
}
