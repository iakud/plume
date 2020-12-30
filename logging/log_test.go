package logging

import (
	"testing"
	"time"
)

func BenchmarkLog(b *testing.B) {
	fw := NewFileWriter(b.Name()+".log", time.Hour, 0)
	defer fw.Flush()
	logger := NewLogger(fw, TraceLevel)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 1; pb.Next(); i++ {
			logger.Info("hello", i)
		}
	})
}
