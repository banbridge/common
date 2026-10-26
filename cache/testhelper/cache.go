// Package testhelper provides test helpers for the cache component.
package testhelper

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	bigcache "github.com/allegro/bigcache/v3"
	redis "github.com/redis/go-redis/v9"
)

// FakeRedisClient is an in-memory  *redis.Client implementation covering
// the commands used by the cache component.
type FakeRedisClient struct {
	*redis.Client

	mu   sync.Mutex
	data map[string]string
}

// NewFakeRedisClient creates an empty in-memory Redis client.
func NewFakeRedisClient() *FakeRedisClient {
	return &FakeRedisClient{data: map[string]string{}}
}

func (f *FakeRedisClient) Get(_ context.Context, key string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	value, ok := f.data[key]
	if !ok {
		return "", redis.Nil
	}
	return value, nil
}

func (f *FakeRedisClient) GetBytes(ctx context.Context, key string) ([]byte, error) {
	value, err := f.Get(ctx, key)
	return []byte(value), err
}

func (f *FakeRedisClient) GetInt64(ctx context.Context, key string) (int64, error) {
	value, err := f.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(value, 10, 64)
}

func (f *FakeRedisClient) Set(_ context.Context, key string, value any, _ time.Duration) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	switch v := value.(type) {
	case []byte:
		f.data[key] = string(v)
	case string:
		f.data[key] = v
	default:
		f.data[key] = fmt.Sprint(v)
	}
	return nil
}

func (f *FakeRedisClient) Del(_ context.Context, keys ...string) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var deleted int64
	for _, key := range keys {
		if _, ok := f.data[key]; ok {
			delete(f.data, key)
			deleted++
		}
	}
	return deleted, nil
}

func (f *FakeRedisClient) Exists(_ context.Context, keys ...string) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	var found int64
	for _, key := range keys {
		if _, ok := f.data[key]; ok {
			found++
		}
	}
	return found, nil
}

func (f *FakeRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return f.IncrBy(ctx, key, 1)
}

func (f *FakeRedisClient) IncrBy(_ context.Context, key string, value int64) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	current := int64(0)
	if raw, ok := f.data[key]; ok {
		current, _ = strconv.ParseInt(raw, 10, 64)
	}
	current += value
	f.data[key] = strconv.FormatInt(current, 10)
	return current, nil
}

func (f *FakeRedisClient) Expire(_ context.Context, key string, _ time.Duration) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.data[key]
	return ok, nil
}

func (f *FakeRedisClient) TTL(_ context.Context, key string) (time.Duration, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.data[key]; !ok {
		return 0, redis.Nil
	}
	return time.Hour, nil
}

func (f *FakeRedisClient) Ping(_ context.Context) error {
	return nil
}

func (f *FakeRedisClient) Close() error {
	return nil
}

// SetupTestCacheClient returns a real miniredis-backed Redis client and a
// bigcache local cache for cache component tests. Tests that need a *redis.Client
// must use this, not the FakeRedisClient: the fake only implements a handful
// of methods on its own struct and embeds a nil client.
func SetupTestCacheClient(t *testing.T) (*redis.Client, *bigcache.BigCache, func()) {
	t.Helper()
	mr := miniredis.RunT(t)
	cli := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	localCache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		t.Fatalf("setup bigcache failed: %v", err)
	}
	return cli, localCache, func() {
		_ = localCache.Close()
		_ = cli.Close()
	}
}
