package main

import (
	"fmt"
	"log"
	"syscall"
	//	"time"
	"os"
	"os/signal"

	"github.com/iakud/falcon/network"
)

type EchoServer struct {
}

func (this *EchoServer) Connected(connection *network.TCPConnection) {
	fmt.Println("Connected")
}

func (this *EchoServer) Disconnected(connection *network.TCPConnection) {
	fmt.Println("Disconnected")
}

func (this *EchoServer) Receive(connection *network.TCPConnection, b []byte) {
	//connection.Send(b)
	//connection.Shutdown()
	//connection.CloseIn(3 * time.Second)
}

func sin() {
	c := make(chan os.Signal)
	signal.Notify(c)
	for s := range c {
		fmt.Println("get signal:", s)
	}
}

func main() {
	go sin()
	server := network.NewTCPServer("127.0.0.1:8000", &EchoServer{})
	if err := server.Start(); err != nil {
		log.Fatalln(err)
	}
	server.Close()
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	select {
	case s := <-c:
		fmt.Println("get signal:", s)
	}

}
