package logging

import (
	"fmt"
	"os"
)

var std = New(os.Stderr, TraceLevel)

func SetLevel(l Level) {
	std.level = l
}

func GetLevel() Level {
	return std.level
}

func SetOutput(out WriteSyncer) {
	std.out = out
}

func Sync() error {
	return std.out.Sync()
}

func Tracef(format string, v ...interface{}) {
	std.log(TraceLevel, fmt.Sprintf(format, v...))
}

func Trace(v ...interface{}) {
	std.log(TraceLevel, fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
	std.log(DebugLevel, fmt.Sprintf(format, v...))
}

func Debug(v ...interface{}) {
	std.log(DebugLevel, fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	std.log(InfoLevel, fmt.Sprintf(format, v...))
}

func Info(v ...interface{}) {
	std.log(InfoLevel, fmt.Sprint(v...))
}

func Warningf(format string, v ...interface{}) {
	std.log(WarningLevel, fmt.Sprintf(format, v...))
}

func Warning(v ...interface{}) {
	std.log(WarningLevel, fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
	std.log(ErrorLevel, fmt.Sprintf(format, v...))
}

func Error(v ...interface{}) {
	std.log(ErrorLevel, fmt.Sprint(v...))
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	std.log(PanicLevel, s)
	panic(s)
}

func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	std.log(PanicLevel, s)
	panic(s)
}

func Fatalf(format string, v ...interface{}) {
	std.log(FatalLevel, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Fatal(v ...interface{}) {
	std.log(FatalLevel, fmt.Sprint(v...))
	os.Exit(1)
}
