package net

import (
	"bufio"
	"net"
	"sync"
)

type TCPHandler interface {
	Connect(*TCPConnection)
	Disconnect(*TCPConnection)
	Receive(*TCPConnection, []byte)
}

type defaultTCPHandler struct {
}

func (*defaultTCPHandler) Connect(*TCPConnection) {

}

func (*defaultTCPHandler) Disconnect(*TCPConnection) {

}

func (*defaultTCPHandler) Receive(*TCPConnection, []byte) {

}

var DefaultTCPHandler *defaultTCPHandler = &defaultTCPHandler{}

type TCPConnection struct {
	conn *net.TCPConn

	handler TCPHandler
	codec   Codec

	bufs   [][]byte
	mutex  sync.Mutex
	cond   *sync.Cond
	closed bool

	Userdata interface{}
}

func newTCPConnection(conn *net.TCPConn, handler TCPHandler, codec Codec) *TCPConnection {
	connection := &TCPConnection{
		conn:    conn,
		handler: handler,
		codec:   codec,
	}
	connection.cond = sync.NewCond(&connection.mutex)
	conn.SetNoDelay(true) // no delay

	return connection
}

func (this *TCPConnection) serve() {
	this.startBackgroundWrite()
	defer this.stopBackgroundWrite()

	this.handler.Connect(this)
	this.loopRead() // loop read
	this.handler.Disconnect(this)
}

func (this *TCPConnection) loopRead() {
	defer this.conn.Close()

	r := bufio.NewReader(this.conn)
	for {
		b, err := this.codec.Read(r)
		if err != nil {
			return
		}
		this.handler.Receive(this, b)
	}
}

func (this *TCPConnection) startBackgroundWrite() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return
	}
	go this.backgroundWrite()
}

func (this *TCPConnection) backgroundWrite() {
	defer this.conn.Close()

	w := bufio.NewWriter(this.conn)
	for closed := false; !closed; {
		var bufs [][]byte

		this.mutex.Lock()
		for !this.closed && len(this.bufs) == 0 {
			this.cond.Wait()
		}
		bufs, this.bufs = this.bufs, bufs // swap
		closed = this.closed
		this.mutex.Unlock()

		for _, b := range bufs {
			if err := this.codec.Write(w, b); err != nil {
				this.closeSend()
				return
			}
		}
		if err := w.Flush(); err != nil {
			this.closeSend()
			return
		}
	}
}

func (this *TCPConnection) stopBackgroundWrite() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return
	}
	this.closed = true
	this.cond.Signal()
}

func (this *TCPConnection) closeSend() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return
	}
	this.closed = true
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

func (this *TCPConnection) close() error {
	this.conn.SetLinger(0)
	return this.conn.Close()
}

func (this *TCPConnection) Shutdown() {
	this.stopBackgroundWrite() // stop write
}
