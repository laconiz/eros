package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/eros/log"
	"os"
)

var (
	addr   string
	client = &api.Client{}
)

func init() {

	for idx, arg := range os.Args {
		if arg == "-consul" && idx+1 < len(os.Args) {
			addr = os.Args[idx+1]
			break
		}
	}

	if addr == "" {
		addr = os.Getenv("CONSUL_HOST")
	}

	if addr == "" {
		addr = "127.0.0.1:8500"
	}

	logger.Infof("connect to %s", addr)

	var err error
	client, err = api.NewClient(&api.Config{
		Address: addr,
	})
	if err != nil {
		panic(fmt.Errorf("new consul client error: %w", err))
	}

	if _, err := client.Catalog().Datacenters(); err != nil {
		panic(fmt.Errorf("check consul client error: %w", err))
	}
}

var logger = log.Std("consul")
