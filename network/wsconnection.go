package network

import (
	"errors"
	"net"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

var (
	ErrWSConnectionPendingSendFull = errors.New("network: WebSocket connection pending send full")
)

type WSConnection struct {
	conn *websocket.Conn

	bufs        [][]byte
	pendingSend int
	mutex       sync.Mutex
	cond        *sync.Cond
	closed      bool
}

func newWSConnection(conn *websocket.Conn) *WSConnection {
	connection := &WSConnection{conn: conn}
	connection.cond = sync.NewCond(&connection.mutex)
	return connection
}

func (c *WSConnection) serve(handler WSHandler) {
	defer c.conn.Close()

	// start write
	c.startBackgroundWrite()
	defer c.stopBackgroundWrite()

	// conn event
	handler.Connect(c, true)
	defer handler.Connect(c, false)
	for {
		var data []byte
		if err := websocket.Message.Receive(c.conn, &data); err != nil {
			c.conn.Close()
			break
		}
		handler.Receive(c, data)
	}
}

func (c *WSConnection) startBackgroundWrite() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return
	}
	go c.backgroundWrite()
}

func (c *WSConnection) backgroundWrite() {
	for closed := false; !closed; {
		var bufs [][]byte

		c.mutex.Lock()
		for !c.closed && len(c.bufs) == 0 {
			c.cond.Wait()
		}
		bufs, c.bufs = c.bufs, bufs // swap
		closed = c.closed
		c.mutex.Unlock()

		for _, message := range bufs {
			if err := websocket.Message.Send(c.conn, message); err != nil {
				c.closeWrite()
				c.conn.Close()
				return
			}
		}
	}
	// not writing now
	c.conn.Close()
}

func (c *WSConnection) stopBackgroundWrite() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	c.cond.Signal()
}

func (c *WSConnection) closeWrite() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return
	}
	c.closed = true
}

func (c *WSConnection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *WSConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *WSConnection) SetPendingSend(pendingSend int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pendingSend = pendingSend
}

func (c *WSConnection) Send(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return nil
	}
	if c.pendingSend > 0 && len(c.bufs) >= c.pendingSend {
		return ErrWSConnectionPendingSendFull
	}
	c.bufs = append(c.bufs, data)
	c.cond.Signal()
	return nil
}

func (c *WSConnection) Shutdown() {
	c.stopBackgroundWrite()
}

func (c *WSConnection) Close() {
	c.conn.Close()
}
func (c *WSConnection) CloseWithTimeout(timeout time.Duration) {
	time.AfterFunc(timeout, c.Close)
}
