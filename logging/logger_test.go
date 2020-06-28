package logging

import (
	"fmt"
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

func BenchmarkLogger(b *testing.B) {
	name := fmt.Sprintf("%s.log", b.Name())
	writer := NewFileWriter(name)
	SetOutput(writer)
	defer Sync()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Infof("%s %d\n", "hello", 1)
		}
	})
}
