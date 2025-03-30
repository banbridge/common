package biz_err

import (
	"context"
	"errors"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownBizCode is unknown bizCode for error info.
	UnknownBizCode = "100000"
	// UnknownReason is unknown reason for error info.
	UnknownReason = "UnknownReason"
	// UnknownMessage is unknown message for error info.
	UnknownMessage = ""
	// UnknownBizMsg is unknown message for error info.
	UnknownBizMsg = "未知原因"
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

func FromError(ctx context.Context, err error) *BizError {
	se := &BizError{}

	if err == nil {
		return se
	}

	if errors.As(err, &se) {
		return se
	}

	se = NewError(ctx, UnknownBizCode, err.Error(),
		WithHttpStatus(UnknownCode),
		WithDepth(3),
	)

	return se
}

func StatusCode(ctx context.Context, err error) string {
	return FromError(ctx, err).Code()
}

func HttpCode(ctx context.Context, err error) int {
	return FromError(ctx, err).HttpCode()
}
