package network_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/iakud/plume/network"
)

type wsEchoServer struct {
}

func newWSEchoServer() *wsEchoServer {
	echoServer := &wsEchoServer{}
	return echoServer
}

func (srv *wsEchoServer) Connect(connection *network.WSConnection, connected bool) {
	if connected {
		log.Printf("echo server: %v connected.\n", connection.RemoteAddr())
	} else {
		log.Printf("echo server: %v disconnected.\n", connection.RemoteAddr())
	}
}

func (srv *wsEchoServer) Receive(connection *network.WSConnection, data []byte) {
	message := string(data)
	log.Printf("echo server: %v receive %v\n", connection.RemoteAddr(), message)
	log.Println("echo server: send", message)
	connection.Send(data)
	connection.Shutdown()
}

type wsEchoClient struct {
	Client *network.WSClient
}

func newWSEchoClient() *wsEchoClient {
	echoClient := &wsEchoClient{}
	return echoClient
}

func (c *wsEchoClient) Connect(connection *network.WSConnection, connected bool) {
	const message string = "hello"
	if connected {
		log.Printf("echo client: %v connected.\n", connection.RemoteAddr())
		log.Println("echo client: send", message)
		connection.Send([]byte(message))
	} else {
		log.Printf("echo client: %v disconnected.\n", connection.RemoteAddr())
		c.Client.Close()
	}
}

func (c *wsEchoClient) Receive(connection *network.WSConnection, data []byte) {
	log.Printf("echo client: %v receive %v\n", connection.RemoteAddr(), string(data))
}

func TestWSEcho(t *testing.T) {
	log.Println("test start")
	server := newWSEchoServer()
	wsServer := network.NewWSServer(server)
	httpServer := &http.Server{Addr: "localhost:8000"}
	go func() {
		client := newWSEchoClient()
		wsClient := network.NewWSClient("ws://localhost:8000", client)
		client.Client = wsClient
		wsClient.EnableRetry() // 启用Retry
		if err := wsClient.DialAndServe(); err != nil {
			log.Println(err)
		}
		wsServer.Close()
		httpServer.Close()
	}()

	http.Handle("/", network.WebsocketHandler(wsServer))
	httpServer.ListenAndServe()
}
