package check

type SelfChecker interface {
	Check() error
}
