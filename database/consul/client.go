package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"

	"github.com/laconiz/eros/utils/command"
)

const DefaultAddr = "127.0.0.1:8500"

var (
	addr   string
	client = &api.Client{}
)

func init() {

	addr = command.ParseStringArg("consul", DefaultAddr)

	if cli, err := api.NewClient(&api.Config{Address: addr}); err != nil {
		panic(fmt.Errorf("new consul client error: %w", err))
	} else {
		client = cli
	}

	if _, err := client.Catalog().Datacenters(); err != nil {
		panic(fmt.Errorf("check consul connection error: %w", err))
	}
}
