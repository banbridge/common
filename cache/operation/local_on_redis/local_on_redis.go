// Package localonredis provides a two-level cache that reads local first and
// falls back to Redis, backfilling the local cache on a Redis hit.
package localonredis

import (
	"context"
	"log/slog"
	"time"

	redis "github.com/redis/go-redis/v9"

	"github.com/banbridge/common/cache"
)

type localOnRedis struct {
	redisCli *redis.Client
	localCli cache.LocalStore
	prefix   string
}

// NewLocalOnRedis creates a two-level BaseCache. prefix is prepended to both
// the local and Redis keys, matching the original deployment behavior.
func NewLocalOnRedis(redisCli *redis.Client, localCli cache.LocalStore, prefix string) cache.BaseCache {
	return &localOnRedis{redisCli: redisCli, localCli: localCli, prefix: prefix}
}

func (l *localOnRedis) key(key string) string {
	return l.prefix + key
}

func (l *localOnRedis) Get(ctx context.Context, key string) ([]byte, error) {
	redisKey := l.key(key)
	data, err := l.localCli.Get(redisKey)
	if err == nil {
		return data, nil
	}

	data, err = l.redisCli.Get(ctx, redisKey).Bytes()
	if err == nil {
		if errL := l.localCli.Set(redisKey, data); errL != nil {
			slog.WarnContext(ctx, "cache local backfill failed", "err", errL)
		}
	}
	return data, err
}

func (l *localOnRedis) Save(ctx context.Context, key string, value []byte, expire time.Duration) error {
	redisKey := l.key(key)
	if err := l.localCli.Set(redisKey, value); err != nil {
		slog.WarnContext(ctx, "cache local save failed", "err", err)
	}
	return l.redisCli.Set(ctx, redisKey, value, expire).Err()
}

func (l *localOnRedis) Delete(ctx context.Context, key string) error {
	redisKey := l.key(key)
	if err := l.localCli.Delete(redisKey); err != nil {
		slog.WarnContext(ctx, "cache local delete failed", "err", err)
	}

	return l.redisCli.Del(ctx, redisKey).Err()
}
