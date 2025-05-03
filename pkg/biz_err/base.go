package biz_err

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/banbridge/common/pkg/logs"
)

// ErrorLevel is the level of error.
type ErrorLevel int8

const (
	// LevelInfo  is debug level.
	LevelInfo ErrorLevel = iota
	// LevelWarn is warn level.
	LevelWarn
	// LevelError is error level.
	LevelError
)

type BizError struct {
	base       error
	msg        string
	code       string
	httpStatus int
	bizMsg     string
	reason     string

	level  ErrorLevel
	fnName string
	stack  []byte

	logger    logs.Logger
	withStack bool
	depth     int
}

func NewError(ctx context.Context, code, msg string, opts ...ErrorOption) *BizError {
	err := &BizError{
		code:      code,
		msg:       msg,
		withStack: true,
		depth:     2,
		logger:    logs.DefaultLogger,
	}

	for _, opt := range opts {
		opt(err)
	}

	if err.withStack {
		err.stack = err.Stack()
	}

	if len(err.fnName) == 0 {
		err.fnName = err.getCurrentLocation()
	}

	if err.logger != nil {
		err.logger.CtxError(ctx, "%s", err.String())
	}

	return err

}

func (e *BizError) String() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("[%s] code=%s, msg=%s, bizMsg=%s ",
		e.fnName, e.code, e.msg, e.bizMsg))
	if len(e.stack) > 0 {
		buf.Write(e.stack)
	}
	return buf.String()
}

func (e *BizError) Error() string {
	return e.msg
}

// Format implements fmt.Formatter.
func (e *BizError) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v', 's':
		f.Write([]byte(e.Error()))
		return
	}
}

func (e *BizError) Code() string {
	return e.code
}

func (e *BizError) HttpCode() int {
	return e.httpStatus
}

func (e *BizError) Reason() string {
	return e.reason
}

func (e *BizError) Stack() []byte {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(e.depth, pcs[:])

	frames := runtime.CallersFrames(pcs[e.depth-1 : n])

	var buffer bytes.Buffer

	for {
		frame, more := frames.Next()

		buffer.WriteString(frame.Function)
		buffer.WriteByte('\n')
		buffer.WriteByte('\t')
		buffer.WriteString(frame.File)
		buffer.WriteByte(':')
		buffer.WriteString(strconv.Itoa(frame.Line))
		buffer.WriteByte('\n')

		if !more {
			break
		}
	}
	return buffer.Bytes()
}

func (e *BizError) Wrap(err error) *BizError {
	e.base = err
	return e
}

func (e *BizError) Unwrap() error {
	return e.base
}

func (e *BizError) getCurrentLocation() string {

	_, file, line, ok := runtime.Caller(e.depth)
	if !ok {
		return "??:0"
	}
	return filepath.Base(file) + ":" + strconv.Itoa(line)
}

// logFunc defines log print functions.
type logFunc func(ctx context.Context, format string, v ...interface{})

func (e *BizError) ctxLog(ctx context.Context) {
	e.getLogFunc()(ctx, "%s", e.String())
}

func (e *BizError) getLogFunc() logFunc {
	switch e.level {
	case LevelInfo:
		return logs.CtxInfo
	case LevelWarn:
		return logs.CtxWarn
	case LevelError:
		return logs.CtxError
	}
	return logs.CtxError
}
