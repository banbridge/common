package iter

const (
	ALL = -1
)

type Iter[T any] interface {
	Next(n int) []T
}
