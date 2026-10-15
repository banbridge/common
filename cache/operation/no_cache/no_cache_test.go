package nocache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/banbridge/common/cache"
)

func TestNoCache(t *testing.T) {
	ctx := context.Background()
	c := NewNoCache()

	_, err := c.Get(ctx, "key")
	require.ErrorIs(t, err, cache.ErrNoCache)
	require.NoError(t, c.Save(ctx, "key", []byte("value"), time.Minute))
	require.ErrorIs(t, c.Delete(ctx, "key"), cache.ErrNoCache)
}
