package iter

type emptyIter[T any] struct {
}

func (e *emptyIter[T]) Next(_ int) []T {
	return nil
}
