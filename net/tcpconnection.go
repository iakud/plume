package net

import (
	"bufio"
	"net"
	"sync"
)

type TCPConnection struct {
	conn  *net.TCPConn
	codec Codec

	connectFunc    func(*TCPConnection)
	disconnectFunc func(*TCPConnection)
	receiveFunc    func(*TCPConnection, []byte)

	bufs   [][]byte
	mutex  sync.Mutex
	cond   *sync.Cond
	closed bool
}

func newTCPConnection(conn *net.TCPConn, codec Codec) *TCPConnection {
	connection := &TCPConnection{
		conn:  conn,
		codec: codec,
	}
	connection.cond = sync.NewCond(&connection.mutex)
	conn.SetNoDelay(true) // no delay

	connection.startBackgroundWrite() // start write
	return connection
}

func (this *TCPConnection) serve() {
	if this.connectFunc != nil {
		this.connectFunc(this)
	}
	this.loopRead() // loop read
	if this.disconnectFunc != nil {
		this.disconnectFunc(this)
	}
}

func (this *TCPConnection) loopRead() {
	defer this.conn.Close()

	r := bufio.NewReader(this.conn)
	for {
		b, err := this.codec.Read(r)
		if err != nil {
			return
		}
		if this.receiveFunc != nil {
			this.receiveFunc(this, b)
		}
	}
}

func (this *TCPConnection) startBackgroundWrite() {
	go this.backgroundWrite()
}

func (this *TCPConnection) backgroundWrite() {
	defer this.conn.Close()

	w := bufio.NewWriter(this.conn)
	for {
		bufs, closed := this.waitForBuffers()
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
		if closed {
			return
		}
	}
}

func (this *TCPConnection) waitForBuffers() ([][]byte, bool) {
	var bufs [][]byte
	var closed bool
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for !this.closed && len(this.bufs) == 0 {
		this.cond.Wait()
	}
	bufs, this.bufs = this.bufs, bufs // swap
	closed = this.closed
	return bufs, closed
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

func (this *TCPConnection) close() {
	this.conn.SetLinger(0)
	this.conn.Close()
	this.stopBackgroundWrite() // stop write
}

func (this *TCPConnection) Shutdown() {
	this.stopBackgroundWrite() // stop write
}
