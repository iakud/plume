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

var DefaultTCPHandler defaultTCPHandler

type TCPConnectionWriter struct {
	conn *net.TCPConn

	bufs    [][]byte
	mutex   sync.Mutex
	cond    *sync.Cond
	started bool
	closed  bool
}

func newTCPConnectionWriter(conn *net.TCPConn) *TCPConnectionWriter {
	writer := &TCPConnectionWriter{conn: conn}
	writer.cond = sync.NewCond(&writer.mutex)
	return writer
}

func (this *TCPConnectionWriter) startBackgroundWrite(cw CodecWriter) {
	this.mutex.Lock()
	if this.started {
		this.mutex.Unlock()
		return
	}
	this.started = true
	this.mutex.Unlock()

	go this.backgroundWrite(cw)
}

func (this *TCPConnectionWriter) backgroundWrite(cw CodecWriter) {
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
			if err := cw.Write(w, b); err != nil {
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

func (this *TCPConnectionWriter) stopBackgroundWrite() {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	this.closed = true
	this.mutex.Unlock()

	this.cond.Signal()
}

func (this *TCPConnectionWriter) closeSend() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return
	}
	this.closed = true
}

func (this *TCPConnectionWriter) Send(b []byte) {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	this.bufs = append(this.bufs, b)
	this.mutex.Unlock()

	this.cond.Signal()
}

type TCPConnection struct {
	conn  *net.TCPConn
	codec Codec

	w *TCPConnectionWriter

	Userdata interface{}
}

func newTCPConnection(conn *net.TCPConn, codec Codec) *TCPConnection {
	connection := &TCPConnection{
		conn:  conn,
		codec: codec,
	}
	connection.cond = sync.NewCond(&connection.mutex)
	conn.SetNoDelay(true) // no delay

	return connection
}

func (this *TCPConnection) serve(handler TCPHandler) {
	if handler == nil {
		handler = DefaultTCPHandler
	}

	this.w = newTCPConnectionWriter(this.conn)
	this.startBackgroundWrite(this.codec)

	handler.Connect(this)

	r := bufio.NewReader(this.conn)
	for {
		b, err := this.codec.Read(r)
		if err != nil {
			return
		}
		handler.Receive(this, b)
	}
	this.conn.Close()

	this.loopRead(handler) // loop read
	this.stopBackgroundWrite()
	handler.Disconnect(this)
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
