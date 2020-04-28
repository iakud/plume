package logger

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type logging struct {
}

func (this *logger) printDepth(s severity, message fmt.Stringer, depth int) {

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

func (this *logging) Debug(args ...interface{}) {
	this.print(debugLog, args...)
}

func (this *logging) Debugf(format string, args ...interface{}) {
	this.printf(debugLog, format, args...)
}

func (this *logging) Info(args ...interface{}) {
	this.print(infoLog, args...)
}

func (this *logging) Infof(format string, args ...interface{}) {
	this.printf(infoLog, format, args...)
}

func (this *logging) Warn(args ...interface{}) {
	this.print(warningLog, args...)
}

func (this *logging) Warnf(format string, args ...interface{}) {
	this.printf(warningLog, format, args...)
}

func (this *logging) Error(args ...interface{}) {
	this.print(errorLog, args...)
}

func (this *logging) Errorf(format string, args ...interface{}) {
	this.printf(errorLog, format, args...)
}

func (this *logging) Fatal(args ...interface{}) {
	this.print(fatalLog, args...)
}

func (this *logging) Fatalf(format string, args ...interface{}) {
	this.printf(fatalLog, format, args...)
}
