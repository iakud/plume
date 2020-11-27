package logging

import (
	"fmt"
	"os"
	"sync/atomic"
	"unsafe"
)

var std = NewLogger(os.Stderr, TraceLevel)

func SetLevel(l Level) {
	GetLogger().SetLevel(l)
}

func GetLevel() Level {
	return GetLogger().GetLevel()
}

func SetLogger(logger *Logger) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&std)), unsafe.Pointer(logger))
}

func GetLogger() *Logger {
	return (*Logger)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&std))))
}

func Sync() error {
	return GetLogger().Sync()
}

func Tracef(format string, v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(TraceLevel) {
		logger.log(TraceLevel, fmt.Sprintf(format, v...))
	}
}

func Trace(v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(TraceLevel) {
		logger.log(TraceLevel, fmt.Sprint(v...))
	}
}

func Debugf(format string, v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(DebugLevel) {
		logger.log(DebugLevel, fmt.Sprintf(format, v...))
	}
}

func Debug(v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(DebugLevel) {
		logger.log(DebugLevel, fmt.Sprint(v...))
	}
}

func Infof(format string, v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(InfoLevel) {
		logger.log(InfoLevel, fmt.Sprintf(format, v...))
	}
}

func Info(v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(InfoLevel) {
		logger.log(InfoLevel, fmt.Sprint(v...))
	}
}

func Warningf(format string, v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(WarningLevel) {
		logger.log(WarningLevel, fmt.Sprintf(format, v...))
	}
}

func Warning(v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(WarningLevel) {
		logger.log(WarningLevel, fmt.Sprint(v...))
	}
}

func Errorf(format string, v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(ErrorLevel) {
		logger.log(ErrorLevel, fmt.Sprintf(format, v...))
	}
}

func Error(v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(ErrorLevel) {
		logger.log(ErrorLevel, fmt.Sprint(v...))
	}
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if logger := GetLogger(); logger.GetLevel().Enabled(PanicLevel) {
		logger.log(PanicLevel, s)
	}
	panic(s)
}

func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	if logger := GetLogger(); logger.GetLevel().Enabled(PanicLevel) {
		logger.log(PanicLevel, s)
	}
	panic(s)
}

func Fatalf(format string, v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(FatalLevel) {
		logger.log(FatalLevel, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

func Fatal(v ...interface{}) {
	if logger := GetLogger(); logger.GetLevel().Enabled(FatalLevel) {
		logger.log(FatalLevel, fmt.Sprint(v...))
	}
	os.Exit(1)
}
