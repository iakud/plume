package logging

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const kBufferSize = 256 * 1024
const kFlushInterval = 3 * time.Second

type bufferedWriter struct {
	mutex     sync.Mutex
	bufWriter *bufio.Writer
	path      string
	name      string
	file      *os.File

	ctx    context.Context
	cancel context.CancelFunc
}

func newBufferedWriter(path string, name string) *bufferedWriter {
	ctx, cancel := context.WithCancel(context.Background())
	bufferedWriter := &bufferedWriter{
		path: path,
		name: name,

		ctx:    ctx,
		cancel: cancel,
	}

	go bufferedWriter.flushPeriodically()
	return bufferedWriter
}

func (this *bufferedWriter) Write(p []byte) (int, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.file == nil {
		this.createFile(filepath.Join(this.path, this.name))
	}
	return this.bufWriter.Write(p)
}

func (this *bufferedWriter) Flush() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.bufWriter.Flush()
}

func (this *bufferedWriter) Sync() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.bufWriter.Flush()
	// return this.file.Sync()
}

func (this *bufferedWriter) Close() error {
	this.Flush()
	this.cancel()
	return nil
}

func (this *bufferedWriter) flushPeriodically() {
	ticker := time.NewTicker(kFlushInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			this.Flush()
		case <-this.ctx.Done():
			return
		}
	}
}

func (this *bufferedWriter) createFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	this.file = file
	if this.bufWriter == nil {
		this.bufWriter = bufio.NewWriterSize(file, kBufferSize)
	} else {
		this.bufWriter.Reset(file)
	}
	return nil
}
