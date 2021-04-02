package log

import (
	"fmt"
	"os"
)

var std = New(os.Stderr, InfoLevel)

func SetLevel(l Level) {
	std.SetLevel(l)
}

func GetLevel() Level {
	return std.GetLevel()
}

func SetOutput(out WriteSyncer) {
	std.SetOutput(out)
}

func AddHook(hook Hook) {
	std.AddHook(hook)
}

func Sync() error {
	return std.Sync()
}

func Tracef(format string, v ...interface{}) {
	if std.Enabled(TraceLevel) {
		std.log(TraceLevel, fmt.Sprintf(format, v...))
	}
}

func Trace(v ...interface{}) {
	if std.Enabled(TraceLevel) {
		std.log(TraceLevel, fmt.Sprint(v...))
	}
}

func Debugf(format string, v ...interface{}) {
	if std.Enabled(DebugLevel) {
		std.log(DebugLevel, fmt.Sprintf(format, v...))
	}
}

func Debug(v ...interface{}) {
	if std.Enabled(DebugLevel) {
		std.log(DebugLevel, fmt.Sprint(v...))
	}
}

func Infof(format string, v ...interface{}) {
	if std.Enabled(InfoLevel) {
		std.log(InfoLevel, fmt.Sprintf(format, v...))
	}
}

func Info(v ...interface{}) {
	if std.Enabled(InfoLevel) {
		std.log(InfoLevel, fmt.Sprint(v...))
	}
}

func Warningf(format string, v ...interface{}) {
	if std.Enabled(WarningLevel) {
		std.log(WarningLevel, fmt.Sprintf(format, v...))
	}
}

func Warning(v ...interface{}) {
	if std.Enabled(WarningLevel) {
		std.log(WarningLevel, fmt.Sprint(v...))
	}
}

func Errorf(format string, v ...interface{}) {
	if std.Enabled(ErrorLevel) {
		std.log(ErrorLevel, fmt.Sprintf(format, v...))
	}
}

func Error(v ...interface{}) {
	if std.Enabled(ErrorLevel) {
		std.log(ErrorLevel, fmt.Sprint(v...))
	}
}

func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if std.Enabled(PanicLevel) {
		std.log(PanicLevel, s)
	}
	panic(s)
}

func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	if std.Enabled(PanicLevel) {
		std.log(PanicLevel, s)
	}
	panic(s)
}

func Fatalf(format string, v ...interface{}) {
	if std.Enabled(FatalLevel) {
		std.log(FatalLevel, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

func Fatal(v ...interface{}) {
	if std.Enabled(FatalLevel) {
		std.log(FatalLevel, fmt.Sprint(v...))
	}
	os.Exit(1)
}
