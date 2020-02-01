// 哈希表

package redis

import (
	"github.com/gomodule/redigo/redis"

	"github.com/laconiz/eros/database/redis/decoder"
)

type Hash struct {
	conn *Redis // redis连接
	key  string // 表名
}

// 删除表中的一个或多个域
func (h *Hash) Delete(fields ...interface{}) error {
	return h.conn.void(HDEL, append([]interface{}{h.key}, fields...)...)
}

// 为表中的域设置值
func (h *Hash) Set(field, value interface{}) error {
	return h.conn.void(HSET, h.key, field, value)
}

// 获取表中指定域的值
func (h *Hash) Get(field, value interface{}) (bool, error) {

	// 获取值
	reply, err := h.conn.Do(HGET, h.key, field)
	if err != nil {
		return false, err
	}

	// 没有此域记录
	switch reply.(type) {
	case nil:
		return false, nil
	}

	// 反序列化数据
	return true, decoder.Decode(value, reply, err)
}

// 获取表中一个或多个域的值
func (h *Hash) Gets(value interface{}, fields ...interface{}) error {
	return h.conn.complex(value, HMGET, append([]interface{}{h.key}, fields...)...)
}

// 获取表中所有的域和值
func (h *Hash) GetAll(value interface{}) error {
	return h.conn.complex(value, HGETALL, h.key)
}

// 为表中指定域的值加上增量
func (h *Hash) Incr(field interface{}, increment int64) (int64, error) {
	return h.conn.int64(HINCRBY, h.key, field, increment)
}

// 为表中指定域加上增量
// 若增量为负数且与原值计算后为负数则不会进行计算
func (h *Hash) UnsignedIncr(field interface{}, increment int64) (int64, bool, error) {

	// 增量为负数
	if increment < 0 {

		result := &struct {
			Value   int64
			Success bool
		}{}

		// 执行脚本并反序列化数据
		replies, err := h.conn.Script().Do(unsignedIncr, h.key, field, increment)
		err = decoder.Decode(result, replies, err)
		return result.Value, result.Success, err
	}

	// 直接进行计算
	value, err := h.Incr(field, increment)
	success := err == nil
	return value, success, err
}

// 无符号增量计算脚本
var unsignedIncr = &Script{
	Name:   "HASH UNSIGNED INCR",
	Script: redis.NewScript(1, luaUnsignedIncr),
}

var luaUnsignedIncr = `

	local new = redis.call('HINCRBY', KEYS[1], ARGV[1], ARGV[2])
	if new >= 0 then
		return {new, true}
	end

	new = redis.call('HINCRBY', KEYS[1], ARGV[1], -ARGV[2])
	return {new, false}
`
