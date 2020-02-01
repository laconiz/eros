// 键

package redis

import (
	"github.com/gomodule/redigo/redis"

	"github.com/laconiz/eros/database/redis/decoder"
)

type Key struct {
	conn *Redis
}

// 检查键是否存在
func (k *Key) Exist(key interface{}) (bool, error) {
	return k.conn.bool(EXISTS, key)
}

// 设置键的值
func (k *Key) Set(key, value interface{}) error {
	return k.conn.void(SET, key, value)
}

// 当键不存在时设置键的值
func (k *Key) SetNX(key, value interface{}) (bool, error) {

	reply, err := k.conn.Do(SET, key, value, NX)
	if err != nil {
		return false, err
	}

	str, err := redis.String(reply, nil)
	return err == nil && str == OK, nil
}

// 设置键的值并设置过期时间
func (k *Key) SetEX(key, value interface{}, second int64) error {
	return k.conn.void(SET, key, value, EX, second)
}

// 当键不存在时设置键的值并设置过期时间
func (k *Key) SetNEX(key, value interface{}, second int64) (bool, error) {

	reply, err := k.conn.Do(SET, key, value, EX, second, NX)
	if err != nil {
		return false, err
	}

	str, err := redis.String(reply, nil)
	return err == nil && str == OK, nil
}

// 获取键的值
func (k *Key) Get(key, value interface{}) (bool, error) {

	reply, err := k.conn.Do(GET, key)
	if err != nil {
		return false, err
	}

	switch reply.(type) {
	case nil:
		return false, nil
	}

	return true, decoder.Decode(value, reply, err)
}

// 删除一个或多个键
func (k *Key) Delete(keys ...interface{}) error {
	return k.conn.void(DEL, keys...)
}

// 为键加上增量
func (k *Key) Incr(key interface{}, value int64) (int64, error) {
	return k.conn.int64(INCRBY, key, value)
}
