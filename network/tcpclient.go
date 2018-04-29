package network

import (
	"log"
	"net"
	"sync"
	"time"
)

type TCPClient struct {
	addr    string
	handler TCPHandler

	mu         sync.Mutex
	connection *TCPConnection
	done       chan struct{}
}

func NewTCPClient(addr string, handler TCPHandler) *TCPClient {
	client := &TCPClient{
		addr:    addr,
		handler: handler,
	}
	return client
}

func (this *TCPClient) connect() error {
	raddr, err := net.ResolveTCPAddr("tcp", this.addr)
	if err != nil {
		return err
	}

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := net.DialTCP("tcp", nil, raddr)
			if err != nil {
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
			connection := newTCPConnection(conn, this.handler)
			this.setConnection(connection)
			time.Sleep(10 * time.Second)
			this.removeConnection()
		}
	}()
	return nil
}

func (this *TCPClient) setConnection(connection *TCPConnection) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.connection = connection
}

func (this *TCPClient) removeConnection() {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.connection = nil
}

func (this *TCPClient) GetConnection() *TCPConnection {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.connection
}

func (this *TCPClient) startConnection(conn *net.TCPConn) {

}
