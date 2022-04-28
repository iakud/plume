package main

import (
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/iakud/plume/etcd"
)

const (
	kService = "/h-server"
)

var serviceMgr *etcd.ServiceManager

func main() {
	client, err := etcd.New([]string{"localhost:2379"})
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	manager, err := client.NewServiceManager(kService)
	if err != nil {
		log.Fatalln(err)
	}
	watcher, err := manager.NewWatcher(UpdateServers)
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	serviceMgr = manager 

	http.HandleFunc("/", GetServers)
	http.HandleFunc("/addServer", AddServer)
	http.ListenAndServe("localhost:80", nil)
}

func GetServers(w http.ResponseWriter, r *http.Request) {
	WriteServers(w)
}

func AddServer(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	if len(name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	port, err := strconv.Atoi(r.FormValue("port"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	server := &Server{Name: name,Addr: tcpAddr.IP.String(),Port: port}
	address := etcd.Address{Addr: server.Addr, Metadata: server}
	serviceMgr.Add(address)
	w.WriteHeader(http.StatusOK)
}
