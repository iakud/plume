package log

import (
	"bytes"
	"fmt"
)

type Logger struct {
}

func (this *Logger) print(s severity, args ...interface{}) {
	b := bytes.Buffer{}
	fmt.Fprint(b, args...)
	b.WriteByte('\n')
}

func (this *Logger) printf(s severity, format string, args ...interface{}) {
	b := bytes.Buffer{}
	fmt.Fprintf(b, format, args...)
	b.WriteByte('\n')
}

type innerLogger interface {
	innerLog()
}
