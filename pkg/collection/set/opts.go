package set

import "github.com/banbridge/common/pkg/choose"

type Option[T comparable] func(*Set[T])

func Members[T comparable](members ...T) Option[T] {
	return func(s *Set[T]) {
		if s.data == nil {
			s.data = make(map[T]struct{}, choose.If(len(members) == 0, initSize, len(members)))
		}
		for _, member := range members {
			s.data[member] = struct{}{}
		}
	}
}

func Selector[F any, T comparable](items []F, selector func(F) T) Option[T] {
	return func(s *Set[T]) {
		if s.data == nil {
			s.data = make(map[T]struct{}, choose.If(len(items) == 0, initSize, len(items)))
		}
		for _, item := range items {
			s.data[selector(item)] = struct{}{}
		}
	}
}
