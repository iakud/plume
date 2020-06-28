package logging

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var ErrClosed = errors.New("logging: file writer already closed")

const kBufferSize = 256 * 1024
const kFlushInterval = 10 * time.Second

type fileWriter struct {
	mutex  sync.Mutex
	buffer *bufio.Writer
	closed bool

	dir  string
	name string
	file *os.File

	filePeriod time.Time
}

func NewFileWriter(path string) *fileWriter {
	dir, name := filepath.Split(path)
	fileWriter := &fileWriter{
		dir:  dir,
		name: name,
	}

	go fileWriter.flushPeriodically()
	return fileWriter
}

func (this *fileWriter) Write(p []byte) (int, error) {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return 0, ErrClosed
	}
	thisPeriod := time.Now().Truncate(time.Hour)
	if thisPeriod != this.filePeriod {
		if err := this.rollFile(thisPeriod); err != nil {
			this.mutex.Unlock()
			return 0, err
		}
	}
	n, err := this.buffer.Write(p)
	this.mutex.Unlock()
	return n, err
}

func (this *fileWriter) Sync() error {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return ErrClosed
	}
	err := this.buffer.Flush()
	this.mutex.Unlock()
	return err
}

func (this *fileWriter) Close() error {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return ErrClosed
	}
	if this.file != nil {
		this.buffer.Flush()
		this.file.Close()
		this.file = nil
	}
	this.closed = true
	this.mutex.Unlock()
	return nil
}

func (this *fileWriter) flushPeriodically() {
	ticker := time.NewTicker(kFlushInterval)
	defer ticker.Stop()
	for _ = range ticker.C {
		err := this.Sync()
		if err != nil && err == ErrClosed {
			return
		}
	}
}

func (this *fileWriter) rollFile(t time.Time) error {
	if this.file != nil {
		this.buffer.Flush()
		this.file.Close()
		this.file = nil
	}
	file, err := createFile(this.dir, this.name, t)
	if err != nil {
		return err
	}
	this.file = file
	this.buffer = bufio.NewWriterSize(file, kBufferSize)
	this.filePeriod = t
	return nil
}
