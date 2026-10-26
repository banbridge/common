// Package redis provides a Redis-backed cache.BaseCache implementation on top
// of the deployment-independent  *redis.Client contract.
package redis

import (
	"context"
	"time"

	"github.com/banbridge/common/cache"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	cli    *redis.Client
	prefix string
}

// NewRedisCache creates a Redis-backed BaseCache. prefix is prepended to every
// key and is expected to end with the separator used by the deployment.
func NewRedisCache(cli *redis.Client, prefix string) cache.RedisBaseCache {
	return &redisCache{cli: cli, prefix: prefix}
}

func (r *redisCache) key(key string) string {
	return r.prefix + key
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	return r.cli.Get(ctx, r.key(key)).Bytes()
}

func (r *redisCache) Save(ctx context.Context, key string, data []byte, expiration time.Duration) error {
	return r.cli.Set(ctx, r.key(key), data, expiration).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	_, err := r.cli.Del(ctx, r.key(key)).Result()
	return err
}

func (r *redisCache) IncrByExpire(ctx context.Context, key string, value int64, expiration time.Duration) (int64, error) {
	redisKey := r.key(key)
	res, err := r.cli.IncrBy(ctx, redisKey, value).Result()
	if err != nil {
		return 0, err
	}
	_, err = r.cli.Expire(ctx, redisKey, expiration).Result()
	return res, err
}
