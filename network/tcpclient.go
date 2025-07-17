package network

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var (
	ErrClientClosed = errors.New("network: Client closed")
)

type TCPClient struct {
	addr  string
	retry bool

	mutex      sync.Mutex
	connection *TCPConnection
	closed     bool
}

func NewTCPClient(addr string) *TCPClient {
	client := &TCPClient{
		addr:  addr,
		retry: false,
	}
	return client
}

func (c *TCPClient) EnableRetry()  { c.retry = true }
func (c *TCPClient) DisableRetry() { c.retry = false }

func dialTCP(addr string) (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return net.DialTCP("tcp", nil, raddr)
}

func (c *TCPClient) DialAndServe(handler TCPHandler, codec Codec) error {
	if c.isClosed() {
		return ErrClientClosed
	}

	if handler == nil {
		handler = DefaultTCPHandler
	}
	if codec == nil {
		codec = DefaultCodec
	}

	var tempDelay time.Duration // how long to sleep on connect failure
	for {
		conn, err := dialTCP(c.addr)
		if err != nil {
			if c.isClosed() {
				return ErrClientClosed
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
			log.Printf("network: TCPClient dial error: %v; retrying in %v", err, tempDelay)
			time.Sleep(tempDelay)
			continue
		}
		tempDelay = 0

		connection := newTCPConnection(conn)
		if err := c.newConnection(connection); err != nil {
			connection.Close()
			return err
		}
		if err := c.serveConnection(connection, handler, codec); err != nil {
			return err
		}
	}
}

func (c *TCPClient) isClosed() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.closed
}

func (c *TCPClient) newConnection(connection *TCPConnection) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return ErrClientClosed
	}
	c.connection = connection
	return nil
}

func (c *TCPClient) serveConnection(connection *TCPConnection, handler TCPHandler, codec Codec) error {
	connection.serve(handler, codec)
	// remove connection
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return ErrClientClosed
	}
	c.connection = nil
	return nil
}

func (c *TCPClient) GetConnection() *TCPConnection {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return nil
	}
	return c.connection
}

func (c *TCPClient) Close() {
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
