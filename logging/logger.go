package logging

import (
	"runtime"
	"sync/atomic"
)

const kCallerSkip int = 2

type Logger struct {
	level Level
}

func (logger *Logger) SetLevel(level Level) {
	atomic.StoreInt32((*int32)(&logger.level), int32(level))
}

func (logger *Logger) GetLevel() Level {
	return Level(atomic.LoadInt32((*int32)(&logger.level)))
}

func (logger *Logger) IsLevelDisabled(level Level) bool {
	return logger.GetLevel() > level
}

func (logger *Logger) logf(level Level, format string, a ...interface{}) {
	if logger.level.Disabled(level) {
		return
	}
	pc, file, line, ok := runtime.Caller(kCallerSkip)
	if !ok {
		file = "???"
		line = 1
	}
	caller := newCaller(pc, file, line)
	_ = caller
}

func (logger *Logger) Tracef(format string, a ...interface{}) {
	logger.logf(TraceLevel, format, a)
}

func (logger *Logger) Debugf(format string, a ...interface{}) {
	logger.logf(DebugLevel, format, a)
}

func (logger *Logger) Infof(format string, a ...interface{}) {
	logger.logf(InfoLevel, format, a)
}

func (logger *Logger) Warnf(format string, a ...interface{}) {
	logger.logf(WarnLevel, format, a)
}

func (logger *Logger) Errorf(format string, a ...interface{}) {
	logger.logf(ErrorLevel, format, a)
}

func (logger *Logger) Panicf(format string, a ...interface{}) {
	logger.logf(PanicLevel, format, a)
}

func (logger *Logger) Fatalf(format string, a ...interface{}) {
	logger.logf(FatalLevel, format, a)
}

type Caller struct {
	PC   uintptr
	File string
	Line int
}

func newCaller(pc uintptr, file string, line int) Caller {
	return Caller{
		PC:   pc,
		File: file,
		Line: line,
	}
}
