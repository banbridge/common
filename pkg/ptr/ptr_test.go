package ptr

import "testing"

type Foo struct {
}

func TestIsPtr(t *testing.T) {
	f := IsPtr[*Foo](nil)
	t.Logf("result: %v", f)

}
