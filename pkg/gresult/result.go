package gresult

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/banbridge/common/pkg/choose"
	"github.com/banbridge/common/pkg/gptr"
	"github.com/banbridge/common/pkg/gvalue"
	"github.com/banbridge/common/pkg/optional"
)

type R[T any] struct {
	val T
	err error
}

func Of[T any](val T, err error) R[T] {
	checkErr(err)
	return R[T]{
		val: val,
		err: err,
	}
}

func OK[T any](val T) R[T] {
	return R[T]{
		val: val,
		err: nil,
	}
}

func Err[T any](e error) R[T] {
	if gvalue.IsNil(e) {
		panic(fmt.Errorf("expected a non-nil error, but found nil error with type %T", e))
	}
	return R[T]{
		err: e,
	}
}

func (r R[T]) Value() T {
	return r.val
}

func (r R[T]) ValueOr(v T) {
	choose.If(r.err == nil, r.val, v)
}

func (r R[T]) ValueOrZero() T {
	return choose.If(r.err == nil, r.val, gvalue.Zero[T]())
}

func (r R[T]) Get() (T, error) {
	return r.val, r.err
}

func (r R[T]) Err() error {
	return r.err
}

func (r R[T]) Must() T {
	if r.err != nil {
		panic(fmt.Errorf("unexpected error in gresult.R[%s]: %T; %s", r.typ(), r.err, r.err))
	}
	return r.val
}
func (r R[T]) IsOK() bool {
	return r.err == nil
}

func (r R[T]) IsErr() bool {
	return r.err != nil
}

func (r R[T]) IfOK(f func(T)) {
	if r.err == nil {
		f(r.val)
	}
}

func (r R[T]) IfErr(f func(error)) {
	if r.err != nil {
		f(r.err)
	}
}

func (r R[T]) typ() string {
	typ := reflect.TypeOf(gvalue.Zero[T]())
	if typ == nil {
		return "any"
	}
	return typ.String()
}

func (r R[T]) String() string {
	if r.err != nil {
		return fmt.Sprintf("gresult.Err[%s](%s)", r.typ(), r.err)
	}
	return fmt.Sprintf("gresult.Ok[%s](%v)", r.typ(), r.val)
}

type jsonR[T any] struct {
	Val *T      `json:"value,omitempty"`
	Err *string `json:"error,omitempty"`
}

func (r R[T]) MarshalJSON() ([]byte, error) {
	jr := jsonR[T]{}
	if r.err != nil {
		jr.Err = gptr.Of(r.err.Error())
	} else {
		jr.Val = gptr.Of(r.val)
	}
	return json.Marshal(jr)
}

func (r R[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	jr := jsonR[T]{}
	err := json.Unmarshal(data, &jr)
	if err != nil {
		return err
	}

	if jr.Val != nil && jr.Err != nil {
		return errors.New("gresult: neither error nor value is nil")
	}

	if jr.Err == nil && jr.Val == nil {
		r.val = gvalue.Zero[T]()
		r.err = nil
	} else if jr.Err != nil {
		r.val = gvalue.Zero[T]()
		r.err = errors.New(*jr.Err)
	} else {
		r.val = *jr.Val
		r.err = nil
	}

	return nil

}

func checkErr(e error) {
	if e != nil && gvalue.IsNil(e) {
		panic(fmt.Errorf("error with type %T is nil", e))
	}
}

func Map[T, V any](r R[T], f func(T) V) R[V] {
	if r.err != nil {
		return Err[V](r.err)
	}
	return OK(f(r.val))
}

func MapErr[T any](r R[T], f func(error) error) R[T] {
	if r.err != nil {
		return Err[T](f(r.err))
	}
	return r
}

func Then[T, V any](r R[T], f func(T) R[V]) R[V] {
	if r.err != nil {
		return Err[V](r.err)
	}
	return f(r.val)
}

func (r R[T]) Optional() optional.O[T] {
	if r.err != nil {
		return optional.Nil[T]()
	}
	return optional.OK(r.val)
}
