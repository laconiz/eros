package incremental

import "sync/atomic"

var value int64 = 0

const step = 1

func Get() int64 {
	return atomic.AddInt64(&value, step)
}
