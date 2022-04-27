package main

import (
	"log"
	"testing"
	"time"

	"github.com/iakud/plume/etcd/registery"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestAddEndpoints(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		log.Fatalln(err)

	}
	defer cli.Close()
	registery.Add(cli, "/myservice", "127.0.0.2")
}
