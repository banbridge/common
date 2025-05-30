package gmap

import "testing"

func TestEqual(t *testing.T) {
	mp1 := map[int]int{
		1: 1,
		2: 2,
		3: 3,
	}
	mp2 := map[int]int{
		1: 1,
		2: 2,
		3: 3,
	}

	if !Equal(mp1, mp2) {
		t.Errorf("expected %v, but found %v", mp1, mp2)
	}
}
