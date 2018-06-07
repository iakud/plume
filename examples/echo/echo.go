package echo

import (
	"log"

	"github.com/iakud/falcon/codec"
	"github.com/iakud/falcon/net"
)

type EchoServer struct {
	server *net.TCPServer
}

func NewEchoServer(addr string) *EchoServer {
	server := net.NewTCPServer(addr)
	echoServer := &EchoServer{
		server: server,
	}
	server.ConnectFunc = echoServer.onConnect
	server.DisconnectFunc = echoServer.onDisconnect
	server.ReceiveFunc = echoServer.onReceive
	return echoServer
}

func (this *EchoServer) Start() {
	this.server.Start(this.newConnection)
}

func (this *EchoServer) Close() {
	this.server.Close()
}

func (this *EchoServer) newConnection(connection *net.TCPConnection) {
	connection.ServeCodec(&codec.StdCodec{})
}

func (this *EchoServer) onConnect(connection *net.TCPConnection) {
	log.Println("server: connected.")
}

func (this *EchoServer) onDisconnect(connection *net.TCPConnection) {
	log.Println("server: disconnected.")
}

func (this *EchoServer) onReceive(connection *net.TCPConnection, b []byte) {
	message := string(b)
	log.Println("server: receive", message)
	connection.Send(b)
	connection.Shutdown()
}

var Message string = "hello"

type EchoClient struct {
	client *net.TCPClient
	done   chan struct{}
}

func NewEchoClient(addr string) *EchoClient {
	client := net.NewTCPClient(addr)
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
	this.client.Start(this.newConnection)
}

func (this *EchoClient) Done() {
	<-this.done
}

func (this *EchoClient) newConnection(connection *net.TCPConnection) {
	connection.ServeCodec(&codec.StdCodec{})
}

func (this *EchoClient) onConnect(connection *net.TCPConnection) {
	log.Println("client: connected.")
	log.Println("client: send", Message)
	connection.Send([]byte(Message))
}

func (this *EchoClient) onDisconnect(connection *net.TCPConnection) {
	log.Println("client: disconnected.")
	this.client.Close()
	close(this.done)
}

func (this *EchoClient) onReceive(connection *net.TCPConnection, b []byte) {
	log.Println("client: receive ", string(b))
}
