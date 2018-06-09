package net

import (
	"log"
	"testing"

	"github.com/iakud/falcon"
)

type EchoServer struct {
	loop   *falcon.EventLoop
	server *TCPServer
}

func NewEchoServer(loop *falcon.EventLoop, addr string) *EchoServer {
	server := NewTCPServer(loop, addr)
	echoServer := &EchoServer{
		loop:   loop,
		server: server,
	}
	server.ConnectFunc = echoServer.onConnect
	server.DisconnectFunc = echoServer.onDisconnect
	server.ReceiveFunc = echoServer.onReceive
	return echoServer
}

func (this *EchoServer) Start() {
	this.server.Start()
}

func (this *EchoServer) Close() {
	this.server.Close()
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
	loop   *falcon.EventLoop
	client *TCPClient
}

func NewEchoClient(loop *falcon.EventLoop, addr string) *EchoClient {
	client := NewTCPClient(loop, addr)
	echoClient := &EchoClient{
		loop:   loop,
		client: client,
	}
	client.ConnectFunc = echoClient.onConnect
	client.DisconnectFunc = echoClient.onDisconnect
	client.ReceiveFunc = echoClient.onReceive
	return echoClient
}

func (this *EchoClient) Start() {
	this.client.Start()
}

func (this *EchoClient) Close() {
	this.client.Close()
}

func (this *EchoClient) onConnect(connection *TCPConnection) {
	log.Println("client: connected.")
	log.Println("client: send", Message)
	connection.Send([]byte(Message))
}

func (this *EchoClient) onDisconnect(connection *TCPConnection) {
	log.Println("client: disconnected.")
	this.client.Close()
	this.loop.Close()
}

func (this *EchoClient) onReceive(connection *TCPConnection, b []byte) {
	log.Println("client: receive ", string(b))
}

func TestEcho(t *testing.T) {
	loop := falcon.NewEventLoop()
	echoServer := NewEchoServer(loop, "localhost:8000")
	echoServer.Start()

	echoClient := NewEchoClient(loop, "localhost:8000")
	echoClient.Start()

	loop.Loop() // loop

	echoClient.Close()
	echoServer.Close()
}
