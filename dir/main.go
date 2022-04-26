package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/iakud/plume/etcd/discovery"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var services []discovery.Address

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		log.Fatalln(err)
	}
	defer cli.Close()

	serviced, err := discovery.New(cli, "/myservice", func(addresses []discovery.Address) {
		services = addresses
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer serviced.Close()

	http.HandleFunc("/", Services)
	http.ListenAndServe("localhost:80", nil)
}

func Services(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(services)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}