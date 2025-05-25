package iter

import (
	"testing"

	"github.com/banbridge/common/pkg/gmap"
	"github.com/banbridge/common/pkg/gslice"
)

func TestFromSlice(t *testing.T) {
	var empty []int

	val := ToSlice(FromSlice[int](nil))
	if !gslice.Equal(val, empty) {

		t.Errorf("expected %v, but found %v", empty, val)
	}

	source := []int{1, 2, 3, 4, 5}
	val = ToSlice(FromSlice[int]([]int{1, 2, 3, 4, 5}))
	if !gslice.Equal(val, source) {
		t.Errorf("expected %v, but found %v", source, val)
	}

}

func TestMap(t *testing.T) {
	source := map[int]string{}
	val := KVToMap(FromMap[int, string](nil))
	t.Logf("expected %v, but found %v", source, val)

	source = map[int]string{1: "1", 2: "2", 3: "3", 4: "4", 5: "5"}
	val = KVToMap(FromMap[int, string](map[int]string{1: "1", 2: "2", 3: "3", 4: "4", 5: "5"}))
	if !gmap.Equal(val, source) {
		t.Errorf("expected %v, but found %v", source, val)
	}
}
