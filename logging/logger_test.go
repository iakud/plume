package logging

import (
	"log"
	"testing"
)

func TestLog(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("123 %d", 111)
}

func TestLogger(t *testing.T) {
	logger := New()
	logger.Debugf("%s%d\n", "gda", 123)
}
