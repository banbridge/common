package gormrepo

import "testing"

func TestOptions(t *testing.T) {
	nums := []int{1}

	t.Logf("index 0, %d", nums[0])
	t.Logf("index 1, %v", nums[1:])

}
