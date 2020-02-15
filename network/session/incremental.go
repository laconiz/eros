package session

import "sync/atomic"

var incremental uint64

func Increment() uint64 {
	return atomic.AddUint64(&incremental, 1)
}
