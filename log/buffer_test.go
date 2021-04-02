package log

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestBuffer(t *testing.T) {
	buf := newBuffer()
	defer buf.free()
	_, file, line, _ := runtime.Caller(0)
	buf.formatHeader(time.Now(), InfoLevel, file, line)
	buf.appendString("test buffer")
	fmt.Println(string(buf.bytes()))
}

func BenchmarkBuffer(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := newBuffer()
			buf.free()
		}
	})
}

func BenchmarkFormatHeader(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := newBuffer()
			_, file, line, _ := runtime.Caller(0)
			buf.formatHeader(time.Now(), InfoLevel, file, line)
			buf.free()
		}
	})
}
