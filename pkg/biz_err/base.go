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

type BizError struct {
	msg        string
	code       string
	httpStatus int
	bizMsg     string

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
		err.logger.CtxError(ctx, "%s", err.Error())
	}

	return err

}

func (e *BizError) Error() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("[%s] code=%s, msg=%s, bizMsg=%s ",
		e.fnName, e.code, e.msg, e.bizMsg))
	if len(e.stack) > 0 {
		buf.Write(e.stack)
	}
	return buf.String()
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

func (e *BizError) getCurrentLocation() string {

	_, file, line, ok := runtime.Caller(e.depth)
	if !ok {
		return "??:0"
	}
	return filepath.Base(file) + ":" + strconv.Itoa(line)
}
