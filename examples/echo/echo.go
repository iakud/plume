package echo

import (
	"log"

	"github.com/iakud/falcon/network"
)

type EchoServer struct {
	server *network.TCPServer
}

func NewEchoServer(addr string) *EchoServer {
	server := network.NewTCPServer(addr)
	echoServer := &EchoServer{
		server: server,
	}
	server.ConnectFunc = echoServer.onConnect
	server.DisconnectFunc = echoServer.onDisconnect
	server.ReceiveFunc = echoServer.onReceive
	return echoServer
}

func (this *EchoServer) Start() error {
	return this.server.Start()
}

func (this *EchoServer) Close() error {
	return this.server.Close()
}

func (this *EchoServer) onConnect(connection *network.TCPConnection) {
	log.Println("server: on connect.")
}

func (this *EchoServer) onDisconnect(connection *network.TCPConnection) {
	log.Println("server: on disconnect.")
}

func (this *EchoServer) onReceive(connection *network.TCPConnection, b []byte) {
	log.Println("server: on receive", string(b))
	log.Println("server: send", string(b))
	connection.Send(b)
	connection.Shutdown()
}

type EchoClient struct {
	client *network.TCPClient
	done   chan struct{}
}

func NewEchoClient(addr string) *EchoClient {
	client := network.NewTCPClient(addr)
	echoClient := &EchoClient{
		client: client,
		done:   make(chan struct{}),
	}
	client.ConnectFunc = echoClient.onConnect
	client.DisconnectFunc = echoClient.onDisconnect
	client.ReceiveFunc = echoClient.onReceive
	return echoClient
}

func (this *EchoClient) Start() error {
	return this.client.Start()
}

func (this *EchoClient) Wait() {
	<-this.done
}

func (this *EchoClient) onConnect(connection *network.TCPConnection) {
	log.Println("client: on connect.")
	message := "hello"
	log.Println("client: send", message)
	connection.Send([]byte(message))
}

func (this *EchoClient) onDisconnect(connection *network.TCPConnection) {
	log.Println("client: on disconnect.")
	this.client.Close()
	log.Println("client: closed.")
	close(this.done)
}

func (this *EchoClient) onReceive(connection *network.TCPConnection, b []byte) {
	log.Println("client: on receive ", string(b))
	connection.Shutdown()
}
