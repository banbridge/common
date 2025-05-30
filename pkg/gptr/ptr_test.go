package gptr

import "testing"

type Foo struct {
}

func TestIsNil(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		if !IsNil(nil) {
			t.Fail()
		}
		if !IsNil((*Foo)(nil)) {
			t.Fail()
		}

		var foo *Foo
		if !IsNil(foo) {
			t.Fail()
		}
	})

}
