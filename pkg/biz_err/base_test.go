package biz_err

import (
	"context"
	"fmt"
	"testing"

	"github.com/banbridge/common/pkg/logs"
)

func TestStack(t *testing.T) {
	err := getErr(context.Background())
	if err != nil {
		return
	}
}

func getErr(ctx context.Context) error {
	return NewError(ctx, "1001", "error")
}

func TestLog1(t *testing.T) {
	logs.Info("msg")

}

func TestFromError(t *testing.T) {
	ctx := context.Background()

	err := FromError(ctx, nil)

	fmt.Println(err.Error())
}

func TestWrapError(t *testing.T) {

}
