package network

import (
	"bufio"
	"net"
	"sync"
)

type TCPConnection struct {
	conn    *net.TCPConn
	handler TCPHandler
	codec   Codec

	mu          sync.Mutex
	cond        *sync.Cond
	sendBuffers [][]byte
	closed      bool
}

func newTCPConnection(conn *net.TCPConn, handler TCPHandler) *TCPConnection {
	connection := &TCPConnection{
		conn:    conn,
		handler: handler,
	}
	connection.cond = sync.NewCond(&connection.mu)
	return connection
}

func (this *TCPConnection) serve() {
	defer this.conn.Close()

	func() {
		defer this.conn.Close()

		w := bufio.NewWriter(this.conn)
		var buffers [][]byte
		for {
			this.mu.Lock()
			for len(this.sendBuffers) == 0 {
				if this.closed {
					this.mu.Unlock()
					return
				}
				this.cond.Wait()
			}
			this.sendBuffers, buffers = buffers, this.sendBuffers // swap
			this.mu.Unlock()

			for _, b := range buffers {
				if err := this.codec.Write(w, b); err != nil {
					this.Close()
					return
				}
			}
			if err := w.Flush(); err != nil {
				this.Close()
				return
			}
			buffers = buffers[0:0] // clear buffers
		}
	}()

	rd := bufio.NewReader(this.conn)
	for {
		if b, err := this.codec.Read(rd); err == nil {
			this.handler.Receive(this, b)
		} else {
			//
			this.Close()
			break
		}
	}
}

func (this *TCPConnection) Send(b []byte) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.closed {
		return
	}
	this.sendBuffers = append(this.sendBuffers, b)
	this.cond.Signal()
}

func (this *TCPConnection) Close() {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.closed {
		return
	}
	this.closed = true
	this.cond.Signal()
}
