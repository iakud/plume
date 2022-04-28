package etcd

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

type ServiceManager struct {
	ctx     context.Context
	em      endpoints.Manager
	service string
}

func newServiceManager(c *clientv3.Client, service string) (*ServiceManager, error) {
	em, err := endpoints.NewManager(c, service)
	if err != nil {
		return nil, fmt.Errorf("etcd: failed to new endpoint manager: %s", err)
	}
	sm := &ServiceManager{
		ctx:     c.Ctx(),
		em:      em,
		service: service,
	}
	return sm, nil
}

func (sm *ServiceManager) Add(address Address) error {
	if err := sm.em.AddEndpoint(sm.ctx, sm.service+"/"+address.Addr, endpoints.Endpoint{Addr: address.Addr, Metadata: address.Metadata}); err != nil {
		return fmt.Errorf("registery: failed to add endpoint: %s", err)
	}
	return nil
}

func (sm *ServiceManager) NewWatcher(handler func([]Address)) (*ServiceWatcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	sw := &ServiceWatcher{
		handler: handler,
		cancel:  cancel,
	}
	wch, err := sm.em.NewWatchChannel(ctx)
	if err != nil {
		return nil, fmt.Errorf("etcd: failed to new watch channer: %s", err)
	}
	sw.wg.Add(1)
	go sw.watch(ctx, wch)
	return sw, nil
}

type ServiceWatcher struct {
	handler func([]Address)
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func (sw *ServiceWatcher) watch(ctx context.Context, wch endpoints.WatchChannel) {
	defer sw.wg.Done()

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
			sw.handler(addrs)
		}
	}
}

func (sw *ServiceWatcher) Close() {
	sw.cancel()
	sw.wg.Wait()
}
