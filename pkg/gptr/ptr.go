package gptr

import "github.com/banbridge/common/pkg/gvalue"

func Of[T any](v T) *T {
	return &v
}

func OfNotZero[T comparable](v T) *T {
	if gvalue.IsZero(v) {
		return nil
	}
	return &v
}

func OfPositive[T gvalue.Numeric](v T) *T {
	if v <= 0 {
		return nil
	}
	return &v
}

func OfNegative[T gvalue.Numeric](v T) *T {
	if v >= 0 {
		return nil
	}
	return &v
}

func Indirect[T any](v *T) (r T) {

	if v == nil {
		return
	}
	return *v
}

func IndirectOr[T any](v *T, or T) T {
	if v == nil {
		return or
	}
	return *v
}

func IsNil(v any) bool {
	return gvalue.IsNil(v)
}
