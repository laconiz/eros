package sql

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 连接参数
type Option struct {
	Dialect     string // 连接类型
	Arg         string // 连接参数
	MaxConn     int    // 最大连接数
	IdleConn    int    // 空闲连接数
	LifeSeconds int    // 连接生命期
	SingleTable bool   // 表命名规则
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
