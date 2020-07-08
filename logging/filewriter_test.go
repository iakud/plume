package logging

import (
	"fmt"
	"testing"
	"time"
)

func TestFileWriter(t *testing.T) {
	fw := NewFileWriter("test.log")
	defer fw.Close()
	s := fmt.Sprintf("open file: %s", time.Now())
	fw.Write([]byte(s))
}

func TestFileWriterFlush(t *testing.T) {
	fw := NewFileWriter("testflush.log")
	s := fmt.Sprintf("open file: %s\n", time.Now())
	fw.Write([]byte(s))
	fw.Write([]byte("flush\n"))
	fw.Flush()
	fw.Write([]byte("after flush\n"))
}
