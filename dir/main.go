package main

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
	}
	defer cli.Close()

	em, _ := endpoints.NewManager(cli, "etcd:///")
	ch, _ := em.NewWatchChannel(cli.Ctx())
	allEndPoints := make(map[string]endpoints.Endpoint)
	for {
		select {
			case updates, ok := <-ch:
				if !ok {
					return
				}
				for _, update := range updates {
					switch update.Op {
					case endpoints.Add:
						allEndPoints[update.Key] = update.Endpoint
					case endpoints.Delete:
						delete(allEndPoints, update.Key)
					default:
						// do nothing
					}
				}
		}
	}
}

func etcdAdd(c *clientv3.Client, service, addr string) error {
	em, _ := endpoints.NewManager(c, service)
	return em.AddEndpoint(c.Ctx(), service+"/"+addr, endpoints.Endpoint{Addr:addr});
}

func etcdDial(c *clientv3.Client, service string) (*grpc.ClientConn, error) {
	etcdResolver, err := resolver.NewBuilder(c);
	if err != nil { return nil, err }
	return  grpc.Dial("etcd:///" + service, grpc.WithResolvers(etcdResolver))
}