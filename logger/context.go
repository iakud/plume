package logger

import (
	"runtime"
	"strings"
)

type logContext struct {
	file string
	line int
}

func getContext(depth int) *logContext {
	_, file, line, ok := runtime.Caller(3 + depth)
	if !ok {
		file = "???"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return &logContext{file, line}
}
