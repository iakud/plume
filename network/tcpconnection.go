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
	ErrConnectionClosed        = errors.New("network: Connection closed")
	ErrConnectionHighWaterMark = errors.New("network: Connection high watermark")
	DefaultHighWaterMark       = 64 * 1024 * 1024
)

type TCPConnection struct {
	conn *net.TCPConn

	bufs          [][]byte
	pendingWrite  int
	highWaterMark int
	mutex         sync.Mutex
	cond          *sync.Cond
	closed        bool

	Userdata interface{}
}

func newTCPConnection(conn *net.TCPConn) *TCPConnection {
	connection := &TCPConnection{
		conn:          conn,
		highWaterMark: DefaultHighWaterMark,
	}
	connection.cond = sync.NewCond(&connection.mutex)
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
		this.pendingWrite = 0             // clear
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
	defer this.mutex.Unlock()
	if this.closed {
		return
	}
	this.closed = true
	this.cond.Signal()
}

func (this *TCPConnection) closeWrite() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return
	}
	this.closed = true
}

func (this *TCPConnection) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *TCPConnection) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *TCPConnection) SetNoDelay(noDelay bool) error {
	return this.conn.SetNoDelay(noDelay)
}

func (this *TCPConnection) SetHighWaterMark(highWaterMark int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.highWaterMark = highWaterMark
}

func (this *TCPConnection) Send(b []byte) error {
	n := len(b)
	if n == 0 {
		return nil
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return ErrConnectionClosed
	}
	var err error
	if this.pendingWrite+n >= this.highWaterMark && this.pendingWrite < this.highWaterMark {
		err = ErrConnectionHighWaterMark
	}
	this.bufs = append(this.bufs, b)
	this.pendingWrite += n
	this.cond.Signal()
	return err
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
