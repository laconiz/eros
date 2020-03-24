package postgresMgr

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/laconiz/eros/database/consul"
	"github.com/laconiz/eros/database/consul/consulor"
	"github.com/laconiz/eros/database/sql"
	"sync"
)

const (
	prefix   = "config/database/sql/"
	defaults = "default"
)

var (
	conn  *gorm.DB
	mutex sync.Mutex
)

func Connect(name string) (*gorm.DB, error) {

	mutex.Lock()
	defer mutex.Unlock()

	option := &sql.Option{}
	err := consulor.KV().Load(prefix+name, option)
	if err == nil {
		return sql.New(option)
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

	conn, err = sql.New(option)
	return conn, err
}
