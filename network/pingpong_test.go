package network

import (
	"log"
	"sync/atomic"
	"testing"
	"time"
)

const (
	kBlockSize   = 1024 * 16
	kClientCount = 10
	kTimeout     = time.Second * 10
)

// Pingpong Server
type pingpongServer struct {
	server *TCPServer
}

func newPingpongServer(addr string) *pingpongServer {
	srv := &pingpongServer{
		server: NewTCPServer(addr),
	}
	return srv
}

func (srv *pingpongServer) ListenAndServe() {
	if err := srv.server.ListenAndServe(srv, nil); err != nil {
		if err == ErrServerClosed {
			return
		}
		log.Println(err)
	}
}

func (srv *pingpongServer) Close() {
	srv.server.Close()
}

func (srv *pingpongServer) Connect(connection *TCPConnection, connected bool) {
	if connected {
		connection.SetNoDelay(true)
	}
}

func (srv *pingpongServer) Receive(connection *TCPConnection, b []byte) {
	connection.Send(b)
}

// Pingpong Client
type pingpongClient struct {
	clients    []*TCPClient
	message    []byte
	nConnected int32

	bytesRead    int64
	messagesRead int64
	done         chan struct{}
}

func newPingpongClient(addr string) *pingpongClient {
	// build message
	message := make([]byte, kBlockSize)
	for i := 0; i < kBlockSize; i++ {
		message[i] = byte(i % 128)
	}
	c := &pingpongClient{
		message: message,
		done:    make(chan struct{}),
	}
	clients := make([]*TCPClient, kClientCount)
	for i := 0; i < kClientCount; i++ {
		client := NewTCPClient(addr)
		go c.serveClient(client)
		clients[i] = client
	}
	c.clients = clients
	time.AfterFunc(kTimeout, c.handleTimeout)
	return c
}

func (c *pingpongClient) serveClient(client *TCPClient) {
	client.EnableRetry() // 启用retry
	if err := client.DialAndServe(c, nil); err != nil {
		if err == ErrClientClosed {
			return
		}
		log.Println(err)
	}
}

func (c *pingpongClient) handleTimeout() {
	for _, client := range c.clients {
		client.Close()
	}
}

func (c *pingpongClient) Connect(connection *TCPConnection, connected bool) {
	if connected {
		connection.SetNoDelay(true)
		connection.Send(c.message)
		if atomic.AddInt32(&c.nConnected, 1) != kClientCount {
			return
		}
		log.Println("all connected")
	} else {
		if atomic.AddInt32(&c.nConnected, -1) != 0 {
			return
		}
		bytesRead := atomic.LoadInt64(&c.bytesRead)
		messagesRead := atomic.LoadInt64(&c.messagesRead)
		log.Println(bytesRead, "total bytes read")
		log.Println(messagesRead, "total messages read")
		log.Println(bytesRead/messagesRead, "average message size")
		timeout := int64(kTimeout / time.Second)
		log.Println(bytesRead/(timeout*1024*1024), "MiB/s throughput")
		close(c.done)
	}
}

func (c *pingpongClient) Receive(connection *TCPConnection, b []byte) {
	connection.Send(b)
	atomic.AddInt64(&c.messagesRead, 1)
	atomic.AddInt64(&c.bytesRead, int64(len(b)))
}

func (c *pingpongClient) Done() {
	<-c.done
}

func TestPingpong(t *testing.T) {
	srv := newPingpongServer("localhost:8000")
	go srv.ListenAndServe()
	defer srv.Close()
	c := newPingpongClient("localhost:8000")
	c.Done()
}
