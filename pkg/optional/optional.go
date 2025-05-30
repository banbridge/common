package optional

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/banbridge/common/pkg/choose"
	"github.com/banbridge/common/pkg/gvalue"
)

type O[T any] struct {
	val T
	ok  bool
}

func Of[T any](value T) O[T] {
	return O[T]{
		val: value,
		ok:  true,
	}
}

func OfPtr[T any](value *T) O[T] {
	if value == nil {
		return Nil[T]()
	}
	return OK(*value)
}

func OK[T any](v T) O[T] {
	return O[T]{
		ok:  true,
		val: v,
	}
}

func Nil[T any]() O[T] {
	return O[T]{}
}

func (o O[T]) Must() T {
	if !o.ok {
		panic(fmt.Errorf("no valid value in optional.O[%s]", o.typ()))
	}
	return o.val
}

func (o O[T]) Value() T {
	return o.val
}

func (o O[T]) ValueOr(v T) T {
	return choose.If(o.ok, o.val, v)
}

func (o O[T]) ValueOrZero() T {
	return choose.If(o.ok, o.val, gvalue.Zero[T]())
}

func (o O[T]) Ptr() *T {
	if !o.ok {
		return nil
	}
	return &o.val
}

func (o O[T]) Get() (T, bool) {
	return o.val, o.ok
}

func (o O[T]) IsOK() bool {
	return o.ok
}

func (o O[T]) IsNil() bool {
	return !o.ok
}

func (o O[T]) IfOk(f func(T)) {
	if o.ok {
		f(o.val)
	}
}

func (o O[T]) IfNil(f func()) {
	if !o.ok {
		f()
	}
}

func (o O[T]) typ() string {
	typ := reflect.TypeOf(gvalue.Zero[T]())
	if typ == nil {
		return "any"
	}
	return typ.String()
}

func (o O[T]) String() string {
	if !o.ok {
		return fmt.Sprintf("optional.O[%s](nil)", o.typ())
	}
	return fmt.Sprintf("optional.O[%s](%v)", o.typ(), o.val)
}

func (o O[T]) MarshalJSON() ([]byte, error) {
	if !o.ok {
		return []byte("null"), nil
	}
	return json.Marshal(o.val)
}

func (o O[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	err := json.Unmarshal(data, &o.val)
	if err != nil {
		return err
	}
	o.ok = true
	return nil
}

func Map[T, V any](o O[T], f func(T) V) O[V] {
	if !o.ok {
		return Nil[V]()
	}
	return OK(f(o.val))
}

func Then[T, V any](o O[T], f func(T) O[V]) O[V] {
	if !o.ok {
		return Nil[V]()
	}
	return f(o.val)
}
