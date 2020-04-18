package log

import (
	"bufio"
	"io"
	"sync"
	"time"
)

type logWriter struct {
	writer *bufio.Writer
	mutex  sync.Mutex
	done   chan struct{}
}

func newLogWriter(w io.Writer, bufferSize int, period time.Duration) *logWriter {
	newWriter := new(logWriter)
	newWriter.writer = bufio.NewWriterSize(w, bufferSize)
	newWriter.done = make(chan struct{})
	go newWriter.periodicalFlush()
	return newWriter
}

func (this *logWriter) Write(p []byte) (n int, err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.writer.Write(p)
}

func (this *logWriter) Flush() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	return this.writer.Flush()
}

func (this *logWriter) periodicalFlush(period time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			this.Flush()
		case <-this.done:
			return
		}
	}
}
