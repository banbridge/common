package cache

import (
	"context"
	"time"
)

type fakeBaseCache struct{}

func (f *fakeBaseCache) Get(context.Context, string) ([]byte, error) {
	return nil, nil
}

func (f *fakeBaseCache) Save(context.Context, string, []byte, time.Duration) error {
	return nil
}

func (f *fakeBaseCache) Delete(context.Context, string) error {
	return nil
}

func (f *fakeBaseCache) IncrByExpire(context.Context, string, int64, time.Duration) (int64, error) {
	return 0, nil
}
