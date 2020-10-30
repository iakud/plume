package logging

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

const kCallerSkip int = 2

type WriteSyncer interface {
	io.Writer
	Sync() error
}

type Logger struct {
	out   WriteSyncer
	level Level
	hooks Hooks
}

func NewLogger(out WriteSyncer, l Level, hooks ...Hook) *Logger {
	logger := &Logger{
		out:   out,
		level: l,
		hooks: hooks,
	}
	return logger
}

func (logger *Logger) SetLevel(l Level) {
	atomic.StoreInt32((*int32)(&logger.level), int32(l))
}

func (logger *Logger) GetLevel() Level {
	return Level(atomic.LoadInt32((*int32)(&logger.level)))
}

func (logger *Logger) Sync() error {
	return logger.out.Sync()
}

func (logger *Logger) log(l Level, s string) {
	if !logger.GetLevel().Enabled(l) {
		return
	}
	now := time.Now() // get this early.
	pc, file, line, ok := runtime.Caller(kCallerSkip)
	if !ok {
		file = "???"
		line = 1
	}

	entry := Entry{now, l, s, pc, file, line}
	// hook
	if logger.hooks != nil {
		logger.hooks.log(&entry)
	}
	// write
	buf := newBuffer()
	buf.formatHeader(entry.Time, entry.Level, entry.File, entry.Line)
	buf.appendString(entry.Message)
	if len(entry.Message) == 0 || s[len(entry.Message)-1] != '\n' {
		buf.appendByte('\n')
	}
	if _, err := logger.out.Write(buf.bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write: %v\n", err)
	}
	buf.free()

	if entry.Level > ErrorLevel {
		logger.out.Sync()
	}
}

func (logger *Logger) Tracef(format string, v ...interface{}) {
	logger.log(TraceLevel, fmt.Sprintf(format, v...))
}

func (logger *Logger) Trace(v ...interface{}) {
	logger.log(TraceLevel, fmt.Sprint(v...))
}

func (logger *Logger) Debugf(format string, v ...interface{}) {
	logger.log(DebugLevel, fmt.Sprintf(format, v...))
}

func (logger *Logger) Debug(v ...interface{}) {
	logger.log(DebugLevel, fmt.Sprint(v...))
}

func (logger *Logger) Infof(format string, v ...interface{}) {
	logger.log(InfoLevel, fmt.Sprintf(format, v...))
}

func (logger *Logger) Info(v ...interface{}) {
	logger.log(InfoLevel, fmt.Sprint(v...))
}

func (logger *Logger) Warningf(format string, v ...interface{}) {
	logger.log(WarningLevel, fmt.Sprintf(format, v...))
}

func (logger *Logger) Warning(v ...interface{}) {
	logger.log(WarningLevel, fmt.Sprint(v...))
}

func (logger *Logger) Errorf(format string, v ...interface{}) {
	logger.log(ErrorLevel, fmt.Sprintf(format, v...))
}

func (logger *Logger) Error(v ...interface{}) {
	logger.log(ErrorLevel, fmt.Sprint(v...))
}

func (logger *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logger.log(PanicLevel, s)
	panic(s)
}

func (logger *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	logger.log(PanicLevel, s)
	panic(s)
}

func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.log(FatalLevel, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (logger *Logger) Fatal(v ...interface{}) {
	logger.log(FatalLevel, fmt.Sprint(v...))
	os.Exit(1)
}
