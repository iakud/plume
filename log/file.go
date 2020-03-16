package log

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

var pid = os.Getpid()

const bufferSize = 256 * 1024

const (
	defaultFilePermissions      = 0666
	defaultDirectoryPermissions = 0767
)

type file struct {
	w *bufio.Writer

	dirPath string

	file     *os.File
	fileSize int64 // The number of bytes written to this file
}

func newFile(path string) io.WriteCloser {

}

func (this *file) Write(b []byte) (int, error) {

	if this.f == nil {
		if err := this.createFile(); err != nil {
			return 0, err
		}
	}
	n, err := this.f.Write(b)
	return n, err
}

func (this *file) Close() error {
	if this.f != nil {
		if err := this.f.Close(); err != nil {
			return err
		}
		this.f = nil
	}
	return nil
}

func (this *file) rollFile() {

}

func logName(name string, t time.Time) string {
	return fmt.Sprintf("%s.%04d%02d%02d-%02d%02d%02d.%d.log",
		name,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		pid)
}

func (this *file) rotate() {
	if this.f != nil {
		this.w.Flush()
		this.f.Close()
	}
}

func (this *file) createFile() error {
	if len(this.currentDirPath) != 0 {
		if err := os.MkdirAll(rw.currentDirPath, defaultDirectoryPermissions); err != nil {
			return err
		}
	}
	rollname := time.Now().Format("20060102-150405")
	filePath := filepath.Join(currentDirPath, rollname)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, defaultFilePermissions)
	if err != nil {
		return err
	}
	stat, err := f.Stat()
	if err != nil {
		rw.currentFile.Close()
		rw.currentFile = nil
		return err
	}
	this.file = file
	this.fileSize = stat.Size()
	return nil
}
