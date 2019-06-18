package net

import (
	"log"
	"testing"
)

type EchoServer struct {
	server *TCPServer
}

func NewEchoServer(addr string) *EchoServer {
	echoServer := &EchoServer{
		server: NewTCPServer(addr),
	}
	return echoServer
}

func (this *EchoServer) ListenAndServe() {
	if err := this.server.ListenAndServe(this, nil); err != nil {
		log.Println(err)
	}
}

func (this *EchoServer) Close() {
	this.server.Close()
}

func (this *EchoServer) Connect(connection *TCPConnection) {
	log.Println("echo server: connected.")
}

func (this *EchoServer) Disconnect(connection *TCPConnection) {
	log.Println("echo server: disconnected.")
}

func (this *EchoServer) Receive(connection *TCPConnection, b []byte) {
	message := string(b)
	log.Println("echo server: receive", message)
	connection.Send(b)
	connection.Shutdown()
}

var Message string = "hello"

type EchoClient struct {
	client *TCPClient
}

func NewEchoClient(addr string) *EchoClient {
	echoClient := &EchoClient{
		client: NewTCPClient(addr),
	}
	return echoClient
}

func (this *EchoClient) ConnectAndServe() {
	this.client.EnableRetry() // 启用retry
	if err := this.client.DialAndServe(this, nil); err != nil {
		log.Println(err)
	}
}

func (this *EchoClient) Connect(connection *TCPConnection) {
	log.Println("echo client: connected.")
	log.Println("echo client: send", Message)
	connection.Send([]byte(Message))
}

func (this *EchoClient) Disconnect(connection *TCPConnection) {
	log.Println("echo client: disconnected.")
	this.client.Close()
}

func (this *EchoClient) Receive(connection *TCPConnection, b []byte) {
	log.Println("echo client: receive", string(b))
}

func TestEcho(t *testing.T) {
	echoServer := NewEchoServer("localhost:8000")
	go func() {
		echoClient := NewEchoClient("localhost:8000")
		echoClient.ConnectAndServe()
		echoServer.Close()
	}()
	echoServer.ListenAndServe()
}
