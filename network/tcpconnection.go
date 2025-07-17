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

func (c *TCPConnection) serve(handler TCPHandler, codec Codec) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("network: panic serving %v: %v\n%s", c.RemoteAddr(), err, buf)
			c.Close()
		}
	}()

	// start write
	c.startBackgroundWrite(codec)
	defer c.stopBackgroundWrite()
	// conn event
	handler.Connect(c, true)
	defer handler.Connect(c, false)
	// loop read
	r := bufio.NewReader(c.conn)
	for {
		b, err := codec.Read(r)
		if err != nil {
			c.Close()
			return
		}
		handler.Receive(c, b)
	}
}

func (c *TCPConnection) startBackgroundWrite(codec Codec) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return
	}
	go c.backgroundWrite(codec)
}

func (c *TCPConnection) backgroundWrite(codec Codec) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("network: panic serving %v: %v\n%s", c.RemoteAddr(), err, buf)
			c.Close()
		}
	}()

	// loop write
	w := bufio.NewWriter(c.conn)
	for closed := false; !closed; {
		var bufs [][]byte

		c.mutex.Lock()
		for !c.closed && len(c.bufs) == 0 {
			c.cond.Wait()
		}
		bufs, c.bufs = c.bufs, bufs // swap
		closed = c.closed
		c.mutex.Unlock()

		for _, b := range bufs {
			if err := codec.Write(w, b); err != nil {
				c.closeWrite()
				c.Close()
				return
			}
		}
		if err := w.Flush(); err != nil {
			c.closeWrite()
			c.Close()
			return
		}
	}
	// not writing now
	c.conn.CloseWrite() // only SHUT_WR
}

func (c *TCPConnection) stopBackgroundWrite() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	c.cond.Signal()
}

func (c *TCPConnection) closeWrite() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return
	}
	c.closed = true
}

func (c *TCPConnection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *TCPConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *TCPConnection) SetNoDelay(noDelay bool) error {
	return c.conn.SetNoDelay(noDelay)
}

func (c *TCPConnection) SetPendingSend(pendingSend int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pendingSend = pendingSend
}

func (c *TCPConnection) Send(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return nil
	}
	if c.pendingSend > 0 && len(c.bufs) >= c.pendingSend {
		return ErrConnectionPendingSendFull
	}
	c.bufs = append(c.bufs, b)
	c.cond.Signal()
	return nil
}

func (c *TCPConnection) Shutdown() {
	c.stopBackgroundWrite() // stop write
}

func (c *TCPConnection) Close() {
	c.conn.Close()
}

func (c *TCPConnection) CloseWithTimeout(timeout time.Duration) {
	time.AfterFunc(timeout, c.Close)
}
