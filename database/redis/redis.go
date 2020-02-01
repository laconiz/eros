package redis

import (
	"github.com/gomodule/redigo/redis"

	"github.com/laconiz/eros/database/redis/decoder"
)

type Redis struct {
	pool *redis.Pool
	conf Config
}

// 执行redis命令
func (r *Redis) Do(cmd string, args ...interface{}) (interface{}, error) {

	// 获取连接
	conn := r.pool.Get()
	defer conn.Close()

	// 序列化参数
	arguments, err := decoder.FormatArgs(args)
	if err != nil {
		return nil, err
	}

	if r.conf.Log {
		// log.Infof("request: %v - %v", cmd, arguments)
	}

	// 执行命令
	reply, err := conn.Do(cmd, arguments...)

	if r.conf.Log {
		// log.Infof("response: %v - %v", decoder.FormatReply(reply), err)
	}

	return reply, err
}

func (r *Redis) void(cmd string, args ...interface{}) error {
	_, err := r.Do(cmd, args...)
	return err
}

func (r *Redis) bool(cmd string, args ...interface{}) (bool, error) {
	return redis.Bool(r.Do(cmd, args...))
}

func (r *Redis) string(cmd string, args ...interface{}) (string, error) {
	return redis.String(r.Do(cmd, args...))
}

func (r *Redis) int64(cmd string, args ...interface{}) (int64, error) {
	return redis.Int64(r.Do(cmd, args...))
}

func (r *Redis) complex(recv interface{}, cmd string, args ...interface{}) error {
	reply, err := r.Do(cmd, args...)
	return decoder.Decode(recv, reply, err)
}

func (r *Redis) Key() *Key {
	return &Key{conn: r}
}

func (r *Redis) Hash(key string) *Hash {
	return &Hash{conn: r, key: key}
}

func (r *Redis) ZOrder(key string) *ZOrder {
	return &ZOrder{conn: r, key: key}
}

func (r *Redis) Set(key string) *Set {
	return &Set{redis: r, key: key}
}

func (r *Redis) Script() *script {
	return &script{conn: r}
}

func (r *Redis) Singleton(key string) *Singleton {
	return &Singleton{conn: r, key: key}
}

func (r *Redis) Atomic(key string) *Atomic {
	return (&Atomic{conn: r, key: key}).Expired(3).Timeout(6).Ticker(50)
}

type Config struct {
	Network   string // 网络类型
	Address   string // 地址
	Password  string // 密码
	Database  int    // 数据库
	MaxIdle   int    // 最大空闲连接数
	MaxActive int    // 最大活跃连接数
	Log       bool   // 显示日志
}

func New(conf Config) (*Redis, error) {

	dial := func() (redis.Conn, error) {
		return redis.Dial(
			conf.Network,
			conf.Address,
			redis.DialPassword(conf.Password),
			redis.DialDatabase(conf.Database),
		)
	}

	r := &Redis{
		pool: &redis.Pool{
			Dial:            dial,
			TestOnBorrow:    nil,
			MaxIdle:         conf.MaxIdle,
			MaxActive:       conf.MaxActive,
			IdleTimeout:     0,
			Wait:            true,
			MaxConnLifetime: 0,
		},
		conf: conf,
	}

	if _, err := r.string(PING); err != nil {
		return nil, err
	}

	return r, nil
}
