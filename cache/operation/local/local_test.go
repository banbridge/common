package local

import (
	"context"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/stretchr/testify/require"
)

func TestLocalCache(t *testing.T) {
	ctx := context.Background()
	bigCache, err := bigcache.New(ctx, bigcache.DefaultConfig(10*time.Minute))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, bigCache.Close()) })

	cli := NewLocalCache(bigCache, "cn_phone:local:")

	_, err = cli.Get(ctx, "1")
	require.ErrorIs(t, err, bigcache.ErrEntryNotFound)

	require.NoError(t, cli.Save(ctx, "1", []byte("value"), 0))
	got, err := cli.Get(ctx, "1")
	require.NoError(t, err)
	require.Equal(t, "value", string(got))

	require.NoError(t, cli.Delete(ctx, "1"))
	_, err = cli.Get(ctx, "1")
	require.ErrorIs(t, err, bigcache.ErrEntryNotFound)
}
