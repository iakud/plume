package log

import (
	"bytes"
	"fmt"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	logDepth(s severity, message fmt.Stringer, depth int)
}

type logger struct {
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
