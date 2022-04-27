package discovery

import (
	"context"
	"fmt"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type Address struct {
	Addr     string
	Metadata interface{}
}

type Discovery struct {
	c      *clientv3.Client
	update func([]Address)
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New(c *clientv3.Client, service string, update func(addresses []Address)) (*Discovery, error) {
	ctx, cancel := context.WithCancel(context.Background())
	d := &Discovery{
		update: update,
		cancel: cancel,
	}

	em, err := endpoints.NewManager(c, service)
	if err != nil {
		return nil, fmt.Errorf("discover: failed to new endpoint manager: %s", err)
	}
	wch, err := em.NewWatchChannel(ctx)
	if err != nil {
		return nil, fmt.Errorf("discover: failed to new watch channer: %s", err)
	}

	d.wg.Add(1)
	go d.watch(ctx, wch)
	return d, nil
}

func (d *Discovery) watch(ctx context.Context, wch endpoints.WatchChannel) {
	defer d.wg.Done()

	allUps := make(map[string]*endpoints.Update)
	for {
		select {
		case <-ctx.Done():
			return
		case ups, ok := <-wch:
			if !ok {
				return
			}

			for _, up := range ups {
				switch up.Op {
				case endpoints.Add:
					allUps[up.Key] = up
				case endpoints.Delete:
					delete(allUps, up.Key)
				}
			}

			var addrs []Address
			for _, up := range allUps {
				addr := Address{
					Addr:     up.Endpoint.Addr,
					Metadata: up.Endpoint.Metadata,
				}
				addrs = append(addrs, addr)
			}
			d.update(addrs)
		}
	}
}

func (d *Discovery) Close() {
	d.cancel()
	d.wg.Wait()
}
