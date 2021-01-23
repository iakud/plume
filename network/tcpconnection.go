package network

import (
	"bufio"
	"errors"
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

var (
	ErrConnectionClosed = errors.New("network: Connection closed")
	ErrSendFull         = errors.New("network: Connection send full")
	kPendingSendSize    = 128
)

type TCPConnection struct {
	conn *net.TCPConn

	bufs   [][]byte
	mutex  sync.Mutex
	cond   *sync.Cond
	closed bool

	Userdata interface{}
}

func newTCPConnection(conn *net.TCPConn) *TCPConnection {
	connection := &TCPConnection{
		conn: conn,
	}
	connection.cond = sync.NewCond(&connection.mutex)
	conn.SetNoDelay(true) // no delay

	return connection
}

func (this *TCPConnection) serve(handler TCPHandler, codec Codec) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("network: panic serving %v: %v\n%s", this.RemoteAddr(), err, buf)
		}
		this.conn.Close()
	}()

	// start write
	this.startBackgroundWrite(codec)
	defer this.stopBackgroundWrite()
	// conn event
	handler.Connect(this, true)
	defer handler.Connect(this, false)
	// loop read
	r := bufio.NewReader(this.conn)
	for {
		b, err := codec.Read(r)
		if err != nil {
			return
		}
		handler.Receive(this, b)
	}
}

func (this *TCPConnection) startBackgroundWrite(codec Codec) {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	this.mutex.Unlock()

	go this.backgroundWrite(codec)
}

func (this *TCPConnection) backgroundWrite(codec Codec) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("network: panic serving %v: %v\n%s", this.RemoteAddr(), err, buf)
		}
		this.conn.Close()
	}()

	// loop write
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
			if err := codec.Write(w, b); err != nil {
				this.closeWrite()
				return
			}
		}
		if err := w.Flush(); err != nil {
			this.closeWrite()
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

func (this *TCPConnection) closeWrite() {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	this.closed = true
	this.mutex.Unlock()
}

func (this *TCPConnection) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *TCPConnection) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *TCPConnection) Send(b []byte) error {
	this.mutex.Lock()
	if this.closed {
		this.mutex.Unlock()
		return ErrConnectionClosed
	}
	if len(this.bufs) >= kPendingSendSize {
		this.mutex.Unlock()
		return ErrSendFull
	}
	this.bufs = append(this.bufs, b)
	this.mutex.Unlock()

	this.cond.Signal()
	return nil
}

func (this *TCPConnection) close() {
	this.conn.Close()
}

func (this *TCPConnection) Shutdown() {
	this.stopBackgroundWrite() // stop write
}

func (this *TCPConnection) ForceClose() {
	this.conn.SetLinger(0)
	this.conn.Close()
}

func (this *TCPConnection) AfterForceClose(d time.Duration) {
	time.AfterFunc(d, this.ForceClose)
}
