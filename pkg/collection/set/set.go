package set

import (
	"fmt"
	"sort"
	"strings"

	"github.com/banbridge/common/pkg/choose"
)

const (
	initSize = 32
)

type Set[T comparable] struct {
	data map[T]struct{}
}

func New[T comparable](members ...T) *Set[T] {
	s := &Set[T]{}

	s.data = make(map[T]struct{}, choose.If(len(members) != 0, len(members), initSize))

	for _, member := range members {
		s.data[member] = struct{}{}
	}

	return s
}

func NewWith[T comparable](opts ...Option[T]) *Set[T] {
	s := &Set[T]{}
	for _, opt := range opts {
		opt(s)
	}

	if s.data == nil {
		s.data = make(map[T]struct{}, initSize)
	}
	return s
}

func (s *Set[T]) Len() int {
	return len(s.data)
}

func (s *Set[T]) Add(member T) bool {
	if _, ok := s.data[member]; !ok {
		s.data[member] = struct{}{}
		return true
	}
	return false
}

func (s *Set[T]) Remove(member T) bool {
	if _, ok := s.data[member]; ok {
		delete(s.data, member)
		return true
	}
	return false
}

func (s *Set[T]) AddN(members ...T) {
	for _, member := range members {
		s.data[member] = struct{}{}
	}
}

func (s *Set[T]) RemoveN(members ...T) {
	for _, member := range members {
		delete(s.data, member)
	}
}

func (s *Set[T]) Contains(member T) bool {
	_, ok := s.data[member]
	return ok
}

func (s *Set[T]) ContainsAny(members ...T) bool {
	for _, member := range members {
		if _, ok := s.data[member]; ok {
			return true
		}
	}
	return false
}

func (s *Set[T]) ContainsAll(members ...T) bool {
	for _, member := range members {
		if _, ok := s.data[member]; !ok {
			return false
		}
	}
	return true
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	res := New[T]()
	for member := range s.data {
		res.data[member] = struct{}{}
	}
	for member := range other.data {
		res.data[member] = struct{}{}
	}
	return res
}

func (s *Set[T]) Diff(other *Set[T]) *Set[T] {
	res := New[T]()
	for member := range s.data {
		if _, ok := other.data[member]; !ok {
			res.data[member] = struct{}{}
		}
	}
	return res
}

func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	res := New[T]()
	for member := range s.data {
		if _, ok := other.data[member]; ok {
			res.data[member] = struct{}{}
		}
	}
	return res
}

func (s *Set[T]) Update(other *Set[T]) {
	s.UnionInPlace(other)
}

func (s *Set[T]) UnionInPlace(other *Set[T]) {
	for member := range other.data {
		s.data[member] = struct{}{}
	}
}

func (s *Set[T]) DiffInPlace(other *Set[T]) {
	for member := range other.data {
		if _, ok := s.data[member]; ok {
			delete(s.data, member)
		}
	}
}

func (s *Set[T]) IntersectInPlace(other *Set[T]) {
	for member := range s.data {
		if _, ok := other.data[member]; !ok {
			delete(s.data, member)
		}
	}
}

func (s *Set[T]) Euqls(other *Set[T]) bool {
	if s.Len() != other.Len() {
		return false
	}
	for member := range s.data {
		if _, ok := other.data[member]; !ok {
			return false
		}
	}
	return true
}

// IsSuperSet 检查s是否是other的超集
func (s *Set[T]) IsSuperSet(other *Set[T]) bool {
	if s.Len() < other.Len() {
		return false
	}
	for member := range other.data {
		if _, ok := s.data[member]; !ok {
			return false
		}
	}
	return true
}

func (s *Set[T]) String() string {
	members := make([]string, 0, len(s.data))

	for member := range s.data {
		members = append(members, fmt.Sprintf("%v", member))
	}

	sort.Strings(members)
	return fmt.Sprintf("set{%s}", strings.Join(members, ", "))

}

func (s *Set[T]) ToSlice() []T {
	members := make([]T, 0, len(s.data))
	for member := range s.data {
		members = append(members, member)
	}
	return members
}

func (s *Set[T]) Clear() {
	s.data = make(map[T]struct{}, initSize)
}

func (s *Set[T]) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Set[T]) Clone() *Set[T] {
	res := New[T]()
	for member := range s.data {
		res.data[member] = struct{}{}
	}
	return res
}
