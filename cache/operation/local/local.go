// Package local provides a bigcache-backed cache.BaseCache implementation.
package local

import (
	"context"
	"time"

	"github.com/banbridge/common/cache"
)

type localCache struct {
	cli    cache.LocalStore
	prefix string
}

// NewLocalCache creates a local BaseCache. prefix is prepended to every key.
func NewLocalCache(cli cache.LocalStore, prefix string) cache.BaseCache {
	return &localCache{cli: cli, prefix: prefix}
}

func (l *localCache) key(key string) string {
	return l.prefix + key
}

func (l *localCache) Get(_ context.Context, key string) ([]byte, error) {
	return l.cli.Get(l.key(key))
}

func (l *localCache) Save(_ context.Context, key string, data []byte, _ time.Duration) error {
	return l.cli.Set(l.key(key), data)
}

func (l *localCache) Delete(_ context.Context, key string) error {
	return l.cli.Delete(l.key(key))
}
