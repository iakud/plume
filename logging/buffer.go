package logging

import (
	"bytes"
	"time"
)

type buffer struct {
	*bytes.Buffer
}

func newBuffer() *buffer {
	b := &buffer{
		Buffer: &bytes.Buffer{},
	}
	return b
}

func (b *buffer) formatHeader(now time.Time, l Level) {

}
