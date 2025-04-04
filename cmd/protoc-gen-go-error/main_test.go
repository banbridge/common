package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	errors "github.com/banbridge/common/pkg/biz_err"
)

func TestReadProto(t *testing.T) {

	reader, _ := os.Open("/Users/banbridge/VScodeProjects/banbridge/petal/errorx/error.proto")
	defer reader.Close()

	//parser := proto.NewParser(reader)
	//definition, _ := parser.Parse()
	//
	//for _, element := range definition.Elements {
	//	t.Logf("%v\n", element)
	//}

}

func IsParamsErr(ctx context.Context, err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(ctx, err)
	return e.Code() == "100001"
}

// 参数错误
func ErrorParamsErr(ctx context.Context, format string, args ...any) *errors.BizError {
	return errors.NewError(ctx, "100001", fmt.Sprintf(format, args...),
		errors.WithHttpStatus(500), errors.WithBizMsg("参数相关错误，请检查相关参数。"), errors.WithReason(""), errors.WithDepth(3))
}
