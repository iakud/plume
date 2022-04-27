package registery

import (
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

func Add(c *clientv3.Client, service string, addr string) error {
	em, err := endpoints.NewManager(c, service)
	if err != nil {
		return fmt.Errorf("registery: failed to new endpoint manager: %s", err)
	}
	if err := em.AddEndpoint(c.Ctx(), service+"/"+addr, endpoints.Endpoint{Addr: addr}); err != nil {
		return fmt.Errorf("registery: failed to add endpoint: %s", err)
	}
	return nil
}
