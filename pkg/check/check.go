package check

import (
	"errors"
	"fmt"
	"strings"

	"github.com/banbridge/common/pkg/ptr"
)

type ErrContext struct {
	Err error
}

func (e *ErrContext) IsError() bool {
	return e.Err != nil
}

// IsTrue 检查是否为true
func IsTrue(ok bool, format string, args ...any) error {
	if !ok {
		return errors.New(fmt.Sprintf(format, args...))
	}
	return nil
}

// IsFalse 检查是否为false
func IsFalse(ok bool, format string, args ...any) error {
	if ok {
		return errors.New(fmt.Sprintf(format, args...))
	}
	return nil
}

// NotBlank 检查是否为空
func NotBlank(str string, format string, args ...any) error {

	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return errors.New(fmt.Sprintf(format, args...))
	}
	return nil
}

// NotNil 检查是否为空
func NotNil(obj any, format string, args ...any) error {
	if gptr.IsNil(obj) {
		return errors.New(fmt.Sprintf(format, args...))
	}
	return nil
}

// WrapErrorWithRes
func WrapErrorWithRes[R any](f func() (R, error), ec *ErrContext) R {
	var res R
	if ec.Err != nil {
		return res
	}
	res, ec.Err = f()
	return res
}

// WrapError 包装错误
func WrapError(f func() error, ec *ErrContext) {
	if ec.Err != nil {
		return
	}
	ec.Err = f()
}

// SliceNotEmpty 检查切片是否为空
func SliceNotEmpty[T any](list []T, format string, args ...any) error {
	if len(list) == 0 {
		return errors.New(fmt.Sprintf(format, args...))
	}
	return nil
}

// MapNotEmpty 检查map是否为空
func MapNotEmpty[T any](m map[string]T, format string, args ...any) error {
	if len(m) == 0 {
		return errors.New(fmt.Sprintf(format, args...))
	}
	return nil
}
