package logs

import "context"

const (
	kvCtxKey = "CTX_KVS"
)

type kvCtx struct {
	kvs []any
	pre *kvCtx
}

func CtxAddKVs(ctx context.Context, kvs ...any) context.Context {
	return ctxAddKVs(ctx, kvs...)
}

func ctxAddKVs(ctx context.Context, kvs ...any) context.Context {
	if len(kvs) == 0 || len(kvs)&1 == 1 {
		return ctx
	}

	kvList := make([]any, 0, len(kvs))
	kvList = append(kvList, kvs...)

	return context.WithValue(ctx, kvCtxKey, &kvCtx{
		kvs: kvList,
		pre: getKVS(ctx),
	})
}

func getKVS(ctx context.Context) *kvCtx {
	v := ctx.Value(kvCtxKey)
	if v == nil {
		return nil
	}
	if kvs, ok := v.(*kvCtx); ok {
		return kvs
	}
	return nil
}

func GetAllKVs(ctx context.Context) []any {
	if ctx == nil {
		return nil
	}
	kvs := getKVS(ctx)
	var result []any
	recursiveAllKVs(&result, kvs, 0)
	return result
}

func recursiveAllKVs(res *[]any, kvs *kvCtx, total int) {
	if kvs == nil {
		*res = make([]any, 0, total)
		return
	}
	recursiveAllKVs(res, kvs.pre, total+len(kvs.kvs))
	*res = append(*res, kvs.kvs...)
}
