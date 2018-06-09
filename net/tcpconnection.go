package net

import (
	"bufio"
	"net"
	"sync"

	"github.com/iakud/falcon"
)

type TCPConnection struct {
	loop *falcon.EventLoop
	conn *net.TCPConn

	dec    *Decoder
	enc    *Encoder
	decBuf *bufio.Reader
	encBuf *bufio.Writer

	bufs   [][]byte
	mutex  sync.Mutex
	cond   *sync.Cond
	closed bool
}

func newTCPConnection(loop *falcon.EventLoop, conn *net.TCPConn) *TCPConnection {
	decBuf := bufio.NewReader(conn)
	encBuf := bufio.NewWriter(conn)
	connection := &TCPConnection{
		loop:   loop,
		conn:   conn,
		dec:    NewDecoder(decBuf),
		enc:    NewEncoder(encBuf),
		decBuf: decBuf,
		encBuf: encBuf,
	}
	connection.cond = sync.NewCond(&connection.mutex)
	conn.SetNoDelay(true) // no delay

	connection.startBackgroundWrite() // start write
	return connection
}

func (this *TCPConnection) serve(connectFunc, disconnectFunc func(*TCPConnection), receiveFunc func(*TCPConnection, []byte)) {
	if connectFunc != nil {
		this.loop.RunInLoop(func() { connectFunc(this) })
	}
	if receiveFunc != nil {
		this.loopRead(receiveFunc) // loop read
	}
	if disconnectFunc != nil {
		this.loop.RunInLoop(func() { disconnectFunc(this) })
	}
}

func (this *TCPConnection) loopRead(receiveFunc func(*TCPConnection, []byte)) {
	defer this.conn.Close()

	for {
		b, err := this.dec.Decode()
		if err != nil {
			return
		}
		this.loop.RunInLoop(func() { receiveFunc(this, b) })
	}
}

func (this *TCPConnection) startBackgroundWrite() {
	go this.backgroundWrite()
}

func (this *TCPConnection) backgroundWrite() {
	defer this.conn.Close()

	for {
		var bufs [][]byte
		var closed bool
		// wait bufs
		this.mutex.Lock()
		for !this.closed && len(this.bufs) == 0 {
			this.cond.Wait()
		}
		bufs, this.bufs = this.bufs, nil // swap
		closed = this.closed
		this.mutex.Unlock()

		for _, b := range bufs {
			if err := this.enc.Encode(b); err != nil {
				this.stopBackgroundWrite()
				return
			}
		}
		if err := this.encBuf.Flush(); err != nil {
			this.stopBackgroundWrite()
			return
		}
		if closed {
			return
		}
	}
}

func (this *TCPConnection) stopBackgroundWrite() {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	this.closed = true
	this.mutex.Unlock()

	this.cond.Signal()
}

func (this *TCPConnection) Send(b []byte) {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	this.bufs = append(this.bufs, b)
	this.mutex.Unlock()

	this.cond.Signal()
}

func (this *TCPConnection) close() {
	this.conn.SetLinger(0)
	this.conn.Close()
	this.stopBackgroundWrite() // stop write
}

func (this *TCPConnection) Shutdown() {
	this.stopBackgroundWrite() // stop write
}
