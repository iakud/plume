package logging

import (
	"io"
	"io/ioutil"
	"testing"
)

func TestLogger(t *testing.T) {
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

type nullWriter struct {
	io.Writer
}

func (w *nullWriter) Sync() error {
	return nil
}

func BenchmarkLogger(b *testing.B) {
	logger := NewLogger(&nullWriter{ioutil.Discard}, TraceLevel)
	SetLogger(logger)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 1; pb.Next(); i++ {
			Info("hello", i)
		}
	})
}
