package network

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"sync"
)

type TCPConnection struct {
	conn *net.TCPConn

	connectFunc    func(*TCPConnection)
	disconnectFunc func(*TCPConnection)
	receiveFunc    func(*TCPConnection, []byte)

	mu      sync.Mutex
	cond    *sync.Cond
	buffers [][]byte
	closed  bool
}

func newTCPConnection(conn *net.TCPConn) *TCPConnection {
	connection := &TCPConnection{
		conn: conn,
	}
	conn.SetNoDelay(true) // no delay
	connection.cond = sync.NewCond(&connection.mu)
	return connection
}

func (this *TCPConnection) serve() {
	defer this.conn.Close()

	// on connect
	if this.connectFunc != nil {
		this.connectFunc(this)
	}

	done := make(chan struct{})
	wait := func() {
		<-done
	}
	go this.serveWrite(done)

	this.serveRead()

	wait() // wait write

	// on disconnect
	if this.disconnectFunc != nil {
		this.disconnectFunc(this)
	}
}

func (this *TCPConnection) serveRead() {
	defer this.conn.Close()
	defer this.closeWrite()

	rd := bufio.NewReader(this.conn)
	h := make([]byte, 2)
	for {
		if _, err := io.ReadFull(rd, h); err != nil {
			return
		}
		n := binary.BigEndian.Uint16(h)
		b := make([]byte, n)
		if _, err := io.ReadFull(rd, b); err != nil {
			return
		}
		if this.receiveFunc != nil {
			this.receiveFunc(this, b)
		}
	}
}

func (this *TCPConnection) serveWrite(done chan struct{}) {
	defer close(done)
	defer this.conn.Close()

	w := bufio.NewWriter(this.conn)
	h := make([]byte, 2)
	for {
		var buffers [][]byte
		this.mu.Lock()
		for len(this.buffers) == 0 {
			if this.closed {
				this.mu.Unlock()
				return
			}
			this.cond.Wait()
		}
		this.buffers, buffers = buffers, this.buffers // swap
		this.mu.Unlock()

		if err := func() error {
			for _, b := range buffers {
				binary.BigEndian.PutUint16(h, uint16(len(b)))
				if _, err := w.Write(h); err != nil {
					return err
				}
				if _, err := w.Write(b); err != nil {
					return err
				}
			}
			return w.Flush()
		}(); err != nil {
			this.closeSend()
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
	this.conn.SetLinger(0)
	this.conn.Close()
	this.closeWrite()
}

func (this *TCPConnection) closeWrite() {
	this.mu.Lock()
	if this.closed {
		this.mu.Unlock()
		return
	}
	this.closed = true
	this.mu.Unlock()

	this.cond.Signal()
}

func (this *TCPConnection) closeSend() {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.closed {
		return
	}
	this.closed = true
}

func (this *TCPConnection) Shutdown() {
	this.closeWrite()
}
