// 集合

package redis

type Set struct {
	redis *Redis
	key   string
}

// 将一个或多个成员加入到集合中
func (s *Set) Add(members ...interface{}) error {
	return s.redis.void(SADD, append([]interface{}{s.key}, members...)...)
}

// 移除集合中的一个或多个成员
func (s *Set) Remove(members ...interface{}) error {
	return s.redis.void(SREM, append([]interface{}{s.key}, members...)...)
}

// 获取集合中的所有成员
func (s *Set) Keys(value interface{}) error {
	return s.redis.complex(value, SMEMBERS, s.key)
}
