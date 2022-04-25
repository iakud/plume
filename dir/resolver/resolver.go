package resolver

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

type Resolver struct {
	c  *clientv3.Client
	wch    endpoints.WatchChannel
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (r *Resolver) watch() {
	defer r.wg.Done()

	allUps := make(map[string]*endpoints.Update)
	allEndPoints := make(map[string]endpoints.Endpoint)
	for {
		select {
		case <-r.ctx.Done():
			return
		case ups, ok := <-r.wch:
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

			addrs := convertToGRPCAddress(allUps)
			r.cc.UpdateState(gresolver.State{Addresses: addrs})
		}
	}
}

func (r *resolver) Close() {
	r.cancel()
	r.wg.Wait()
}