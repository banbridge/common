package biz_err

type ErrorOption func(*BizError)

func WithErrorStack(withStack bool) ErrorOption {
	return func(e *BizError) {
		e.withStack = withStack
	}
}

func WithSkipDepth(depth int) ErrorOption {
	return func(e *BizError) {
		e.depth = depth
	}
}

func WithHttpStatus(status int) ErrorOption {
	return func(e *BizError) {
		e.httpStatus = status
	}
}
