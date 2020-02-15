package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

func New(addr string) (*Consul, error) {

	client, err := api.NewClient(&api.Config{Address: addr})
	if err != nil {
		return nil, fmt.Errorf("new consul client error: %w", err)
	}

	if _, err = client.Catalog().Datacenters(); err != nil {
		return nil, fmt.Errorf("check consul connection error: %w", err)
	}

	return &Consul{addr: addr, Client: client}, nil
}

type Consul struct {
	addr string
	*api.Client
}

func (c *Consul) KV() *KV {
	return &KV{KV: c.Client.KV()}
}

func (c *Consul) Watcher() *Watcher {
	return &Watcher{consul: c}
}
