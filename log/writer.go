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
	ticker *time.Ticker
	done   chan struct{}
}

func newLogWriter(w io.Writer, bufferSize int, period time.Duration) *logWriter {
	newWriter := new(logWriter)
	newWriter.writer = bufio.NewWriterSize(w, bufferSize)
	newWriter.ticker = time.NewTicker(period)
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

func (this *logWriter) periodicalFlush() {
	for {
		select {
		case <-this.ticker.C:
			this.Flush()
		case <-this.done:
			return
		}
	}
}
