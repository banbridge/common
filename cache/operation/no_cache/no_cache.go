// Package nocache provides a no-op cache.BaseCache implementation.
package nocache

import (
	"context"
	"time"

	"github.com/banbridge/common/cache"
)

type noCache struct{}

// NewNoCache creates a no-op BaseCache. Reads and deletes return
// cache.ErrNoCache so helper functions can treat them as cache misses.
func NewNoCache() cache.BaseCache {
	return &noCache{}
}

func (e *noCache) Get(_ context.Context, _ string) ([]byte, error) {
	return nil, cache.ErrNoCache
}

func (e *noCache) Save(_ context.Context, _ string, _ []byte, _ time.Duration) error {
	return nil
}

func (e *noCache) Delete(_ context.Context, _ string) error {
	return cache.ErrNoCache
}
