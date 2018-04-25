package network

import (
	"log"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	addr        string
	handler     TCPHandler
	mu          sync.Mutex
	connections map[*TCPConnection]struct{}
	listener    *net.TCPListener
	done        chan struct{}
}

func NewTCPServer(addr string, handler TCPHandler) *TCPServer {
	server := &TCPServer{
		addr:        addr,
		handler:     handler,
		connections: make(map[*TCPConnection]struct{}),
	}
	return server
}

func listenTCP(addr string) (*net.TCPListener, error) {
	if addr == "" {
		addr = ":0"
	}
	la, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", la)
}

func (this *TCPServer) Start() error {
	if this.listener != nil {
		return nil
	}
	ln, err := listenTCP(this.addr)
	if err != nil {
		return err
	}
	this.listener = ln

	done := make(chan struct{})
	this.done = done

	go func() {
		defer ln.Close()
		defer close(done)

		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					log.Printf("TCPServer: accept error: %v; retrying in %v", err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				log.Printf("TCPServer: error: %v", err)
				return
			}
			tempDelay = 0
			conn.Close()
			//connection := newTcpConnection(conn, this.Handler)
			//		connection.onClose = func() {
			//		this.trackConnection(connection, false)
			//}
			//this.trackConnection(connection, true)
			//		go connection.serve(this.Handler)
		}
	}()
	return nil
}

func (this *TCPServer) addConnection(connection *TCPConnection) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.connections[connection] = struct{}{}
}

func (this *TCPServer) removeConnection(connection *TCPConnection) {
	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.connections, connection)
}

func (this *TCPServer) Close() error {
	var err error
	if ln := this.listener; ln != nil {
		this.listener = nil
		err = ln.Close()
		<-this.done
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	for connection := range this.connections {
		// connection.conn.Close()
		delete(this.connections, connection)
	}
	return err
}
