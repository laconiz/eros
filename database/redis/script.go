// 脚本

package redis

import (
	"github.com/gomodule/redigo/redis"

	"github.com/laconiz/eros/database/redis/decoder"
)

type Script struct {
	Name   string
	Script *redis.Script
}

type script struct {
	conn *Redis
}

// 预加载脚本
func (s *script) Load(script *Script) error {

	conn := s.conn.pool.Get()
	defer conn.Close()

	return script.Script.Load(conn)
}

func (s *script) Do(script *Script, args ...interface{}) (interface{}, error) {

	conn := s.conn.pool.Get()
	defer conn.Close()

	// 序列化参数
	arguments, err := decoder.FormatArgs(args)
	if err != nil {
		return nil, err
	}

	if s.conn.option.Log {
		// log.Infof("script %v - %v", script.Name, args)
	}

	// 执行脚本
	reply, err := script.Script.Do(conn, arguments...)

	if s.conn.option.Log {
		// log.Infof("result %v - %v", decoder.FormatReply(reply), err)
	}

	return reply, err
}
