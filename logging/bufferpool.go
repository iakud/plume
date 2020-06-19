package logging

import (
	"sync"
)

type BufferPool struct {
	p *sync.Pool
}

func NewBufferPool() *BufferPool {
	pool := &BufferPool{
		p: &sync.Pool{
			New: func() interface{} {
				return newBuffer()
			},
		},
	}
	return pool
}

func (p *BufferPool) Get() *buffer {
	buf := p.p.Get().(*buffer)
	buf.Reset()
	return buf
}

func (p *BufferPool) Put(buf *buffer) {
	p.p.Put(buf)
}
