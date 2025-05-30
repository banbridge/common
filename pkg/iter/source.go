package iter

import (
	"context"
	"math"

	"github.com/banbridge/common/pkg/gvalue"
	"github.com/banbridge/common/pkg/tuple"
)

type sliceIter[T any] struct {
	s []T
}

func (i *sliceIter[T]) Next(n int) []T {
	if n == ALL || n >= len(i.s) {
		n = len(i.s)
	}
	next := make([]T, n)
	copy(next, i.s[:n])
	i.s = i.s[n:]
	return next
}

func FromSlice[T any](s []T) Iter[T] {
	return &sliceIter[T]{
		s: s,
	}
}

type stealSliceIter[T any] struct {
	s []T
}

func (i *stealSliceIter[T]) Next(n int) []T {
	if n == ALL || n >= len(i.s) {
		n = len(i.s)
	}
	next := i.s[:n]
	i.s = i.s[n:]
	return next
}

func StealSlice[T any](s []T) Iter[T] {
	return &stealSliceIter[T]{
		s: s,
	}
}

type mapKeyValueIter[K comparable, V any] struct {
	i *unsafeMapIter[K, V]
}

func (i *mapKeyValueIter[K, V]) Next(n int) []tuple.T2[K, V] {
	keys, values := i.i.Next(n, true, true)

	if len(keys) == 0 {
		return nil
	}
	next := make([]tuple.T2[K, V], len(keys))
	for idx := range keys {
		next[idx] = tuple.Make2(keys[idx], values[idx])
	}
	return next
}

func FromMap[K comparable, V any](m map[K]V) Iter[tuple.T2[K, V]] {
	return &mapKeyValueIter[K, V]{
		newUnsafeMapIter[K, V](m),
	}
}

type mapKeyIter[K comparable, V any] struct {
	i *unsafeMapIter[K, V]
}

func (i *mapKeyIter[K, V]) Next(n int) []K {
	keys, _ := i.i.Next(n, true, false)
	return keys
}

func FromMapKeys[K comparable, V any](m map[K]V) Iter[K] {
	return &mapKeyIter[K, V]{
		newUnsafeMapIter[K, V](m),
	}
}

type mapValueIter[K comparable, V any] struct {
	i *unsafeMapIter[K, V]
}

func (i *mapValueIter[K, V]) Next(n int) []V {
	_, values := i.i.Next(n, false, true)
	return values
}

func FromMapValues[K comparable, V any](m map[K]V) Iter[V] {
	return &mapValueIter[K, V]{
		newUnsafeMapIter[K, V](m),
	}
}

type chanIter[T any] struct {
	ch  <-chan T
	ctx context.Context
}

func (i *chanIter[T]) Next(n int) (r []T) {
	if n == 0 {
		return
	}

	if n == ALL {
		for {
			select {
			case <-i.ctx.Done():
				return
			case v, ok := <-i.ch:
				if !ok {
					return
				}
				r = append(r, v)

			}
		}
	}

	r = make([]T, 0, n)
	for j := 0; j < n; j++ {
		select {
		case <-i.ctx.Done():
			return
		case v, ok := <-i.ch:
			if !ok {
				return
			}
			r = append(r, v)
		}
	}
	return
}

func FromChan[T any](ctx context.Context, ch <-chan T) Iter[T] {
	return &chanIter[T]{
		ch:  ch,
		ctx: ctx,
	}
}

type rangeIter[T gvalue.Numeric] struct {
	cur  T
	stop T
	step T
}

func (i *rangeIter[T]) Next(n int) (r []T) {
	intervalLen := math.Abs(float64(i.stop - i.cur))

	step := math.Abs(float64(i.step))

	l := int(math.Ceil(intervalLen / step))
	if n == 0 || l == 0 {
		return
	}

	if n == ALL || n >= l {
		n = l
	}

	j := 0

	r = make([]T, n)
	for j < n {
		r[j] = i.cur
		i.cur += i.step
		j++
	}
	return
}

func Range[T gvalue.Numeric](start, stop T) Iter[T] {
	return RangeWithStep(start, stop, 1)

}

func RangeWithStep[T gvalue.Numeric](start, stop, step T) Iter[T] {
	return &rangeIter[T]{
		cur:  start,
		stop: stop,
		step: step,
	}
}

type repeatIter[T any] struct {
	val T
}

func (i *repeatIter[T]) Next(n int) (r []T) {
	if n == ALL {
		panic("repeat iter can't be used with ALL")
	}
	r = make([]T, n)
	for j := 0; j < n; j++ {
		r[j] = i.val
	}
	return
}

func Repeat[T any](val T) Iter[T] {
	return &repeatIter[T]{
		val: val,
	}
}
