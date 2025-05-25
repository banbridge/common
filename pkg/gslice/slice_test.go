package gslice

import (
	"strconv"
	"testing"
)

func TestFilterMap(t *testing.T) {

	res := []int{1, 2, 3, 0, 0, 0}
	ans := FilterMap(res, func(item int) (string, bool) {
		return strconv.Itoa(item), item != 0
	})

	t.Log(ans)
}

func TestRemove(t *testing.T) {

	nums := []int{1}

	res := Remove(nums, 1)
	t.Log(res)

}

func TestSort(t *testing.T) {

	nums := []int{20, 11, 3, 4, 5, 6, 7, 8, 9, 10}
	Sort(nums, func(a, b int) bool {
		return a < b
	})

	t.Log(nums)

}

func TestMap(t *testing.T) {
	mp1 := map[int]int{
		1: 1,
		2: 2,
		3: 3,
	}

	mp1[5] = 5
	mp1[4] = 4
	t.Log(mp1)
}

func TestDivide(t *testing.T) {

	s := []int{0, 1, 2, 3, 4}
	res := Divide(s, 3)
	t.Log(res)

}
