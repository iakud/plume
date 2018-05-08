package network

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"sync"
)

type TCPConnection struct {
	conn    *net.TCPConn
	handler TCPHandler

	mu      sync.Mutex
	cond    *sync.Cond
	buffers [][]byte
	closed  bool
}

func newTCPConnection(conn *net.TCPConn, handler TCPHandler) *TCPConnection {
	connection := &TCPConnection{
		conn:    conn,
		handler: handler,
	}
	conn.SetNoDelay(true) // no delay
	connection.cond = sync.NewCond(&connection.mu)
	return connection
}

func (this *TCPConnection) serve() {
	defer this.conn.Close()

	done := make(chan struct{})
	go this.serveSend(done)

	this.handler.Connected(this)

	rd := bufio.NewReader(this.conn)
	h := make([]byte, 2)
	for {
		if _, err := io.ReadFull(rd, h); err != nil {
			this.Close()
			break
		}
		n := binary.BigEndian.Uint16(h)
		b := make([]byte, n)
		if _, err := io.ReadFull(rd, b); err != nil {
			this.Close()
			break
		}
		this.handler.Receive(this, b)
	}
	<-done // wait done

	this.handler.Disconnected(this)
}

func (this *TCPConnection) serveSend(done chan struct{}) {
	defer close(done)
	defer this.conn.Close()

	w := bufio.NewWriter(this.conn)
	h := make([]byte, 2)
	for {
		var buffers [][]byte
		if this.wait(&buffers); len(buffers) == 0 {
			return
		}

		for _, b := range buffers {
			binary.BigEndian.PutUint16(h, uint16(len(b)))
			if _, err := w.Write(h); err != nil {
				this.close()
				return
			}
			if _, err := w.Write(b); err != nil {
				this.close()
				return
			}
		}
		if err := w.Flush(); err != nil {
			this.close()
			return
		}
	}
}

func (this *TCPConnection) wait(buffers *[][]byte) {
	this.mu.Lock()
	defer this.mu.Unlock()
	for len(this.buffers) == 0 {
		if this.closed {
			return
		}
		this.cond.Wait()
	}
	this.buffers, *buffers = *buffers, this.buffers // swap
}

func (this *TCPConnection) Send(b []byte) {
	this.mu.Lock()
	if this.closed {
		this.mu.Unlock()
		return
	}
	this.buffers = append(this.buffers, b)
	this.mu.Unlock()
	this.cond.Signal()
}

func (this *TCPConnection) close() {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.closed {
		return
	}
	this.closed = true
}

func (this *TCPConnection) Close() {
	this.mu.Lock()
	if this.closed {
		this.mu.Unlock()
		return
	}
	this.closed = true
	this.mu.Unlock()
	this.cond.Signal()
}
