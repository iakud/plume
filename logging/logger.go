package logging

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const kCallerSkip int = 2

type WriteSyncer interface {
	io.Writer
	Sync() error
}

type Logger struct {
	mu    sync.Mutex
	out   WriteSyncer
	level Level
	hooks Hooks
}

func New(out WriteSyncer, l Level, hooks ...Hook) *Logger {
	logger := &Logger{
		out:   out,
		level: l,
		hooks: hooks,
	}
	return logger
}

func (logger *Logger) SetOutput(out WriteSyncer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.out = out
}

func (logger *Logger) SetLevel(l Level) {
	atomic.StoreInt32((*int32)(&logger.level), int32(l))
}

func (logger *Logger) GetLevel() Level {
	return Level(atomic.LoadInt32((*int32)(&logger.level)))
}

func (logger *Logger) Enabled(level Level) bool {
	return logger.GetLevel() <= level
}

func (logger *Logger) AddHook(hook Hook) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.hooks.Add(hook)
}

func (logger *Logger) Sync() error {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	return logger.out.Sync()
}

func (logger *Logger) log(l Level, s string) {
	now := time.Now() // get this early.
	pc, file, line, ok := runtime.Caller(kCallerSkip)
	if !ok {
		file = "???"
		line = 1
	}
	entry := Entry{now, l, s, pc, file, line}
	// hook
	logger.logHooks(&entry)
	// write
	buf := newBuffer()
	defer buf.free()
	buf.formatHeader(entry.Time, entry.Level, entry.File, entry.Line)
	buf.appendString(entry.Message)
	if len(entry.Message) == 0 || entry.Message[len(entry.Message)-1] != '\n' {
		buf.appendByte('\n')
	}
	logger.mu.Lock()
	defer logger.mu.Unlock()
	if _, err := logger.out.Write(buf.bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write: %v\n", err)
	}
	if entry.Level > ErrorLevel {
		logger.out.Sync()
	}
}

func (logger *Logger) logHooks(entry *Entry) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.hooks.log(entry)
}

func (logger *Logger) Tracef(format string, v ...interface{}) {
	if logger.Enabled(TraceLevel) {
		logger.log(TraceLevel, fmt.Sprintf(format, v...))
	}
}

func (logger *Logger) Trace(v ...interface{}) {
	if logger.Enabled(TraceLevel) {
		logger.log(TraceLevel, fmt.Sprint(v...))
	}
}

func (logger *Logger) Debugf(format string, v ...interface{}) {
	if logger.Enabled(DebugLevel) {
		logger.log(DebugLevel, fmt.Sprintf(format, v...))
	}
}

func (logger *Logger) Debug(v ...interface{}) {
	if logger.Enabled(DebugLevel) {
		logger.log(DebugLevel, fmt.Sprint(v...))
	}
}

func (logger *Logger) Infof(format string, v ...interface{}) {
	if logger.Enabled(InfoLevel) {
		logger.log(InfoLevel, fmt.Sprintf(format, v...))
	}
}

func (logger *Logger) Info(v ...interface{}) {
	if logger.Enabled(InfoLevel) {
		logger.log(InfoLevel, fmt.Sprint(v...))
	}
}

func (logger *Logger) Warningf(format string, v ...interface{}) {
	if logger.Enabled(WarningLevel) {
		logger.log(WarningLevel, fmt.Sprintf(format, v...))
	}
}

func (logger *Logger) Warning(v ...interface{}) {
	if logger.Enabled(WarningLevel) {
		logger.log(WarningLevel, fmt.Sprint(v...))
	}
}

func (logger *Logger) Errorf(format string, v ...interface{}) {
	if logger.Enabled(ErrorLevel) {
		logger.log(ErrorLevel, fmt.Sprintf(format, v...))
	}
}

func (logger *Logger) Error(v ...interface{}) {
	if logger.Enabled(ErrorLevel) {
		logger.log(ErrorLevel, fmt.Sprint(v...))
	}
}

func (logger *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if logger.Enabled(PanicLevel) {
		logger.log(PanicLevel, s)
	}
	panic(s)
}

func (logger *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	if logger.Enabled(PanicLevel) {
		logger.log(PanicLevel, s)
	}
	panic(s)
}

func (logger *Logger) Fatalf(format string, v ...interface{}) {
	if logger.Enabled(FatalLevel) {
		logger.log(FatalLevel, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

func (logger *Logger) Fatal(v ...interface{}) {
	if logger.Enabled(FatalLevel) {
		logger.log(FatalLevel, fmt.Sprint(v...))
	}
	os.Exit(1)
}
