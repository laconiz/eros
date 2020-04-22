package redisMgr

import (
	"github.com/laconiz/eros/database/consul"
	"github.com/laconiz/eros/database/consul/consulor"
	"github.com/laconiz/eros/database/redis"
	"github.com/laconiz/eros/oceanus/example/model"
	"sync"
)

const (
	prefix   = "database/redis/"
	defaults = "default"
)

var (
	conn  *redis.Redis
	mutex sync.Mutex
)

func Connect(name string) (*redis.Redis, error) {

	mutex.Lock()
	defer mutex.Unlock()

	option := &redis.Option{}
	err := consulor.KV().Load(model.OptionPrefix+prefix+name, option)
	if err == nil {
		return redis.New(option)
	}

	if err != consul.ErrNotFound {
		return nil, err
	}

	if conn != nil {
		return conn, nil
	}

	err = consulor.KV().Load(model.OptionPrefix+prefix+defaults, option)
	if err != nil {
		return nil, err
	}

	conn, err = redis.New(option)
	return conn, err
}
