package network

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrWSConnectionPendingSendFull = errors.New("network: WebSocket connection pending send full")
)

const (
	TextMessage   = websocket.TextMessage
	BinaryMessage = websocket.BinaryMessage
)

type WSMessage struct {
	Type int
	Data []byte
}

type WSConn struct {
	conn *websocket.Conn

	bufs        []WSMessage
	pendingSend int
	mutex       sync.Mutex
	cond        *sync.Cond
	closed      bool
}

func newWSConn(wsconn *websocket.Conn) *WSConn {
	conn := &WSConn{conn: wsconn}
	conn.cond = sync.NewCond(&conn.mutex)
	return conn
}

func (c *WSConn) serve(handler WSHandler) {
	defer c.conn.Close()

	// start write
	c.startBackgroundWrite()
	defer c.stopBackgroundWrite()

	// conn event
	handler.Connect(c, true)
	defer handler.Connect(c, false)
	for {
		messageType, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			c.conn.Close()
			break
		}
		handler.Receive(c, messageType, data)
	}
}

func (c *WSConn) startBackgroundWrite() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return
	}
	go c.backgroundWrite()
}

func (c *WSConn) backgroundWrite() {
	for closed := false; !closed; {
		var bufs []WSMessage

		c.mutex.Lock()
		for !c.closed && len(c.bufs) == 0 {
			c.cond.Wait()
		}
		bufs, c.bufs = c.bufs, bufs // swap
		closed = c.closed
		c.mutex.Unlock()

		for _, message := range bufs {
			if err := c.conn.WriteMessage(message.Type, message.Data); err != nil {
				c.closeWrite()
				c.conn.Close()
				return
			}
		}
	}
	// not writing now
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (c *WSConn) stopBackgroundWrite() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	c.cond.Signal()
}

func (c *WSConn) closeWrite() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.closed {
		return
	}
	c.closed = true
}

func (c *WSConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *WSConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *WSConn) SetPendingSend(pendingSend int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pendingSend = pendingSend
}

func (c *WSConn) Send(messageType int, data []byte) error {
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
	c.bufs = append(c.bufs, WSMessage{messageType, data})
	c.cond.Signal()
	return nil
}

func (c *WSConn) Shutdown() {
	c.stopBackgroundWrite()
}

func (c *WSConn) Close() {
	c.conn.Close()
}
func (c *WSConn) CloseWithTimeout(timeout time.Duration) {
	time.AfterFunc(timeout, c.Close)
}
