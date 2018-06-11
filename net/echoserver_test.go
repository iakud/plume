package net

import (
	"log"
	"testing"
)

type EchoServer struct {
	server *TCPServer
}

func NewEchoServer(addr string) *EchoServer {
	server := NewTCPServer(addr, &DefaultCodec)
	echoServer := &EchoServer{
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
	log.Println("echo server: connected.")
}

func (this *EchoServer) onDisconnect(connection *TCPConnection) {
	log.Println("echo server: disconnected.")
}

func (this *EchoServer) onReceive(connection *TCPConnection, b []byte) {
	message := string(b)
	log.Println("echo server: receive", message)
	connection.Send(b)
	connection.Shutdown()
}

var Message string = "hello"

type EchoClient struct {
	client *TCPClient
	done   chan struct{}
}

func NewEchoClient(addr string) *EchoClient {
	client := NewTCPClient(addr, &DefaultCodec)
	echoClient := &EchoClient{
		client: client,
		done:   make(chan struct{}),
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
	log.Println("echo client: connected.")
	log.Println("echo client: send", Message)
	connection.Send([]byte(Message))
}

func (this *EchoClient) onDisconnect(connection *TCPConnection) {
	log.Println("echo client: disconnected.")
	this.client.Close()
	close(this.done)
}

func (this *EchoClient) onReceive(connection *TCPConnection, b []byte) {
	log.Println("echo client: receive ", string(b))
}

func (this *EchoClient) Done() {
	<-this.done
}

func TestEcho(t *testing.T) {
	echoServer := NewEchoServer("localhost:8000")
	echoServer.Start()

	echoClient := NewEchoClient("localhost:8000")
	echoClient.Start()

	echoClient.Done()
	echoServer.Close()
}
