package net

import (
	"log"
	"testing"

	"github.com/iakud/falcon/codec"
)

type EchoServer struct {
	server *TCPServer
}

func NewEchoServer(addr string) *EchoServer {
	server := NewTCPServer(addr)
	echoServer := &EchoServer{
		server: server,
	}
	return echoServer
}

func (this *EchoServer) Start() {
	this.server.Start(this.newConnection)
}

func (this *EchoServer) Close() {
	this.server.Close()
}

func (this *EchoServer) newConnection(connection *TCPConnection) {
	connection.ServeCodec(&codec.StdCodec{}, this.onConnect, this.onDisconnect, this.onReceive)
}

func (this *EchoServer) onConnect(connection *TCPConnection) {
	log.Println("server: connected.")
}

func (this *EchoServer) onDisconnect(connection *TCPConnection) {
	log.Println("server: disconnected.")
}

func (this *EchoServer) onReceive(connection *TCPConnection, b []byte) {
	message := string(b)
	log.Println("server: receive", message)
	connection.Send(b)
	connection.Shutdown()
}

var Message string = "hello"

type EchoClient struct {
	client *TCPClient
	done   chan struct{}
}

func NewEchoClient(addr string) *EchoClient {
	client := NewTCPClient(addr)
	echoClient := &EchoClient{
		client: client,
		done:   make(chan struct{}),
	}
	return echoClient
}

func (this *EchoClient) Start() {
	this.client.Start(this.newConnection)
}

func (this *EchoClient) Done() {
	<-this.done
}

func (this *EchoClient) newConnection(connection *TCPConnection) {
	connection.ServeCodec(&codec.StdCodec{}, this.onConnect, this.onDisconnect, this.onReceive)
}

func (this *EchoClient) onConnect(connection *TCPConnection) {
	log.Println("client: connected.")
	log.Println("client: send", Message)
	connection.Send([]byte(Message))
}

func (this *EchoClient) onDisconnect(connection *TCPConnection) {
	log.Println("client: disconnected.")
	this.client.Close()
	close(this.done)
}

func (this *EchoClient) onReceive(connection *TCPConnection, b []byte) {
	log.Println("client: receive ", string(b))
}

func TestEcho(t *testing.T) {
	echoServer := NewEchoServer("localhost:8000")
	echoServer.Start()

	echoClient := NewEchoClient("localhost:8000")
	echoClient.Start()

	echoClient.Done()
	echoServer.Close()
}
