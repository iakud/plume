package logging

import (
	"fmt"
	"testing"
	"time"
)

func TestFileWriter(t *testing.T) {
	fw := NewFileWriter(t.Name()+".log", time.Hour, 0)
	defer fw.Close()
	s := fmt.Sprintf("open file: %s", time.Now())
	fw.Write([]byte(s))
}

func TestFileWriterFlush(t *testing.T) {
	fw := NewFileWriter(t.Name()+".log", time.Hour, 0)
	s := fmt.Sprintf("open file: %s\n", time.Now())
	fw.Write([]byte(s))
	fw.Write([]byte("flush\n"))
	fw.Flush()
	fw.Write([]byte("after flush\n"))
}

func TestFileWriterRolls(t *testing.T) {
	maxRolls := 2
	fw := NewFileWriter(t.Name()+".log", time.Second, maxRolls)
	defer fw.Close()
	fw.Write([]byte("test 1"))
	time.Sleep(time.Second)
	fw.Write([]byte("test 2"))
	time.Sleep(time.Second)
	fw.Write([]byte("test 3"))
}
