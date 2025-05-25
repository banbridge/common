package iter

import "testing"

func TestPeeker_Peek(t *testing.T) {
	p := ToPeeker(FromSlice([]int{1, 2, 3, 4, 5}))
	//for i := 0; i < 5; i++ {
	//	t.Logf("peek %d: %d", i, p.Peek(3))
	//}
	//
	s := []int{1, 2, 3, 4}
	p = ToPeeker(FromSlice(s))
	for i := 0; i < 10; i++ {
		t.Logf("peek:%v next:%v", p.Peek(1), p.Next(1))
	}
}
