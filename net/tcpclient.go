package net

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var (
	ErrClientClosed = errors.New("net: Client closed")
)

type TCPClient struct {
	addr    string
	handler TCPHandler
	codec   Codec

	mutex      sync.Mutex
	connection *TCPConnection
	closed     bool
}

func NewTCPClient(addr string, handler TCPHandler, codec Codec) *TCPClient {
	client := &TCPClient{
		addr:    addr,
		handler: handler,
		codec:   codec,
	}
	return client
}

func dialTCP(addr string) (*net.TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return net.DialTCP("tcp", nil, raddr)
}

func (this *TCPClient) ConnectAndServe() error {
	if this.isClosed() {
		return ErrClientClosed
	}

	handler := this.handler
	if handler == nil {
		handler = DefaultTCPHandler
	}
	codec := this.codec
	if codec == nil {
		codec = DefaultCodec
	}

	var tempDelay time.Duration // how long to sleep on connect failure
	for {
		conn, err := dialTCP(this.addr)
		if err != nil {
			if this.isClosed() {
				return ErrClientClosed
			}

			if tempDelay == 0 {
				tempDelay = 1 * time.Second
			} else {
				tempDelay *= 2
			}
			if max := 1 * time.Minute; tempDelay > max {
				tempDelay = max
			}
			log.Printf("TCPClient: connect error: %v; retrying in %v", err, tempDelay)
			time.Sleep(tempDelay)
			continue
		}
		tempDelay = 0

		connection := newTCPConnection(conn, handler, codec)
		if err := this.newConnection(connection); err != nil {
			connection.close()
			return err
		}
		if err := this.serveConnection(connection); err != nil {
			return err
		}
	}
}

func (this *TCPClient) isClosed() bool {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.closed
}

func (this *TCPClient) newConnection(connection *TCPConnection) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return ErrClientClosed
	}
	this.connection = connection
	return nil
}

func (this *TCPClient) serveConnection(connection *TCPConnection) error {
	connection.serve()
	// remove connection
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return ErrClientClosed
	}
	this.connection = nil
	return nil
}

func (this *TCPClient) GetConnection() *TCPConnection {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return nil
	}
	return this.connection
}

func (this *TCPClient) Close() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.closed {
		return
	}
	this.closed = true
	if this.connection == nil {
		return
	}
	this.connection.close()
	this.connection = nil
}
