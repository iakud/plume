package logging

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

	maxRolls   int
	history    []string
	filePeriod time.Time
}

// maxRools: if <= 0, unlimited
func NewFileWriter(path string, period time.Duration, maxRolls int) *FileWriter {
	ctx, cancel := context.WithCancel(context.Background())
	dir, name := filepath.Split(path)
	fw := &FileWriter{
		dir:    filepath.Dir(dir),
		name:   name,
		period: period,
		cancel: cancel,

		maxRolls: maxRolls,
	}
	if maxRolls > 0 {
		if history, err := fw.historyRolls(); err == nil {
			fw.history = history
		}
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
	if fw.maxRolls > 0 {
		if stat, err := file.Stat(); err == nil {
			fw.history = append(fw.history, stat.Name())
			fw.removeOldRolls()
		}
	}
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
	if err := os.Symlink(filename, symlink); err != nil {
		os.Link(filename, symlink)
	}
	return file, nil
}

func (fw *FileWriter) historyRolls() ([]string, error) {
	f, err := os.Open(fw.dir)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	var history []string
	for _, file := range files {
		// regular files
		if file.Mode()&os.ModeType != 0 {
			continue
		}
		// filter
		if strings.HasPrefix(file.Name(), fw.name) {
			history = append(history, file.Name())
		}
	}
	sort.Sort(sort.StringSlice(history)) // sort by string
	return history, nil
}

func (fw *FileWriter) removeOldRolls() {
	if nRolls := len(fw.history); nRolls > fw.maxRolls {
		removeRools := nRolls - fw.maxRolls
		for _, name := range fw.history[:removeRools] {
			filename := filepath.Join(fw.dir, name)
			os.Remove(filename) // ignore err
		}
		fw.history = fw.history[removeRools:]
	}
}
