// 有序集合

package redis

import "github.com/gomodule/redigo/redis"

type ZOrder struct {
	conn *Redis
	key  string
}

// 为成员增加分值
func (z *ZOrder) Incr(member interface{}, increment int64) error {
	return z.conn.void(ZINCRBY, z.key, increment, member)
}

// 返回指定区间内的成员 从0开始
func (z *ZOrder) Range(start, stop int, value interface{}) error {
	return z.conn.complex(value, ZREVRANGE, z.key, start, stop, WITHSCORES)
}

// 获取成员的分值
func (z *ZOrder) Score(member interface{}) (int64, bool, error) {

	// 获取分值
	reply, err := z.conn.Do(ZSCORE, z.key, member)
	if err != nil {
		return 0, false, err
	}

	// 没有此成员记录
	switch reply.(type) {
	case nil:
		return 0, false, nil
	}

	score, err := redis.Int64(reply, err)
	return score, true, err
}

// 获取成员的排名 从0开始
func (z *ZOrder) Rank(member interface{}) (int64, bool, error) {

	// 获取排名
	reply, err := z.conn.Do(ZREVRANK, z.key, member)
	if err != nil {
		return 0, false, err
	}

	// 没有此成员记录
	switch reply.(type) {
	case nil:
		return 0, false, nil
	}

	rank, err := redis.Int64(reply, err)
	return rank, true, nil
}
