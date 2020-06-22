package logging

import (
	"log"
	"testing"
)

type Writer struct {
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestLog(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("123 %d", 111)
}

func TestLogger(t *testing.T) {
	logger := New()
	// logger.SetOutput(&Writer{})
	logger.Debugf("%s%d\n", "gda", 123)
}

func BenchmarkLog(b *testing.B) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(&Writer{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Printf("%s%d\n", "gda", 123)
	}
}

func BenchmarkLogger(b *testing.B) {
	logger := New()
	logger.SetOutput(&Writer{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("%s%d\n", "gda", 123)
	}
}
