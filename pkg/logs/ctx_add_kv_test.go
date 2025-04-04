package logs

import (
	"context"
	"testing"
)

func TestKVs(t *testing.T) {
	ctx := context.Background()

	ctx = ctxAddKVs(ctx, "key", "value")

	kvs := GetAllKVs(ctx)

	CtxInfo(ctx, "kvs:%+v", kvs)
}
