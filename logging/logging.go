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
	GetLogger().log(TraceLevel, fmt.Sprintf(format, v...))
}

func Trace(v ...interface{}) {
	GetLogger().log(TraceLevel, fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	GetLogger().log(DebugLevel, fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	GetLogger().log(DebugLevel, fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	GetLogger().log(InfoLevel, fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	GetLogger().log(InfoLevel, fmt.Sprint(v...))
}

func Warningf(format string, v ...interface{}) {
	GetLogger().log(WarningLevel, fmt.Sprintf(format, v...))
}

func Warning(v ...interface{}) {
	GetLogger().log(WarningLevel, fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	GetLogger().log(ErrorLevel, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	GetLogger().log(ErrorLevel, fmt.Sprint(v...))
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	GetLogger().log(PanicLevel, s)
	panic(s)
}

func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	GetLogger().log(PanicLevel, s)
	panic(s)
}

func Fatalf(format string, v ...interface{}) {
	GetLogger().log(FatalLevel, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Fatal(v ...interface{}) {
	GetLogger().log(FatalLevel, fmt.Sprint(v...))
	os.Exit(1)
}
