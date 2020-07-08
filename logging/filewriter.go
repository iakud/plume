package logging

import (
	"bufio"
	"context"
	"errors"
	"fmt"
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
	period time.Duration
	cancel context.CancelFunc

	mutex  sync.Mutex
	file   *os.File
	buffer *bufio.Writer
	closed bool

	filePeriod time.Time
}

func NewFileWriter(path string, period time.Duration) *FileWriter {
	ctx, cancel := context.WithCancel(context.Background())
	dir, name := filepath.Split(path)
	fw := &FileWriter{
		dir:    dir,
		name:   name,
		period: period,
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
	thisPeriod := time.Now().Truncate(fw.period)
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
	if fw.buffer == nil {
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
	if fw.buffer == nil {
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
	if fw.buffer == nil {
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
	if fw.buffer != nil {
		fw.buffer.Flush()
		fw.buffer = nil
		fw.file.Close()
		fw.file = nil
	}
	file, err := fw.createFile(t)
	if err != nil {
		return err
	}
	fw.file = file
	fw.buffer = bufio.NewWriterSize(file, kBufferSize)
	return nil
}

func (fw *FileWriter) createFile(t time.Time) (*os.File, error) {
	name := fmt.Sprintf("%s.%04d%02d%02d-%02d%02d%02d", fw.name,
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	filename := filepath.Join(fw.dir, name)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("logging: cannot create log: %v", err)
	}

	symlink := filepath.Join(fw.dir, fw.name)
	os.Remove(symlink) // ignore err
	if err := os.Symlink(name, symlink); err != nil {
		os.Link(name, symlink)
	}
	return file, nil
}
