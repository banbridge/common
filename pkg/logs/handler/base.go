package handler

import (
	"fmt"
	"path"
	"runtime"
)

func getCallInfo(pc uintptr) string {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	file := path.Base(f.File)
	return fmt.Sprintf("%s:%d", file, f.Line)
}
