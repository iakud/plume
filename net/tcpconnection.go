package net

import (
	"bufio"
	"net"
	"sync"

	"github.com/iakud/falcon"
)

type TCPConnection struct {
	loop  *falcon.EventLoop
	conn  *net.TCPConn
	codec Codec

	connectFunc    func(*TCPConnection)
	disconnectFunc func(*TCPConnection)
	receiveFunc    func(*TCPConnection, []byte)

	bufs   [][]byte
	mutex  sync.Mutex
	cond   *sync.Cond
	closed bool

	Userdata interface{}
}

func newTCPConnection(loop *falcon.EventLoop, conn *net.TCPConn, codec Codec) *TCPConnection {
	connection := &TCPConnection{
		loop:  loop,
		conn:  conn,
		codec: codec,
	}
	connection.cond = sync.NewCond(&connection.mutex)
	conn.SetNoDelay(true) // no delay

	connection.startBackgroundWrite() // start write
	return connection
}

func (this *TCPConnection) serve() {
	this.onConnect()
	this.loopRead() // loop read
	this.onDisconnect()
}

func (this *TCPConnection) loopRead() {
	defer this.conn.Close()

	r := bufio.NewReader(this.conn)
	for {
		b, err := this.codec.Read(r)
		if err != nil {
			return
		}
		this.onReceive(b)
	}
}

func (this *TCPConnection) startBackgroundWrite() {
	go this.backgroundWrite()
}

func (this *TCPConnection) backgroundWrite() {
	defer this.conn.Close()

	w := bufio.NewWriter(this.conn)
	for {
		var bufs [][]byte
		var closed bool

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

func (this *TCPConnection) onConnect() {
	if this.connectFunc == nil {
		return
	}
	if this.loop == nil {
		this.connectFunc(this)
		return
	}
	this.loop.RunInLoop(func() { this.connectFunc(this) })
}

func (this *TCPConnection) onDisconnect() {
	if this.disconnectFunc == nil {
		return
	}
	if this.loop == nil {
		this.disconnectFunc(this)
		return
	}
	this.loop.RunInLoop(func() { this.disconnectFunc(this) })
}

func (this *TCPConnection) onReceive(b []byte) {
	if this.receiveFunc == nil {
		return
	}
	if this.loop == nil {
		this.receiveFunc(this, b)
		return
	}
	this.loop.RunInLoop(func() { this.receiveFunc(this, b) })
}
