package net

import (
	"log"
	"testing"
)

type EchoServer struct {
	server *TCPServer
}

func NewEchoServer(addr string) *EchoServer {
	echoServer := &EchoServer{}
	echoServer.server = NewTCPServer(addr, echoServer, DefaultCodec)
	return echoServer
}

func (this *EchoServer) ListenAndServe() {
	if err := this.server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func (this *EchoServer) Close() {
	if err := this.server.Close(); err != nil {
		log.Println(err)
	}
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
	done   chan struct{}
}

func NewEchoClient(addr string) *EchoClient {
	echoClient := &EchoClient{
		done: make(chan struct{}),
	}
	echoClient.client = NewTCPClient(addr, echoClient, DefaultCodec)
	return echoClient
}

func (this *EchoClient) ConnectAndServe() {
	if err := this.client.ConnectAndServe(); err != nil {
		log.Println(err)
	}
}

func (this *EchoClient) Done() {
	<-this.done
}

func (this *EchoClient) Connect(connection *TCPConnection) {
	log.Println("echo client: connected.")
	log.Println("echo client: send", Message)
	connection.Send([]byte(Message))
}

func (this *EchoClient) Disconnect(connection *TCPConnection) {
	log.Println("echo client: disconnected.")

	if err := this.client.Close(); err != nil {
		log.Println(err)
	}
	close(this.done)
}

func (this *EchoClient) Receive(connection *TCPConnection, b []byte) {
	log.Println("echo client: receive", string(b))
}

func TestEcho(t *testing.T) {
	echoServer := NewEchoServer("localhost:8000")
	go echoServer.ListenAndServe()

	echoClient := NewEchoClient("localhost:8000")
	go echoClient.ConnectAndServe()

	echoClient.Done()
	echoServer.Close()
}
