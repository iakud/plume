package echo

import (
	"fmt"
	//	"time"

	"github.com/iakud/falcon/network"
)

type EchoServer struct {
	server *network.TCPServer
}

func NewEchoServer(addr string) *EchoServer {
	echoServer := &EchoServer{
		server: network.NewTCPServer(addr),
	}
	return echoServer
}

func (this *EchoServer) Start() error {
	return this.server.Start(this)
}

func (this *EchoServer) Close() error {
	return this.server.Close()
}

func (this *EchoServer) Connected(connection *network.TCPConnection) {
	fmt.Println("server: connected.")
}

func (this *EchoServer) Disconnected(connection *network.TCPConnection) {
	fmt.Println("server: disconnected.")
}

func (this *EchoServer) Receive(connection *network.TCPConnection, b []byte) {
	fmt.Println("server: receive", string(b))
	fmt.Println("server: send", string(b))
	connection.Send(b)
	connection.Close()
}

type EchoClient struct {
	client *network.TCPClient
	done   chan struct{}
}

func NewEchoClient(addr string) *EchoClient {
	echoClient := &EchoClient{
		client: network.NewTCPClient(addr),
		done:   make(chan struct{}),
	}
	return echoClient
}

func (this *EchoClient) Start() error {
	return this.client.Start(this)
}

func (this *EchoClient) Done() <-chan struct{} {
	return this.done
}

func (this *EchoClient) Connected(connection *network.TCPConnection) {
	fmt.Println("client: connected.")
	message := "hello"
	fmt.Println("client: send", message)
	connection.Send([]byte(message))

}

func (this *EchoClient) Disconnected(connection *network.TCPConnection) {
	fmt.Println("client: disconnected.")
	this.client.Close()
	close(this.done)
}

func (this *EchoClient) Receive(connection *network.TCPConnection, b []byte) {
	fmt.Println("client: receive ", string(b))
	connection.Close()
}
