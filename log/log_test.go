package log

import (
	"testing"
	"time"
)

func BenchmarkLog(b *testing.B) {
	fw := NewFileWriter(b.Name()+".log", time.Hour, 0)
	defer fw.Flush()
	SetOutput(fw)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for i := 1; pb.Next(); i++ {
			Info(b.Name(), "abcdefghijklmnopqrstuvwxyz 1234567890 abcdefghijklmnopqrstuvwxyz", i)
		}
	})
}
