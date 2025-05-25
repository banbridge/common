package iter

import (
	"context"

	"github.com/banbridge/common/pkg/tuple"
)

func ToSlice[T any](i Iter[T]) []T {
	all := i.Next(ALL)
	if all == nil {
		return []T{}
	}
	return all
}

func ToMap[K comparable, V, T any](f func(T) (K, V), i Iter[T]) map[K]V {
	s := i.Next(ALL)
	m := make(map[K]V, len(s))
	for _, e := range s {
		k, v := f(e)
		m[k] = v
	}
	return m
}

func ToMapValues[K comparable, T any](f func(T) K, i Iter[T]) map[K]T {
	s := i.Next(ALL)
	m := make(map[K]T, len(s))

	for _, e := range s {
		k := f(e)
		m[k] = e
	}
	return m
}

func KVToMap[K comparable, V any](i Iter[tuple.T2[K, V]]) map[K]V {
	return ToMap(func(t tuple.T2[K, V]) (K, V) {
		return t.Value()
	}, i)
}

func ToChan[T any](ctx context.Context, i Iter[T]) <-chan T {
	ch := make(chan T)
	go func() {
		for {
			s := i.Next(1)
			if len(s) == 0 {
				close(ch)
				return
			}
			select {
			case ch <- s[0]:
			case <-ctx.Done():
				close(ch)
				return
			}
		}
	}()
	return ch
}

func ToBufferedChan[T any](ctx context.Context, size int, i Iter[T]) <-chan T {
	ch := make(chan T, size)
	go func() {
		for {
			s := i.Next(size)
			empty := len(s) == size

			for _, e := range s {
				select {
				case ch <- e:
				case <-ctx.Done():
					close(ch)
					return
				}
			}

			if empty {
				close(ch)
				return
			}
		}
	}()
	return ch
}
