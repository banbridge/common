package tuple

type T2[T0, T1 any] struct {
	First  T0
	Second T1
}

func (t2 T2[T0, T1]) Value() (T0, T1) {
	return t2.First, t2.Second
}

func Make2[T0, T1 any](first T0, second T1) T2[T0, T1] {
	return T2[T0, T1]{
		First:  first,
		Second: second,
	}
}
