// 分布式锁

package redis

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/satori/go.uuid"
)

var (
	ErrAtomicLockFailed   = errors.New("lock failed")   // 加锁失败
	ErrAtomicUnlockFailed = errors.New("unlock failed") // 解锁失败
)

// Redis分布式锁
type Atomic struct {
	conn    *Redis        // 连接
	key     string        // 锁名
	expired int64         // 过期时间
	ticker  time.Duration // 重试频率
	timeout time.Duration // 超时时间
	value   string        // uuid 解锁需要的value
}

// Key过期时间
func (a *Atomic) Expired(s int64) *Atomic {
	a.expired = s
	return a
}

// 重试频率
func (a *Atomic) Ticker(ms int64) *Atomic {
	a.ticker = time.Duration(ms) * time.Millisecond
	return a
}

// 加锁超时时间
func (a *Atomic) Timeout(s int64) *Atomic {
	a.timeout = time.Duration(s) * time.Second
	return a
}

// 获取分布式锁
func (a *Atomic) Lock() error {

	// 设置随机值
	a.value = uuid.NewV5(uuid.NewV4(), a.key).String()

	//
	ticker := time.NewTicker(a.ticker)
	defer ticker.Stop()
	// 超时时间
	deadline := time.Now().Add(a.timeout)

	key := a.conn.Key()

	for {

		// 尝试加锁
		ok, err := key.SetNEX(a.key, a.value, a.expired)
		if err != nil {
			// 加锁错误
			return err
		} else if ok {
			// 加锁成功
			return nil
		}

		// 加锁超时
		if deadline.Before(time.Now()) {
			return ErrAtomicLockFailed
		}

		// 加锁间隔
		<-ticker.C
	}
}

// 释放分布式锁
func (a *Atomic) Unlock() error {

	// 执行脚本
	script := a.conn.Script()
	reply, err := script.Do(unlockAtomic, a.key, a.value)

	// 执行错误
	ok, err := redis.Bool(reply, err)
	if err != nil {
		return err
	}

	// 解锁失败
	if !ok {
		return ErrAtomicUnlockFailed
	}

	return nil
}

// 执行分布式事务
func (a *Atomic) Exec(f func()) (bool, error) {

	if err := a.Lock(); err != nil {
		return false, err
	}
	defer a.Unlock()

	defer func() {
		if err := recover(); err != nil {
			// log.Errorf("atomic execute error: %v", err)
		}
	}()

	f()

	return true, a.Unlock()
}

// 释放分布式锁的脚本
var unlockAtomic = &Script{
	Name:   "ATOMIC UNLOCK",
	Script: redis.NewScript(1, luaUnlockAtomic),
}

var luaUnlockAtomic = `

	if redis.call('GET', KEYS[1]) == ARGV[1] then
		return redis.call('DEL', KEYS[1])
	end

	return 0
`
