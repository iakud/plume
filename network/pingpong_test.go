package network_test

import (
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/iakud/plume/network"
)

const (
	kBlockSize   = 1024 * 16
	kClientCount = 10
	kTimeout     = time.Second * 10
)

// Pingpong Server
type pingpongServer struct {
	server *network.TCPServer
}

func newPingpongServer(addr string) *pingpongServer {
	srv := &pingpongServer{
		server: network.NewTCPServer(addr),
	}
	return srv
}

func (srv *pingpongServer) ListenAndServe() {
	if err := srv.server.ListenAndServe(srv, nil); err != nil {
		if err == network.ErrServerClosed {
			return
		}
		slog.Info(err.Error())
	}
}

func (srv *pingpongServer) Close() {
	srv.server.Close()
}

func (srv *pingpongServer) Connect(connection *network.TCPConnection, connected bool) {
	if connected {
		connection.SetNoDelay(true)
	}
}

func (srv *pingpongServer) Receive(connection *network.TCPConnection, b []byte) {
	connection.Send(b)
}

// Pingpong Client
type pingpongClient struct {
	clients    []*network.TCPClient
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
	clients := make([]*network.TCPClient, kClientCount)
	for i := 0; i < kClientCount; i++ {
		client := network.NewTCPClient(addr)
		go c.serveClient(client)
		clients[i] = client
	}
	c.clients = clients
	time.AfterFunc(kTimeout, c.handleTimeout)
	return c
}

func (c *pingpongClient) serveClient(client *network.TCPClient) {
	client.EnableRetry() // 启用retry
	if err := client.DialAndServe(c, nil); err != nil {
		if err == network.ErrClientClosed {
			return
		}
		slog.Error(err.Error())
	}
}

func (c *pingpongClient) handleTimeout() {
	for _, client := range c.clients {
		client.Close()
	}
}

func (c *pingpongClient) Connect(connection *network.TCPConnection, connected bool) {
	if connected {
		connection.SetNoDelay(true)
		connection.Send(c.message)
		if atomic.AddInt32(&c.nConnected, 1) != kClientCount {
			return
		}
		slog.Info("all connected")
	} else {
		if atomic.AddInt32(&c.nConnected, -1) != 0 {
			return
		}
		bytesRead := atomic.LoadInt64(&c.bytesRead)
		messagesRead := atomic.LoadInt64(&c.messagesRead)
		slog.Info("", "total bytes read", bytesRead)
		slog.Info("", "total messages read", messagesRead)
		slog.Info("", "average message size", bytesRead/messagesRead)
		timeout := int64(kTimeout / time.Second)
		slog.Info("", "MiB/s throughput", bytesRead/(timeout*1024*1024))
		close(c.done)
	}
}

func (c *pingpongClient) Receive(connection *network.TCPConnection, b []byte) {
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
