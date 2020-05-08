package logger

import (
	"bufio"
	"context"
	"io"
	"sync"
	"time"
)

const bufferSize = 256 * 1024
const flushInterval = 30 * time.Second

type bufferedWriter struct {
	mutex     sync.Mutex
	writer    io.Writer
	bufWriter *bufio.Writer

	ctx    context.Context
	cancel context.CancelFunc
}

func newBufferedWriter(writer io.Writer, bufferSize int, flushPeriod time.Duration) {

	bufWriter := bufio.NewWriterSize(writer, bufferSize)
	ctx, cancel := context.WithCancel(context.Background())
	bufferedWriter := &bufferedWriter{
		writer:    writer,
		bufWriter: bufWriter,

		ctx:    ctx,
		cancel: cancel,
	}

	if flushPeriod > 0 {
		go bufferedWriter.flushPeriodically()
	}
	return bufferedWriter
}

func (this *bufferedWriter) Write(p []byte) (int, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	length := len(p)
	if length > this.bufWriter.Size() {
		if err := this.bufWriter.Flush(); err != nil {
			return 0, err
		}
		return this.writer.Write(p)
	}
	if length > this.bufWriter.Available() {
		if err := this.bufWriter.Flush(); err != nil {
			return 0, err
		}
	}
	return this.bufWriter.Write(p)
}

func (this *bufferedWriter) Flush() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.bufWriter.Flush()
}

func (this *bufferedWriter) Close() error {
	if closer, ok := this.writer.(io.Closer); ok {
		return closer.Close()
	}
	this.cancel()
	return nil
}

func (this *bufferedWriter) flushPeriodically() {
	if this.flushPeriod == 0 {
		return
	}
	ticker := time.NewTicker(this.flushPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			this.Flush()
		case <-ctx.Done():
			return
		}
	}
}
