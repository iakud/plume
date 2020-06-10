package logging

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	p *sync.Pool
}

func NewBufferPool() *BufferPool {
	pool := &BufferPool{
		p: &sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
	return pool
}

func (p *BufferPool) Get() *bytes.Buffer {
	buffer := p.p.Get().(*bytes.Buffer)
	buffer.Reset()
	return buffer
}

func (p *BufferPool) Put(buffer *bytes.Buffer) {
	p.p.Put(buffer)
}
