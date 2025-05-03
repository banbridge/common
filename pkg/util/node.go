package util

import "github.com/banbridge/common/pkg/check"

func NodeExec[C any](nodeFunc func(context C) error, context C, ec *check.ErrContext) {
	if ec.IsError() {
		return
	}
	ec.Err = nodeFunc(context)
}
