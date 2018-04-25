package network

import (
	"net"
	"sync"
)

type TCPConnection struct {
	conn          *net.TCPConn
	handler       TCPHandler
	mu            sync.Mutex
	cond          *sync.Cond
	inWait        int
	sendBytes     [][]byte
	sendClosed    bool
	receiveClosed bool
	closed        bool
	doneChan      chan struct{}

	maxLen int
}

func newTCPConnection(conn *net.TCPConn, handler TCPHandler) *TCPConnection {
	connection := &TCPConnection{
		conn:     conn,
		handler:  handler,
		doneChan: make(chan struct{}),

		maxLen: 4096,
	}
	connection.cond = sync.NewCond(&connection.mu)
	return connection
}
