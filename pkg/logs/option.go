package logs

type LoggerOption func(l *SlogOption)

type SlogOption struct {
	CallDepth int
}

func getDefaultOpt() *SlogOption {
	return &SlogOption{
		CallDepth: 1,
	}
}

func WithCallDepth(callDepth int) LoggerOption {
	return func(l *SlogOption) {
		l.CallDepth = callDepth
	}
}
