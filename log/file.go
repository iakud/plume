package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

var pid = os.Getpid()

type File struct {
	f *os.File
}

func (this *File) Write() {

}

func (this *File) rollFile() {

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
