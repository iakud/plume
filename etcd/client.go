package etcd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	etcdClient *clientv3.Client
}

func New(addrs []string) (*Client, error) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &Client{etcdClient: c}, nil
}

func (c *Client) Client() *clientv3.Client {
	return c.etcdClient
}

func (c *Client) Close() error {
	return c.etcdClient.Close()
}

func (c *Client) NewServiceManager(service string) (*ServiceManager, error) {
	manager, err := newServiceManager(c.etcdClient, service)
	if err != nil {
		return nil, err
	}
	return manager, nil
}
