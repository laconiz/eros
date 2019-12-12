package redis

import (
	"sync"
	"testing"
	"time"
)

func TestAtomic(t *testing.T) {

	assert(r.Key().Delete(key) == nil)

	// 正常执行
	executed, err := r.Atomic(key).Exec(func() {})
	assert(err == nil, executed)

	// 执行超时
	executed, err = r.Atomic(key).Expired(1).Timeout(2).Ticker(100).Exec(func() {
		time.Sleep(time.Millisecond * 1100)
	})
	assert(err == ErrAtomicUnlockFailed, executed)

	// key被占用
	assert(r.Key().Set(key, 1) == nil)
	executed, err = r.Atomic(key).Exec(func() {})
	assert(err == ErrAtomicLockFailed, !executed)

	assert(r.Key().Delete(key) == nil)

	// 测试分布式
	var count int32
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			executed, err := r.Atomic(key).Exec(func() {
				count++
				time.Sleep(time.Millisecond * 100)
			})
			assert(err == nil, executed)
			wg.Done()
		}()
	}

	wg.Wait()
	assert(count == 20)
}
