package network

import (
	"sync"
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
