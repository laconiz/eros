package sql

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Option struct {
	Dialect     string
	Arg         string
	MaxConn     int
	IdleConn    int
	LifeSeconds int
	SingleTable bool
}

func New(option *Option) (*gorm.DB, error) {

	client, err := gorm.Open(option.Dialect, option.Arg)
	if err != nil {
		return nil, err
	}

	client.SingularTable(option.SingleTable)
	client.DB().SetMaxOpenConns(option.MaxConn)
	client.DB().SetMaxIdleConns(option.IdleConn)
	client.DB().SetConnMaxLifetime(time.Duration(option.LifeSeconds) * time.Second)

	return client, nil
}
