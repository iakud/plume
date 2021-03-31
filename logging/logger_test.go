package logging

import (
	"fmt"
	"io"
	"os"
	"testing"
)

type nullWriter struct {
	io.Writer
}

func (*nullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (*nullWriter) Sync() error {
	return nil
}

func TestLogger(t *testing.T) {
	SetLevel(TraceLevel)

	Tracef("hello %d world!", TraceLevel)
	Debugf("hello %d world!", DebugLevel)
	Infof("hello %d world!", InfoLevel)
	Warningf("hello %d world!", WarningLevel)
	Errorf("hello %d world!", ErrorLevel)

	Trace("hello ", int(TraceLevel), " world!")
	Debug("hello ", int(DebugLevel), " world!")
	Info("hello ", int(InfoLevel), " world!")
	Warning("hello ", int(WarningLevel), " world!")
	Error("hello ", int(ErrorLevel), " world!")
}

func TestHook(t *testing.T) {
	warningHook := func(e *Entry) error {
		if WarningLevel == e.Level {
			fmt.Fprintln(os.Stdout, "hook warning:", e.Message)
		}
		return nil
	}
	errorHook := func(e *Entry) error {
		if ErrorLevel == e.Level {
			fmt.Fprintln(os.Stdout, "hook error:", e.Message)
		}
		return nil
	}
	SetOutput(&nullWriter{})
	SetLevel(TraceLevel)
	AddHook(warningHook)
	AddHook(errorHook)

	Info("This is info log!")
	Warning("This is warning log!")
	Error("This is error log!")
}

func BenchmarkLogger(b *testing.B) {
	logger := New(&nullWriter{}, TraceLevel)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 1; pb.Next(); i++ {
			logger.Info("hello", i)
		}
	})
}
