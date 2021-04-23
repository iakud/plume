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
	ErrConnectionPendingSendFull = errors.New("network: Connection pending send full")
)

type TCPConnection struct {
	conn *net.TCPConn

	bufs        [][]byte
	pendingSend int
	mutex       sync.Mutex
	cond        *sync.Cond
	closed      bool

	Userdata interface{}
}

func newTCPConnection(conn *net.TCPConn) *TCPConnection {
	connection := &TCPConnection{
		conn: conn,
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
			this.Close()
		}
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
			this.Close()
			return
		}
		handler.Receive(this, b)
	}
}

func (this *TCPConnection) startBackgroundWrite(codec Codec) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		this.mutex.Unlock()
		return
	}
	go this.backgroundWrite(codec)
}

func (this *TCPConnection) backgroundWrite(codec Codec) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("network: panic serving %v: %v\n%s", this.RemoteAddr(), err, buf)
			this.Close()
		}
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
				this.Close()
				return
			}
		}
		if err := w.Flush(); err != nil {
			this.closeWrite()
			this.Close()
			return
		}
	}
	// not writing now
	this.conn.CloseWrite() // only SHUT_WR
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

func (this *TCPConnection) SetPendingSend(pendingSend int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.pendingSend = pendingSend
}

func (this *TCPConnection) Send(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.closed {
		return nil
	}
	if this.pendingSend > 0 && len(this.bufs) >= this.pendingSend {
		return ErrConnectionPendingSendFull
	}
	this.bufs = append(this.bufs, b)
	this.cond.Signal()
	return nil
}

func (this *TCPConnection) Close() {
	this.conn.Close()
}

func (this *TCPConnection) Shutdown() {
	const delay = time.Second * 3
	this.ShutdownIn(delay)
}

func (this *TCPConnection) ShutdownIn(d time.Duration) {
	this.stopBackgroundWrite() // stop write
	time.AfterFunc(d, this.Close)
}
