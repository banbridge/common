package gvalue

import (
	"sync"
	"unsafe"

	"cmp"
)

func Zero[T any]() T {
	var zero T
	return zero
}

func Max[T cmp.Ordered](a T, b ...T) T {
	res := a
	for _, v := range b {
		if v > res {
			res = v
		}
	}
	return res
}

func Min[T cmp.Ordered](a T, b ...T) T {
	res := a
	for _, v := range b {
		if v < res {
			res = v
		}
	}
	return res
}

func MinMax[T cmp.Ordered](a T, b ...T) (minRes, maxRes T) {
	minRes, maxRes = a, a
	for _, v := range b {
		if v < minRes {
			minRes = v
		}
		if v > maxRes {
			maxRes = v
		}
	}
	return
}

// Clamp returns the value v clamped to the range [min, max].
//
// If v is less than min, Clamp returns min.
// If v is greater than max, Clamp returns max.
// Otherwise, Clamp returns v.
//
// Example usage:
//

func Clamp[T cmp.Ordered](v, min, max T) T {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

type vauleface struct {
	x    uintptr
	data unsafe.Pointer
}

// IsNil reports whether v is nil.
//
// It panics if v's Kind is not Interface or Ptr.
//
// IsNil is equivalent to v == nil, but is more efficient.
func IsNil(v any) bool {
	return (*vauleface)(unsafe.Pointer(&v)).data == nil
}

func IsNotNil(v any) bool {
	return !IsNil(v)
}

func IsZero[T comparable](v T) bool {
	return v == Zero[T]()
}

func IsNotZero[T comparable](v T) bool {
	return !IsZero(v)
}

func TypeAssert[To, From any](v From) To {
	return any(v).(To)
}

func TypeAssertTry[To, From any](v From) (To, bool) {
	res, ok := any(v).(To)
	return res, ok
}

func Cast[To, From Numeric](v From) To {
	return To(v)
}

func Once[T any](v func() T) func() T {

	var (
		once sync.Once
		res  T
	)

	return func() T {
		once.Do(func() {
			res = v()
		})
		return res
	}

}
