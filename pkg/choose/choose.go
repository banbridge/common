package choose

// If returns onTrue if cond is true, otherwise returns onFalse.
func If[T any](cond bool, onTrue, onFalse T) T {
	if cond {
		return onTrue
	}
	return onFalse
}

type Lazy[T any] func() T

func IfLazy[T any](cond bool, onTrue, onFalse Lazy[T]) T {
	if cond {
		return onTrue()
	}
	return onFalse()
}

func IfLazyLeft[T any](cond bool, onTrue Lazy[T], onFalse T) T {
	if cond {
		return onTrue()
	}
	return onFalse
}

func IfLazyRight[T any](cond bool, onTrue T, onFalse Lazy[T]) T {
	if cond {
		return onTrue
	}
	return onFalse()
}
