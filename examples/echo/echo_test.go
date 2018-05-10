package echo

import (
	"log"
	"testing"
	"time"
	//"time"
)

func TestEcho(t *testing.T) {
	echoServer := NewEchoServer("127.0.0.1:9000")
	if err := echoServer.Start(); err != nil {
		log.Fatalln(err)
	}
	echoClient := NewEchoClient("127.0.0.1:9000")
	if err := echoClient.Start(); err != nil {
		log.Fatalln(err)
	}
	//time.Sleep(time.Second)
	select {
	case <-time.After(3 * time.Second):
	case <-echoClient.done:
	}
	//echoClient.Wait()
	echoServer.Close()
}
