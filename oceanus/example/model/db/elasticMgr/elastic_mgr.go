package elasticMgr

import (
	"github.com/laconiz/eros/database/consul"
	"github.com/laconiz/eros/database/consul/consulor"
	"github.com/laconiz/eros/database/elastic"
	"sync"
)

const (
	prefix   = "config/database/elastic/"
	defaults = "default"
)

var (
	conn  *elastic.Elastic
	mutex sync.Mutex
)

func Connect(name string) (*elastic.Elastic, error) {

	mutex.Lock()
	defer mutex.Unlock()

	option := &elastic.Option{}
	err := consulor.KV().Load(prefix+name, option)
	if err == nil {
		return elastic.New(option)
	}

	if err != consul.ErrNotFound {
		return nil, err
	}

	if conn != nil {
		return conn, nil
	}

	err = consulor.KV().Load(prefix+defaults, option)
	if err != nil {
		return nil, err
	}

	conn, err = elastic.New(option)
	return conn, err
}
