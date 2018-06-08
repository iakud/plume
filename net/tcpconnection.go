package net

import (
	"bufio"
	"net"
	"sync"

	"github.com/iakud/falcon/codec"
)

type TCPConnection struct {
	conn *net.TCPConn

	bufs   [][]byte
	mutex  sync.Mutex
	cond   *sync.Cond
	closed bool
}

func newTCPConnection(conn *net.TCPConn) *TCPConnection {
	connection := &TCPConnection{
		conn: conn,
	}
	connection.cond = sync.NewCond(&connection.mutex)
	conn.SetNoDelay(true) // no delay
	return connection
}

func (this *TCPConnection) ServeCodec(codec codec.Codec, connectFunc, disconnectFunc func(*TCPConnection), receiveFunc func(*TCPConnection, []byte)) {
	this.startBackgroundWrite(codec) // start write
	defer this.stopBackgroundWrite() // stop write

	if connectFunc != nil {
		connectFunc(this)
	}
	this.loopRead(codec, receiveFunc) // loop read
	if disconnectFunc != nil {
		disconnectFunc(this)
	}
}

func (this *TCPConnection) loopRead(codec codec.Codec, receiveFunc func(*TCPConnection, []byte)) {
	defer this.conn.Close()

	rd := bufio.NewReader(this.conn)
	for {
		b, err := codec.Read(rd)
		if err != nil {
			return
		}
		if receiveFunc != nil {
			receiveFunc(this, b)
		}
	}
}

func (this *TCPConnection) startBackgroundWrite(codec codec.Codec) {
	go this.backgroundWrite(codec)
}

func (this *TCPConnection) backgroundWrite(codec codec.Codec) {
	defer this.conn.Close()

	w := bufio.NewWriter(this.conn)
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
			if err := codec.Write(w, b); err != nil {
				this.stopBackgroundWrite()
				return
			}
		}
		if err := w.Flush(); err != nil {
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
