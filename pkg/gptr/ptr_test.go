package gptr

import "testing"

type Foo struct {
}

func TestIsPtr(t *testing.T) {
	f := IsPtr(1)
	t.Logf("result: %v", f)

}
