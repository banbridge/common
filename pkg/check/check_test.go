package check

import "testing"

func TestSliceNotEmpty(t *testing.T) {
	f := SliceNotEmpty[string](nil)
	t.Logf("result: %v", f)
}
