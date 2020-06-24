package logging

import (
	"fmt"
	"log"
	"os"
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
	logger := NewLogger("./", "testlogger.log")
	logger.Debugf("123 %d", 111)
}

func BenchmarkLogger(b *testing.B) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	file, err := os.OpenFile("./benchmarLogger.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		b.Fatal(err)
	}
	log.SetOutput(file)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Output(2, fmt.Sprintf("%s%d\n", "gda", 123))
	}
}

func BenchmarkBufferedLogger(b *testing.B) {
	logger := NewLogger("./", "benchmarBufferedLogger.log")
	defer logger.Sync()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("%s%d\n", "gda", 123)
	}
}

func BenchmarkLoggerInParallel(b *testing.B) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	file, err := os.OpenFile("./benchmarLogger2.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		b.Fatal(err)
	}
	log.SetOutput(file)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Output(2, fmt.Sprintf("%s%d\n", "gda", 123))
		}
	})
}

func BenchmarkBufferedLoggerInParallel(b *testing.B) {
	logger := NewLogger("./", "benchmarBufferedLogger2.log")
	defer logger.Sync()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Infof("%s%d\n", "gda", 123)
		}
	})
}
