package logging

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var ErrClosed = errors.New("logging: file writer already closed")

const kBufferSize = 256 * 1024
const kFlushInterval = 10 * time.Second

type FileWriter struct {
	dir    string
	name   string
	cancel context.CancelFunc

	mutex  sync.Mutex
	file   *os.File
	buffer *bufio.Writer
	closed bool

	filePeriod time.Time
}

func NewFileWriter(path string) *FileWriter {
	ctx, cancel := context.WithCancel(context.Background())
	dir, name := filepath.Split(path)
	fw := &FileWriter{
		dir:    dir,
		name:   name,
		cancel: cancel,
	}

	go fw.flushPeriodically(ctx)
	return fw
}

func (fw *FileWriter) Write(p []byte) (int, error) {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()
	if fw.closed {
		return 0, ErrClosed
	}
	thisPeriod := time.Now().Truncate(time.Hour)
	if thisPeriod != fw.filePeriod {
		if err := fw.rollFile(thisPeriod); err != nil {
			return 0, err
		}
		fw.filePeriod = thisPeriod
	}
	return fw.buffer.Write(p)
}

func (fw *FileWriter) Sync() error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()
	if fw.closed {
		return ErrClosed
	}
	if fw.file == nil {
		return nil
	}
	fw.buffer.Flush()
	return fw.file.Sync()
}

func (fw *FileWriter) Flush() error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()
	if fw.closed {
		return ErrClosed
	}
	if fw.file == nil {
		return nil
	}
	return fw.buffer.Flush()
}

func (fw *FileWriter) Close() error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()
	if fw.closed {
		return ErrClosed
	}
	fw.closed = true
	fw.cancel()
	if fw.file == nil {
		return nil
	}
	fw.buffer.Flush()
	return fw.file.Close()
}

func (fw *FileWriter) flushPeriodically(ctx context.Context) {
	ticker := time.NewTicker(kFlushInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fw.Flush()
		case <-ctx.Done():
			return
		}
	}
}

func (fw *FileWriter) rollFile(t time.Time) error {
	if fw.file != nil {
		fw.buffer.Flush()
		fw.file.Close()
		fw.file = nil
	}
	file, err := createFile(fw.dir, fw.name, t)
	if err != nil {
		return err
	}
	fw.file = file
	if fw.buffer == nil {
		fw.buffer = bufio.NewWriterSize(file, kBufferSize)
	} else {
		fw.buffer.Reset(file)
	}
	return nil
}
