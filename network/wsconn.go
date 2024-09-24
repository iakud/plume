package network

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

var (
	ErrWSConnectionPendingSendFull = errors.New("network: WebSocket connection pending send full")
)

type WSMessageType int

const (
	TextMessage   WSMessageType = websocket.TextMessage
	BinaryMessage               = websocket.BinaryMessage
)

type WSMessage struct {
	MessageType int
	Data        []byte
}

type WSConn struct {
	wsconn *websocket.Conn

	bufs        []WSMessage
	pendingSend int
	mutex       sync.Mutex
	cond        *sync.Cond
	closed      bool
}

func newWSConn(wsconn *websocket.Conn) *WSConn {
	conn := &WSConn{wsconn: wsconn}
	conn.cond = sync.NewCond(&conn.mutex)
	return conn
}

func (c *WSConn) serve(handler WSHandler) {
	defer c.wsconn.Close()

	// start write
	c.startBackgroundWrite()
	defer c.stopBackgroundWrite()

	// conn event
	handler.Connect(c, true)
	defer handler.Connect(c, false)
	for {
		messageType, data, err := c.wsconn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			c.wsconn.Close()
			break
		}
		handler.Receive(c, WSMessageType(messageType), data)
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
			if err := c.wsconn.WriteMessage(message.MessageType, message.Data); err != nil {
				c.closeWrite()
				c.wsconn.Close()
				return
			}
		}
	}
	// not writing now
	c.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
	return c.wsconn.LocalAddr()
}

func (c *WSConn) RemoteAddr() net.Addr {
	return c.wsconn.RemoteAddr()
}

func (c *WSConn) SetPendingSend(pendingSend int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pendingSend = pendingSend
}

func (c *WSConn) Send(messageType WSMessageType, data []byte) error {
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
	c.bufs = append(c.bufs, WSMessage{int(messageType), data})
	c.cond.Signal()
	return nil
}

func (c *WSConn) Shutdown() {
	c.stopBackgroundWrite()
}

func (c *WSConn) Close() {
	c.wsconn.Close()
}
func (c *WSConn) CloseWithTimeout(timeout time.Duration) {
	time.AfterFunc(timeout, c.Close)
}
