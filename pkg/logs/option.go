package logs

type LoggerOption func(l *SlogOption)

type SlogOption struct {
	CallDepth int
}

func WithCallDepth(callDepth int) LoggerOption {
	return func(l *SlogOption) {
		l.CallDepth = callDepth
	}
}
