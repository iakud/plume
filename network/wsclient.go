package network

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrWSClientClosed = errors.New("network: Websocket client closed")
)

type WSClient struct {
	Url     string
	Handler WSHandler
	retry   bool

	mutex      sync.Mutex
	connection *WSConnection
	closed     bool
}

func NewWSClient(url string, handler WSHandler) *WSClient {
	client := &WSClient{
		Url:     url,
		Handler: handler,
	}
	return client
}

func (c *WSClient) EnableRetry()  { c.retry = true }
func (c *WSClient) DisableRetry() { c.retry = false }

func DialAndServeWS(url string, handler WSHandler) error {
	client := &WSClient{Url: url, Handler: handler}
	return client.DialAndServe()
}

func (c *WSClient) DialAndServe() error {
	if c.isClosed() {
		return ErrWSClientClosed
	}

	handler := c.Handler
	if handler == nil {
		handler = DefaultWSHandler
	}

	var tempDelay time.Duration // how long to sleep on connect failure
	for {
		conn, _, err := websocket.DefaultDialer.Dial(c.Url, nil)
		if err != nil {
			if c.isClosed() {
				return ErrWSClientClosed
			}
			if !c.retry {
				return err
			}

			if tempDelay == 0 {
				tempDelay = 1 * time.Second
			} else {
				tempDelay *= 2
			}
			if max := 1 * time.Minute; tempDelay > max {
				tempDelay = max
			}
			log.Printf("network: Websocket client dial error: %v; retrying in %v", err, tempDelay)
			time.Sleep(tempDelay)
			continue
		}
		tempDelay = 0

		connection := newWSConnection(conn)
		if err := c.newConnection(connection); err != nil {
			connection.Close()
			return err
		}
		if err := c.serveConnection(connection, handler); err != nil {
			return err
		}
	}
}

func (c *WSClient) isClosed() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.closed
}

func (c *WSClient) newConnection(connection *WSConnection) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return ErrWSClientClosed
	}
	c.connection = connection
	return nil
}

func (c *WSClient) serveConnection(conn *WSConnection, handler WSHandler) error {
	conn.serve(handler)
	// remove connection
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return ErrWSClientClosed
	}
	c.connection = nil
	return nil
}

func (c *WSClient) GetConnection() *WSConnection {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return nil
	}
	return c.connection
}

func (c *WSClient) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return
	}
	c.closed = true
	if c.connection == nil {
		return
	}
	c.connection.Close()
	c.connection = nil
}
