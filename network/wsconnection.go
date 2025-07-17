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

type WSConnection struct {
	conn *websocket.Conn

	bufs        []WSMessage
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

func (c *WSConnection) Send(messageType int, data []byte) error {
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

func (c *WSConnection) Shutdown() {
	c.stopBackgroundWrite()
}

func (c *WSConnection) Close() {
	c.conn.Close()
}
func (c *WSConnection) CloseWithTimeout(timeout time.Duration) {
	time.AfterFunc(timeout, c.Close)
}
