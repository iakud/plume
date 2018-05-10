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

	mu      sync.Mutex
	cond    *sync.Cond
	buffers [][]byte
	closed  bool

	connectFunc    func(*TCPConnection)
	disconnectFunc func(*TCPConnection)
	receiveFunc    func(*TCPConnection, []byte)
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

	this.onConnect()         // on connect
	defer this.onDisonnect() // on disconnect

	done := make(chan struct{})
	go this.serveWrite(done) // write

	wait := func() { <-done }
	defer wait() // wait write

	this.serveRead() // read
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
		this.onReceive(b)
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
		// swap
		buffers, this.buffers = this.buffers, buffers
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

func (this *TCPConnection) closeSend() {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.closed {
		return
	}
	this.closed = true
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

func (this *TCPConnection) close() {
	this.conn.SetLinger(0)
	this.conn.Close()
	this.closeWrite()
}

func (this *TCPConnection) Shutdown() {
	this.closeWrite()
}

func (this *TCPConnection) onConnect() {
	if this.connectFunc != nil {
		this.connectFunc(this)
	}
}

func (this *TCPConnection) onDisonnect() {
	if this.disconnectFunc != nil {
		this.disconnectFunc(this)
	}
}

func (this *TCPConnection) onReceive(b []byte) {
	if this.receiveFunc != nil {
		this.receiveFunc(this, b)
	}
}
