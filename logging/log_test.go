package logging

import (
	"testing"
	"time"
)

func BenchmarkLog(b *testing.B) {
	fw := NewFileWriter(b.Name()+".log", time.Hour, 0)
	defer fw.Flush()
	logger := New(fw, TraceLevel)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 1; pb.Next(); i++ {
			logger.Info(b.Name(), "abcdefghijklmnopqrstuvwxyz 1234567890 abcdefghijklmnopqrstuvwxyz", i)
		}
	})
}
