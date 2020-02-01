// 分布式单次执行锁

package redis

import (
	"time"
)

type Singleton struct {
	conn *Redis
	key  string
}

func (s *Singleton) Exec(f func()) (bool, error) {

	value := time.Now().String()

	// 尝试加锁
	ok, err := s.conn.Key().SetNX(s.key, value)
	if err != nil {
		return false, err
	}

	// 加锁成功 执行函数
	if ok {
		f()
	}

	return ok, nil
}
