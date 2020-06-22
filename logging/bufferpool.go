package logging

import (
	"sync"
)

type bufferPool struct {
	p *sync.Pool
}

func newBufferPool() *bufferPool {
	pool := &bufferPool{
		p: &sync.Pool{
			New: func() interface{} {
				return newBuffer()
			},
		},
	}
	return pool
}

func (p *bufferPool) get() *buffer {
	buf := p.p.Get().(*buffer)
	buf.reset()
	return buf
}

func (p *bufferPool) put(buf *buffer) {
	p.p.Put(buf)
}
