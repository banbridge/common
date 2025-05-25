package iter

import "testing"

func TestReverse(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}

	it := Take(2, Reserve(FromSlice(a)))

	re := it.Next(ALL)
	t.Logf("result: %v", re)
}

func TestZip(t *testing.T) {
	p1 := ToPeeker(FromSlice([]int{1, 2, 3, 4, 5}))
	p2 := ToPeeker(FromSlice([]int{1, 2, 3, 4, 5}))

	zi := Zip(func(t1, t2 int) int {
		return t1 + t2
	}, p1, p2)
	re := zi.Next(2)
	t.Logf("res:%v", re)
}

func TestIntersperse(t *testing.T) {
	it := Intersperse(0, FromSlice([]int{1, 2, 3, 4, 5}))
	re := ToSlice(it)
	t.Logf("res:%v", re)
}

func TestPrepend(t *testing.T) {
	//nums := []int{0, 1, 2, 3, 4}
	it := Prepend(0, FromSlice([]int{1, 2, 3}))
	re := ToSlice(it)
	t.Logf("res:%v", re)
}

func TestUniq(t *testing.T) {
	nums := []int{3, 1, 2, 3, 3, 5}

	it := Uniq(FromSlice(nums))
	re := ToSlice(it)
	t.Logf("res:%v", re)
}

func TestDup(t *testing.T) {
	nums := []int{3, 1, 2, 3, 3, 5}
	it := Dup(FromSlice(nums))
	re := ToSlice(it)
	t.Logf("res:%v", re)
}

func TestChunk(t *testing.T) {
	nums := []int{3, 1, 2, 3, 3, 5}
	it := Chunk(2, FromSlice(nums))
	re := ToSlice(it)
	t.Logf("res:%v", re)
}

func TestCompact(t *testing.T) {
	nums := []int{3, 1, 0, 2, 3, 3, 5}
	it := Compact(FromSlice(nums))
	re := ToSlice(it)
	t.Logf("res:%v", re)
}
