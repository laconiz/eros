package consulor

import (
	"github.com/laconiz/eros/database/consul"
	"github.com/laconiz/eros/utils/command"
)

func KV() *consul.KV {
	return client.KV()
}

func Watcher() *consul.Watcher {
	return client.Watcher()
}

var client *consul.Consul

const defaultAddr = "192.168.0.110:8500"

func init() {

	addr := command.ParseStringArg("consul", defaultAddr)

	var err error
	if client, err = consul.New(addr); err != nil {
		panic(err)
	}
}
