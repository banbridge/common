package ctx_values

import (
	"context"
	"testing"

	"github.com/banbridge/common/pkg/logs"
)

func TestKVs(t *testing.T) {
	ctx := context.Background()

	ctx = ctxAddKVs(ctx, "key", "value")

	kvs := GetAllKVs(ctx)

	logs.CtxInfo(ctx, "kvs:%+v", kvs)
}
