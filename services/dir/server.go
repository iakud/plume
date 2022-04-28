package main

import (
	"encoding/json"
	"io"
	"log"
	"sync/atomic"
	"unsafe"

	"github.com/iakud/plume/etcd"
)

var cacheServers unsafe.Pointer

type Server struct {
	Name string `json:"name`
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

func UpdateServers(addresses []etcd.Address) {
	var servers []*Server
	for _, address := range addresses {
		server, ok := address.Metadata.(*Server)
		if !ok {
			log.Printf("dir: failed to update address: %s", address)
			continue
		}
		servers = append(servers, server)
	}

	var info struct {
		servers []*Server `json:"server"`
	}
	info.servers = servers
	data, err := json.Marshal(info)
	if err != nil {
		log.Printf("dir: failed to json marshal: %s", err)
	}
	atomic.StorePointer(&cacheServers, unsafe.Pointer(&data))
}

func WriteServers(w io.Writer) {
	data := *(*[]byte)(atomic.LoadPointer(&cacheServers))
	w.Write(data)
}
