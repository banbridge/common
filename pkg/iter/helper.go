package iter

import (
	"reflect"
	"sync"
	"unsafe"
)

type Peeker[T any] interface {
	Iter[T]
	Peek(n int) []T
}

type peeker[T any] struct {
	Iter[T]
	buf []T
}

func ToPeeker[T any](iter Iter[T]) Peeker[T] {
	p := &peeker[T]{
		Iter: iter,
	}
	p.buf = iter.Next(1)
	return p
}

func (p *peeker[T]) Next(n int) []T {
	next := p.Peek(n)
	if len(next) == len(p.buf) {
		p.buf = nil
	} else {
		p.buf = p.buf[len(next):]
	}
	return next
}

func (p *peeker[T]) Peek(n int) []T {
	if n == 0 {
		return []T{}
	}

	if n == ALL {
		p.buf = append(p.buf, p.Iter.Next(ALL)...)
		return p.buf
	}

	if n > len(p.buf) {
		p.buf = append(p.buf, p.Iter.Next(n-len(p.buf))...)
		return p.buf
	}
	return p.buf[:n]
}

func hasNext[T any](p Peeker[T]) bool {
	return len(p.Peek(1)) != 0
}

type reflectValue struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

type hiter struct {
	key         unsafe.Pointer
	elem        unsafe.Pointer
	t           unsafe.Pointer
	h           unsafe.Pointer
	buckets     unsafe.Pointer
	overflow    *[]unsafe.Pointer
	oldoverflow *[]unsafe.Pointer
	startBucket uintptr
	offset      uint8
	wrapped     bool
	B           uint8
	i           uint8
	bucket      uintptr
	checkBucket uintptr
}

func (h *hiter) initialized() bool {
	return h.t != nil
}

//go:noescape
//go:linkname mapiternext reflect.mapiternext
func mapiternext(it *hiter)

//go:noescape
//go:linkname mapiterinit reflect.mapiterinit
func mapiterinit(rtype unsafe.Pointer, m unsafe.Pointer, it *hiter)

type unsafeMapIter[K comparable, V any] struct {
	m map[K]V
	h hiter
}

func newUnsafeMapIter[K comparable, V any](m map[K]V) *unsafeMapIter[K, V] {
	it := &unsafeMapIter[K, V]{
		m: m,
	}
	return it
}

func (it *unsafeMapIter[K, V]) Next(n int, needKey, needValue bool) ([]K, []V) {
	if n == 0 || len(it.m) == 0 {
		return nil, nil
	}

	if !it.h.initialized() && (n >= len(it.m) || n == ALL) {
		m := it.m
		it.m = nil
		if needKey && needValue {
			keys := make([]K, 0, len(m))
			values := make([]V, 0, len(m))
			for k, v := range m {
				keys = append(keys, k)
				values = append(values, v)
			}
			return keys, values
		} else if needKey {
			keys := make([]K, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			return keys, nil
		} else if needValue {
			values := make([]V, 0, len(m))
			for _, v := range m {
				values = append(values, v)
			}
			return nil, values
		}
	}

	if !it.h.initialized() {
		rv := reflect.ValueOf(it.m)
		v := (*reflectValue)(unsafe.Pointer(&rv))
		mapiterinit(v.typ, v.ptr, &it.h)
	}

	var (
		keys   []K
		values []V
	)

	if n != ALL {
		if needKey {
			keys = make([]K, 0, n)
		}
		if needValue {
			values = make([]V, 0, n)
		}
	}

	for n == ALL || n > 0 {
		if it.h.key == nil {
			break
		}
		if needKey {
			keys = append(keys, *(*K)(it.h.key))
		}
		if needValue {
			values = append(values, *(*V)(it.h.elem))
		}
		n--
		mapiternext(&it.h)
	}
	return keys, values
}

type fastpathIter[T any] struct {
	Iter[T]
	once    sync.Once
	full    func() Iter[T]
	partial func() Iter[T]
}

func newFastpathIter[T any](full func() Iter[T], partial func() Iter[T]) *fastpathIter[T] {
	return &fastpathIter[T]{
		full:    full,
		partial: partial,
	}
}

func (it *fastpathIter[T]) Next(n int) []T {
	if n == 0 {
		return []T{}
	}
	it.once.Do(func() {
		if n == ALL {
			it.Iter = it.full()
		} else {
			it.Iter = it.partial()
		}
	})
	return it.Iter.Next(n)
}
