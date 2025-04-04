package biz_err

type ErrorOption func(*BizError)

func WithErrorStack(withStack bool) ErrorOption {
	return func(e *BizError) {
		e.withStack = withStack
	}
}

func WithDepth(depth int) ErrorOption {
	return func(e *BizError) {
		e.depth = depth
	}
}

func WithHttpStatus(status int) ErrorOption {
	return func(e *BizError) {
		e.httpStatus = status
	}
}

func WithBizMsg(msg string) ErrorOption {
	return func(e *BizError) {
		e.bizMsg = msg
	}
}

func WithReason(reason string) ErrorOption {
	return func(e *BizError) {
		e.reason = reason
	}
}

func WithLogLevel(level ErrorLevel) ErrorOption {
	return func(e *BizError) {
		e.level = level
	}
}
