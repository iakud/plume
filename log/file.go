package log

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var pid = os.Getpid()

const bufferSize = 256 * 1024

type file struct {
	w *bufio.Writer
	f *os.File

	nbytes uint64 // The number of bytes written to this file
}

func (this *file) Write() {

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
