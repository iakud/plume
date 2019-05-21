package net

import (
	"context"
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
	ctx    context.Context
	cancel context.CancelFunc
	addr   string

	mutex      sync.Mutex
	connection *TCPConnection
	closed     bool
}

func NewTCPClient(addr string) *TCPClient {
	ctx, cancel := context.WithCancel(context.Background())
	client := &TCPClient{
		ctx:    ctx,
		cancel: cancel,
		addr:   addr,
	}
	return client
}

func dialTCPContext(ctx context.Context, addr string) (*net.TCPConn, error) {
	var dialer net.Dialer
	c, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, ok := c.(*net.TCPConn)
	if !ok {
		c.Close() // close
		return nil, errors.New("unexpected type")
	}
	return conn, nil
}

func (this *TCPClient) DialAndServe(handler TCPHandler, codec Codec) error {
	if this.isClosed() {
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
		conn, err := dialTCPContext(this.ctx, this.addr)
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
	this.cancel()

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
