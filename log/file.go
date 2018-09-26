package log

import (
	"fmt"
	"os"
	"time"
)

var pid = os.Getpid()

type file struct {
	f *os.File
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
	this.f.Sync()
	this.f.Close()
}
